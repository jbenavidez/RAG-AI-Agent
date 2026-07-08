package routes

import (
	"net/http"
	"rag/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func SetUpReoutes(ragHandler *handlers.RagHandler) http.Handler {

	mux := chi.NewRouter()

	mux.Get("/", ragHandler.UploadDoc)

	return mux
}
