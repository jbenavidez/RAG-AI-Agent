package routes

import (
	"net/http"
	"rag/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func SetUpReoutes(h *handlers.Handlers) http.Handler {

	mux := chi.NewRouter()

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Get("/", h.Home)
	mux.Get("/docs", h.GetAllUploadFiles)
	mux.Get("/docs", h.GetAllUploadFiles)
	mux.Get("/upload", h.UploadDoc)
	mux.Post("/upload", h.ProcessDoc)
	return mux
}
