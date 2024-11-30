package test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/jsusmachaca/fileserver/api/handler"
	"github.com/jsusmachaca/go-router/pkg/router"
)

func TestFileHandle(t *testing.T) {
	tmpDir := t.TempDir()
	os.Mkdir(filepath.Join(tmpDir, "subdir"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("test file"), 0644)

	h := handler.ListDir{PathDir: tmpDir}

	t.Run("Valid Directory", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/list/", nil)
		req.Host = "localhost:8080"
		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expected :=
			"{\"dirs\":[{\"url\":\"http://localhost:8080/api/list/subdir\",\"name\":\"subdir\"}],\"files\":[{\"url\":\"http://localhost:8080/api/fs/file1.txt\",\"name\":\"file1.txt\"}]}\n"
		body := rr.Body.String()
		if body != expected {
			t.Errorf("handler returned unexpected body: got %s want %s", rr.Body.String(), expected)
		}
	})

	t.Run("Directory Traversal", func(t *testing.T) {
		hl := &handler.ListDir{PathDir: tmpDir}

		route := router.NewRouter()
		route.Get("/api/list/", hl)

		req, _ := http.NewRequest("GET", "/api/list/../../../../../../../../../../../../../../etc", nil)
		req.Host = "localhost:8080"
		rr := httptest.NewRecorder()

		route.ServeMux.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusMovedPermanently {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMovedPermanently)
		}
	})

	t.Run("Invalid Directory", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/list/nonexistent", nil)
		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})

	t.Run("File Instead of Directory", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/list/file1.txt", nil)
		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusMovedPermanently {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMovedPermanently)
		}
	})
}
