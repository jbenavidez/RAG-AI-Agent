package routes

import (
	"net/http"
	"rag/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func SetUpReoutes(ragHandler *handlers.RagHandler) http.Handler {

	mux := chi.NewRouter()
	mux.Get("/docs", ragHandler.GetAllUploadFiles)
	mux.Get("/upload", ragHandler.UploadDoc)
	mux.Post("/upload", ragHandler.ProcessDoc)
	return mux
}
