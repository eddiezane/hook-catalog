package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/eddiezane/hook/pkg/hook"
)

type NopHandler struct{}

func (NopHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestHook(t *testing.T) {
	srv := httptest.NewServer(NopHandler{})
	defer srv.Close()

	if err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		// Exclude all hidden directories.
		if path != "." && info.IsDir() && strings.HasPrefix(path, ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
		}

		if !strings.HasSuffix(path, ".yaml") {
			return nil
		}

		t.Run(path, func(t *testing.T) {
			hooks, err := hook.NewFromPath(path)
			if err != nil {
				t.Fatal(err)
			}
			for _, h := range hooks {
				if _, err := h.Fire(srv.URL); err != nil {
					t.Error(err)
				}
			}
		})
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}
