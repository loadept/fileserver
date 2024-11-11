package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jsusmachaca/fileserver/pkg/dirs"
)

func GetFiles(w http.ResponseWriter, r *http.Request) {
	pathDir := os.Getenv("PATH_DIR")

	filename := r.URL.Path[len("/fs/"):]
	cleanedFilename := filepath.Clean(filename)
	file := filepath.Join(pathDir, cleanedFilename)

	fileStat, err := os.Stat(file)
	if os.IsNotExist(err) {
		log.Printf("Error at requested file for %s in %s: %s", r.RemoteAddr, cleanedFilename, err)
		http.NotFound(w, r)
		return
	}
	if fileStat.IsDir() {
		log.Printf("Error, directory by %s was requested at %s: %s", r.RemoteAddr, cleanedFilename, err)
		http.Redirect(w, r, "/list/"+cleanedFilename, 301)
		return
	}

	fileStream, err := os.Open(file)
	if err != nil {
		log.Printf("Error opening file %s: %v", cleanedFilename, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer fileStream.Close()

	buffer := make([]byte, 512)
	_, err = fileStream.Read(buffer)
	if err != nil {
		log.Printf("Error reading file %s: %v", cleanedFilename, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
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

func UploadFiles(w http.ResponseWriter, r *http.Request) {
	pathDir := os.Getenv("PATH_DIR")
	directory := ""
	if dirs, ok := r.URL.Query()["directory"]; ok && len(dirs) > 0 {
		directory = dirs[0]
	}

	srcFile, header, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "file don't upload"}`))
		return
	}
	defer srcFile.Close()

	var uploadPath string
	if len(directory) > 0 {
		uploadPath = filepath.Join(pathDir, directory, header.Filename)
	} else {
		uploadPath = filepath.Join(pathDir, header.Filename)
	}

	dstFile, err := os.Create(uploadPath)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "file don't upload"}`))
		return
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "file don't upload"}`))
		return
	}
	log.Printf("File uploaded for %s", r.RemoteAddr)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).
		Encode(map[string]string{
			"message": "file uploaded " + header.Filename,
		}); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func ListDir(w http.ResponseWriter, r *http.Request) {
	var FileSystem dirs.DirectoryStructure
	FileSystem.Dirs = []string{}
	FileSystem.Files = []string{}

	pathDir := os.Getenv("PATH_DIR")
	directories := r.URL.Path[6:]
	cleanedFilename := filepath.Clean(directories)
	fullPath := filepath.Join(pathDir, cleanedFilename)

	stat, err := os.Stat(fullPath)
	if err != nil {
		log.Printf("Error accessing path %s: %v", fullPath, err)
		http.Error(w, "Directory not found", http.StatusNotFound)
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
		http.Error(w, "Failed to read directory", http.StatusInternalServerError)
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
				fmt.Sprintf("%s/list/%s", baseURL, filepath.Join(cleanedFilename, data.Name())),
			)
			continue
		}
		FileSystem.Files = append(
			FileSystem.Files,
			fmt.Sprintf("%s/fs/%s", baseURL, filepath.Join(cleanedFilename, data.Name())),
		)
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(FileSystem); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "failed to list directory"}`))
	}
}
