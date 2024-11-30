package test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/jsusmachaca/fileserver/api/handler"
)

func TestGetFileHandle(t *testing.T) {
	tmpDir := t.TempDir()
	os.Mkdir(filepath.Join(tmpDir, "subdir"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("test file"), 0644)

	h := handler.GetFiles{PathDir: tmpDir}

	t.Run("Valid File", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/fs/file1.txt", nil)
		req.Host = "localhost:8080"
		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		contentExpected := "text/plain; charset=utf-8"
		if contentType := rr.Header().Get("Content-Type"); contentType != contentExpected {
			t.Errorf("handler returned wrong status code: got %v want %v", contentType, contentExpected)
		}
	})

	t.Run("Invalid File", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/fs/nonexistent", nil)
		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}

		expected := "{\"error\":\"The file does not exist\"}\n"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %s want %s", rr.Body.String(), expected)
		}
	})

	t.Run("Directory Instead of File", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/fs/subdir", nil)
		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusMovedPermanently {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMovedPermanently)
		}
	})
}
