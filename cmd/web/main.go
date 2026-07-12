package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"rag/internal/config"
	"rag/internal/db"
	"rag/internal/handlers"
	"rag/internal/render"
	dbrepo "rag/internal/repository/db_repo"
	"rag/internal/routes"
	"rag/internal/services"
	"rag/internal/storage"
)

const (
	portNumber = "8080"
)

func main() {

	//init project
	appConfig := config.AppConfig{
		UseCache: true,
	}
	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
	}
	appConfig.TemplateCache = templateCache
	fmt.Println("starting application on port", portNumber)

	weaviateClient, err := db.NewWeaviateClient()
	if err != nil {
		panic(err)
	}

	//wire everything up
	weaviateRepo := dbrepo.NewWeaviateDBRepo(weaviateClient)
	uploadDir := os.Getenv("UPLOAD_DIR")
	fileStorage := storage.NewLocalStorage(uploadDir)
	uploadService := services.NewUploadService(weaviateRepo, fileStorage)
	renderer := render.NewRenderer(&appConfig)
	ragHandlers := handlers.NewRagHandler(uploadService, renderer)
	mux := routes.SetUpReoutes(ragHandlers)

	// init server
	err = http.ListenAndServe(fmt.Sprintf(":%s", portNumber), mux)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("server up and running")
}
