package handler

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jsusmachaca/go-router/pkg/response"
)

type GetFiles struct {
	PathDir string
}

func (h *GetFiles) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathDir := h.PathDir

	filename := r.URL.Path[len("/api/fs/"):]
	cleanedFilename := filepath.Clean(filename)
	file := filepath.Join(pathDir, cleanedFilename)

	fileStat, err := os.Stat(file)
	if os.IsNotExist(err) {
		log.Printf("Error at requested file for %s in %s: %s", r.RemoteAddr, cleanedFilename, err)
		response.JsonErrorFromString(w, "The file does not exist", http.StatusNotFound)
		return
	}
	if fileStat.IsDir() {
		log.Printf("Error, directory by %s was requested at %s: %s", r.RemoteAddr, cleanedFilename, err)
		http.Redirect(w, r, "/api/list/"+cleanedFilename, 301)
		return
	}

	fileStream, err := os.Open(file)
	if err != nil {
		log.Printf("Error opening file %s: %v", cleanedFilename, err)
		response.JsonErrorFromString(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer fileStream.Close()

	buffer := make([]byte, 512)
	_, err = fileStream.Read(buffer)
	if err != nil {
		log.Printf("Error reading file %s: %v", cleanedFilename, err)
		response.JsonErrorFromString(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	contentType := http.DetectContentType(buffer)

	if contentType == "application/octet-stream" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	} else {
		w.Header().Set("Content-Type", contentType)
	}

	log.Printf("resource requested for %s in %s, type %s", r.RemoteAddr, cleanedFilename, contentType)
	http.ServeFile(w, r, file)
}
