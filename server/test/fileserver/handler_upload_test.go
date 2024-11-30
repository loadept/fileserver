package test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/jsusmachaca/fileserver/api/handler"
)

func TestUploadFileHandle(t *testing.T) {
	tmpDir := t.TempDir()
	err := os.Mkdir(filepath.Join(tmpDir, "subdir"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	part, err := writer.CreateFormFile("file", "testfile.txt")
	if err != nil {
		t.Fatal(err)
	}

	_, err = part.Write([]byte("Testing content"))
	if err != nil {
		t.Fatal(err)
	}

	err = writer.Close()
	if err != nil {
		t.Fatal(err)
	}

	content := buf.Bytes()

	h := handler.UploadFiles{PathDir: tmpDir}

	t.Run("Valid Upload File", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/api/fs/upload", bytes.NewBuffer(content))
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Host = "localhost:8080"
		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expected := "{\"message\":\"File uploaded testfile.txt\"}\n"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %s want %s", rr.Body.String(), expected)
		}

		filePath := filepath.Join(tmpDir, "testfile.txt")
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("file not found: %v", filePath)
		}
	})

	t.Run("Valid Upload File In Directory", func(t *testing.T) {
		req, err := http.NewRequest("PUT", "/api/fs/upload?directory=subdir", bytes.NewBuffer(content))
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Host = "localhost:8080"
		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expected := "{\"message\":\"File uploaded testfile.txt\"}\n"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %s want %s", rr.Body.String(), expected)
		}

		filePath := filepath.Join(tmpDir, "subdir", "testfile.txt")
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("file not found: %v", filePath)
		}
	})
}
