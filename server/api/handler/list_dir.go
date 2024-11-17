package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	dirs "github.com/jsusmachaca/fileserver/pkg/file_system"
	"github.com/jsusmachaca/go-router/pkg/response"
)

type ListDir struct {
	PathDir string
}

func (h *ListDir) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathDir := h.PathDir

	var FileSystem dirs.DirectoryStructure
	FileSystem.Dirs = []dirs.Content{}
	FileSystem.Files = []dirs.Content{}

	directories := r.URL.Path[6:]
	cleanedFilename := filepath.Clean(directories)
	fullPath := filepath.Join(pathDir, cleanedFilename)

	stat, err := os.Stat(fullPath)
	if err != nil {
		log.Printf("Error accessing path %s: %v", fullPath, err)
		response.JsonErrorFromString(w, "Directory not found", http.StatusNotFound)
		return
	}
	if !stat.IsDir() {
		log.Printf("Error, directory by %s was requested at %s: %s", r.RemoteAddr, cleanedFilename, err)
		http.Redirect(w, r, "/fs/"+cleanedFilename, http.StatusMovedPermanently)
		return
	}

	fs, err := os.ReadDir(fullPath)
	if err != nil {
		log.Printf("Error reading directory %s: %v", fullPath, err)
		response.JsonErrorFromString(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	baseURL := fmt.Sprintf("%s://%s", scheme, r.Host)

	for _, data := range fs {
		if data.Type().IsDir() {
			FileSystem.Dirs = append(
				FileSystem.Dirs,
				dirs.Content{
					URL:  fmt.Sprintf("%s/list/%s", baseURL, filepath.Join(cleanedFilename, data.Name())),
					Name: data.Name(),
				},
			)
			continue
		}
		FileSystem.Files = append(
			FileSystem.Files,
			dirs.Content{
				URL:  fmt.Sprintf("%s/fs/%s", baseURL, filepath.Join(cleanedFilename, data.Name())),
				Name: data.Name(),
			},
		)
	}
	response.JsonResponse(w, FileSystem, http.StatusAccepted)
}
