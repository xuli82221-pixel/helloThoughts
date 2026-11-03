package memo

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "golang.org/x/image/webp"
)

const MaxImageSize = 10 * 1024 * 1024
const maxImagePixels = 16_000_000

func (s *Store) assetVersion() string {
	parts := []string{}
	fs.WalkDir(s.root.FS(), ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if p != "." && strings.HasPrefix(d.Name(), ".") {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		if d.Type()&os.ModeSymlink != 0 {
			return nil
		}
		switch strings.ToLower(path.Ext(p)) {
		case ".png", ".jpg", ".jpeg", ".webp":
			if st, e := d.Info(); e == nil && st.Mode().IsRegular() {
				parts = append(parts, fmt.Sprintf("%s:%d:%d", p, st.Size(), st.ModTime().UnixNano()))
			}
		}
		return nil
	})
	sort.Strings(parts)
	return digest([]byte(strings.Join(parts, "\n")))
}

type ImageAsset struct {
	Path   string `json:"path"`
	Source string `json:"source"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	MIME   string `json:"mime"`
}

func inspectImage(body []byte) (string, string, image.Config, error) {
	bad := func() (string, string, image.Config, error) {
		return "", "", image.Config{}, problem(422, "invalid_image", "Use a valid PNG, JPEG or WebP image, up to 10 MiB and 16 megapixels.")
	}
	if len(body) == 0 || len(body) > MaxImageSize {
		return bad()
	}
	config, format, err := image.DecodeConfig(bytes.NewReader(body))
	if err != nil || config.Width < 1 || config.Height < 1 || config.Width > 8192 || config.Height > 8192 || int64(config.Width)*int64(config.Height) > maxImagePixels {
		return bad()
	}
	ext, mime := "", ""
	switch format {
	case "png":
		ext, mime = ".png", "image/png"
	case "jpeg":
		ext, mime = ".jpg", "image/jpeg"
	case "webp":
		ext, mime = ".webp", "image/webp"
	default:
		return bad()
	}
	decoded, actual, err := image.Decode(bytes.NewReader(body))
	if err != nil || actual != format || decoded.Bounds().Dx() != config.Width || decoded.Bounds().Dy() != config.Height {
		return bad()
	}
	return ext, mime, config, nil
}

func (s *Store) AddImage(notePath string, body []byte) (ImageAsset, error) {
	ext, mime, config, err := inspectImage(body)
	if err != nil {
		return ImageAsset{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, err = s.read(notePath); err != nil {
		return ImageAsset{}, err
	}
	base, _ := classify(notePath)
	if base == "" {
		base = "Inbox"
	}
	dir := path.Join(base, "assets")
	if err = s.root.MkdirAll(dir, 0755); err != nil {
		return ImageAsset{}, storageError(err)
	}
	target := path.Join(dir, time.Now().UTC().Format("20060102")+"-"+randomID()[:24]+ext)
	f, err := s.root.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return ImageAsset{}, storageError(err)
	}
	ok := false
	defer func() {
		f.Close()
		if !ok {
			s.root.Remove(target)
		}
	}()
	if _, err = f.Write(body); err != nil {
		return ImageAsset{}, storageError(err)
	}
	if err = f.Sync(); err != nil {
		return ImageAsset{}, storageError(err)
	}
	if err = f.Close(); err != nil {
		return ImageAsset{}, storageError(err)
	}
	rel, err := filepath.Rel(filepath.FromSlash(path.Dir(notePath)), filepath.FromSlash(target))
	if err != nil {
		return ImageAsset{}, storageError(err)
	}
	ok = true
	return ImageAsset{target, filepath.ToSlash(rel), config.Width, config.Height, mime}, nil
}

func resolveImagePath(notePath, source string) (string, error) {
	bad := problem(400, "invalid_image_path", "Use a relative PNG, JPEG or WebP path inside the library.")
	if validPath(notePath) != nil || len(source) > 1500 || strings.ContainsAny(source, "\\\x00") {
		return "", bad
	}
	u, err := url.Parse(source)
	if err != nil || u.IsAbs() || u.Host != "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || strings.HasPrefix(u.Path, "/") || strings.ContainsAny(u.Path, "\\:\x00") {
		return "", bad
	}
	p := path.Join(path.Dir(notePath), u.Path)
	if !fs.ValidPath(p) || strings.HasPrefix(p, "../") {
		return "", bad
	}
	for _, part := range strings.Split(p, "/") {
		if strings.HasPrefix(part, ".") || strings.TrimRight(part, ". ") != part || strings.IndexFunc(part, func(r rune) bool { return r < 32 }) >= 0 {
			return "", bad
		}
	}
	switch strings.ToLower(path.Ext(p)) {
	case ".png", ".jpg", ".jpeg", ".webp":
	default:
		return "", bad
	}
	return p, nil
}

func (s *Store) ReadImage(notePath, source string) ([]byte, string, error) {
	p, err := resolveImagePath(notePath, source)
	if err != nil {
		return nil, "", err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, err = s.read(notePath); err != nil {
		return nil, "", err
	}
	return s.readImageFile(p)
}

func (s *Store) readImageFile(p string) ([]byte, string, error) {
	f, err := s.root.Open(p)
	if err != nil {
		return nil, "", storageError(err)
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil || !st.Mode().IsRegular() || st.Size() > MaxImageSize {
		return nil, "", problem(422, "invalid_image", "Invalid local image.")
	}
	body, err := io.ReadAll(io.LimitReader(f, MaxImageSize+1))
	if err != nil {
		return nil, "", storageError(err)
	}
	_, mime, _, err := inspectImage(body)
	return body, mime, err
}

func (s *Server) registerImages(mux *http.ServeMux) {
	mux.Handle("POST /api/images", s.authorized(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, MaxImageSize)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fail(w, problem(413, "image_too_large", "Images must be at most 10 MiB."))
			return
		}
		value, err := s.store.AddImage(r.URL.Query().Get("note"), body)
		if err != nil {
			fail(w, err)
			return
		}
		s.broadcast()
		jsonResponse(w, 201, value)
	})))
	mux.Handle("GET /api/images", s.authorized(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, mime, err := s.store.ReadImage(r.URL.Query().Get("note"), r.URL.Query().Get("src"))
		if err != nil {
			fail(w, err)
			return
		}
		w.Header().Set("Content-Type", mime)
		w.Header().Set("Content-Disposition", "inline")
		w.Write(body)
	})))
}
