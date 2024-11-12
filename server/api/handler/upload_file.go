package handler

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jsusmachaca/go-router/pkg/response"
)

type UploadFiles struct {
	PathDir string
}

func (h *UploadFiles) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathDir := h.PathDir

	directory := ""
	if dirs, ok := r.URL.Query()["directory"]; ok && len(dirs) > 0 {
		directory = dirs[0]
	}

	srcFile, header, err := r.FormFile("file")
	if err != nil {
		response.JsonErrorFromString(w, "Bad request", http.StatusBadRequest)
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
		response.JsonErrorFromString(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		response.JsonErrorFromString(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	log.Printf("File uploaded for %s", r.RemoteAddr)

	response.JsonResponse(w, map[string]string{
		"message": "file uploaded " + header.Filename,
	}, http.StatusAccepted)
}
