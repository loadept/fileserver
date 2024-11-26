package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/jsusmachaca/fileserver/api/handler"
	"github.com/jsusmachaca/fileserver/api/middleware"
	"github.com/jsusmachaca/fileserver/internal/database"
	"github.com/jsusmachaca/go-router/pkg/router"
)

var db *sql.DB

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("\033[31mNot .env file found. Using system variables\033[0m")
	}

	var err error
	db, err = database.GetConnection()
	if err != nil {
		log.Fatalf("Error to connect database: %v", err)
	}

	err = database.Migrate(db)
	if err != nil {
		log.Fatalf("Error to migrate database: %v", err)
	}
}

func main() {
	PATH_DIR := os.Getenv("PATH_DIR")
	PORT := os.Getenv("PORT")

	route := router.NewRouter()

	getFiles := &handler.GetFiles{PathDir: PATH_DIR}
	uploadFiles := &handler.UploadFiles{PathDir: PATH_DIR}
	listDir := &handler.ListDir{PathDir: PATH_DIR}
	login := &handler.Login{DB: db}
	register := &handler.Register{DB: db}

	route.Get("/api/fs/", getFiles, middleware.AuthMiddlewareQuery)
	route.Put("/api/fs/upload", uploadFiles, middleware.AuthMiddleware)
	route.Get("/api/list/", listDir, middleware.AuthMiddleware)
	route.Post("/api/login", login)
	route.Post("/api/register", register)

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: route.ServeMux,
	}
	fmt.Printf("Server listen on port %s\n", PORT)
	server.ListenAndServe()
}
