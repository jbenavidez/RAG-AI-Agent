package main

import (
	"fmt"
	"log"
	"net/http"
	"rag/internal/config"
	"rag/internal/db"
	"rag/internal/handlers"
	"rag/internal/render"
	dbrepo "rag/internal/repository/db_repo"
	"rag/internal/routes"
	"rag/internal/services"
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

	conn, err := db.NewWeaviateClient()
	if err != nil {
		panic(err)
	}

	//wire everything up
	uploadService := services.NewUploadService(&dbrepo.WeaviateDBRepo{DB: conn})
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
