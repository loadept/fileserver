package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/jsusmachaca/fileserver/api/handler"
	"github.com/jsusmachaca/fileserver/api/middleware"
	"github.com/jsusmachaca/go-router/pkg/router"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("\033[31mNot .env file found. Using system variables\033[0m")
	}
}

func main() {
	PATH_DIR := os.Getenv("PATH_DIR")
	PORT := os.Getenv("PORT")

	route := router.NewRouter()

	getFiles := &handler.GetFiles{PathDir: PATH_DIR}
	uploadFiles := &handler.UploadFiles{PathDir: PATH_DIR}
	listDir := &handler.ListDir{PathDir: PATH_DIR}

	route.Get("/fs/", middleware.AuthMiddleware, getFiles)
	route.Put("/fs/upload/", middleware.AuthMiddleware, uploadFiles)
	route.Get("/list/", middleware.AuthMiddleware, listDir)

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: route.ServeMux,
	}
	fmt.Printf("Server listen on port %s\n", PORT)
	server.ListenAndServe()
}
