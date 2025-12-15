package main

import (
	"clipService/internal/config"
	"clipService/internal/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func routes(app *config.AppConfig) http.Handler {

	mux := chi.NewRouter()
	// display chat-room
	mux.Get("/", handlers.Repo.ChatRoom)
	return mux
}
