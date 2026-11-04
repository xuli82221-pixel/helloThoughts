package thoughts

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestRefuseWrongInstanceAndRedirect(t *testing.T) {
	loc := Location{Dir: t.TempDir(), Vault: t.TempDir()}
	for _, mode := range []string{"wrong_instance", "redirect"} {
		t.Run(mode, func(t *testing.T) {
			var reached bool
			target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { reached = true; w.WriteHeader(200) }))
			defer target.Close()
			fake := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if mode == "redirect" {
					http.Redirect(w, r, target.URL, 302)
					return
				}
				json.NewEncoder(w).Encode(State{Product: "memo", Version: Version, Vault: loc.Vault, Instance: "wrong"})
			}))
			defer fake.Close()
			s := State{Product: "memo", Version: Version, Vault: loc.Vault, Instance: randomID(), Token: randomID(), URL: fake.URL}
			b, _ := json.Marshal(s)
			if e := os.WriteFile(filepath.Join(loc.Dir, "runtime.json"), b, 0600); e != nil {
				t.Fatal(e)
			}
			if _, e := authenticate(context.Background(), loc); e == nil {
				t.Fatal("unverified service accepted")
			}
			if reached {
				t.Fatal("followed untrusted redirect")
			}
		})
	}
}
