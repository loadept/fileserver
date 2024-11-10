package main

import (
	"fmt"
	"net/http"

	"github.com/jsusmachaca/fileserver/api/handler"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /fs/", handler.GetFiles)
	mux.HandleFunc("PUT /fs/upload", handler.UploadFiles)
	mux.HandleFunc("GET /list/", handler.ListDir)

	server := http.Server{
		Addr:    ":8082",
		Handler: mux,
	}
	fmt.Println("Server listen on port 8082")
	server.ListenAndServe()
}
