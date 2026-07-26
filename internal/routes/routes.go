package routes

import (
	"net/http"
	"rag/internal/handlers"
	"rag/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func SetUpReoutes(h *handlers.Handlers) http.Handler {

	mux := chi.NewRouter()
	mux.Use(middleware.SessionMiddleware)
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Get("/", h.Home)
	mux.Get("/ws", h.WebSocket)
	mux.Get("/docs", h.GetAllUploadFiles)
	mux.Get("/upload", h.UploadDoc)
	mux.Post("/upload", h.ProcessDoc)
	return mux
}
