package memo

import (
	"bytes"
	"errors"
	"os"
	"path"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// Moves preserve Markdown bytes and copy referenced local images to the paths
// those same relative references resolve to at the destination. Shared originals
// are kept; an existing destination attachment is reused only if bytes match.
func (s *Store) Move(source, directory, revision string) (Note, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	note, err := s.read(source)
	if err != nil {
		return Note{}, err
	}
	if revision == "" || note.Revision != revision {
		return Note{}, problem(409, "revision_conflict", "The file changed.")
	}
	target := path.Join(directory, path.Base(source))
	if err = validPath(target); err != nil {
		return Note{}, err
	}
	if target == source {
		return note, nil
	}
	info, err := s.root.Stat(directory)
	if err != nil || !info.IsDir() {
		return Note{}, problem(400, "invalid_path", "Choose an existing destination folder.")
	}
	if err = s.uniqueName(target, ""); err != nil {
		return Note{}, err
	}
	seen := map[string]bool{}
	document := goldmark.DefaultParser().Parse(text.NewReader([]byte(note.Content)))
	err = ast.Walk(document, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		image, ok := node.(*ast.Image)
		if !entering || !ok {
			return ast.WalkContinue, nil
		}
		ref := string(util.ResolveEntityNames(util.ResolveNumericReferences(util.UnescapePunctuations(image.Destination))))
		if seen[ref] {
			return ast.WalkContinue, nil
		}
		seen[ref] = true
		if strings.HasPrefix(ref, "/") || strings.Contains(ref, ":") {
			return ast.WalkContinue, nil
		}
		original, e := resolveImagePath(source, ref)
		if e != nil {
			return ast.WalkStop, e
		}
		destination, e := resolveImagePath(target, ref)
		if e != nil {
			return ast.WalkStop, e
		}
		if original == destination {
			return ast.WalkContinue, nil
		}
		body, _, e := s.readImageFile(original)
		if e != nil {
			var p *Error
			if errors.As(e, &p) && p.Status == 404 {
				return ast.WalkContinue, nil
			}
			return ast.WalkStop, e
		}
		if _, e = s.root.Stat(destination); e == nil {
			existing, _, e := s.readImageFile(destination)
			if e != nil || !bytes.Equal(existing, body) {
				return ast.WalkStop, problem(409, "image_collision", "A different image occupies the destination. Both notes were kept.")
			}
			return ast.WalkContinue, nil
		} else if !errors.Is(e, os.ErrNotExist) {
			return ast.WalkStop, storageError(e)
		}
		if e = s.publishImage(destination, body); e != nil {
			return ast.WalkStop, e
		}
		return ast.WalkContinue, nil
	})
	if err != nil {
		return Note{}, err
	}
	entry, err := s.stage(source, revision, "move")
	if err != nil {
		return Note{}, err
	}
	payload, metadata := trashPaths(entry.ID)
	if err = s.root.Link(payload, target); err != nil {
		s.rollbackStage(entry)
		return Note{}, storageError(err)
	}
	moved, err := s.read(target)
	if err != nil {
		return Note{}, problem(409, "operation_conflict", "The captured note remains recoverable in Trash.")
	}
	s.root.Remove(payload)
	s.root.Remove(metadata)
	return moved, nil
}

func (s *Store) publishImage(target string, body []byte) error {
	if err := s.root.MkdirAll(path.Dir(target), 0755); err != nil {
		return storageError(err)
	}
	tmp := path.Join(path.Dir(target), ".image-"+randomID()+".tmp")
	f, err := s.root.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return storageError(err)
	}
	defer s.root.Remove(tmp)
	if _, err = f.Write(body); err == nil {
		err = f.Sync()
	}
	closeErr := f.Close()
	if err != nil {
		return storageError(err)
	}
	if closeErr != nil {
		return storageError(closeErr)
	}
	if err = s.root.Link(tmp, target); err != nil {
		return storageError(err)
	}
	return nil
}
