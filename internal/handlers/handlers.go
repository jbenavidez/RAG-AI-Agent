package handlers

import (
	"fmt"
	"net/http"
	"rag/internal/render"
	"rag/internal/services"
)

type RagHandler struct {
	service  *services.UploadService
	Renderer *render.Renderer
}

func NewRagHandler(s *services.UploadService, r *render.Renderer) *RagHandler {
	return &RagHandler{
		service:  s,
		Renderer: r,
	}

}

func (h *RagHandler) UploadDoc(w http.ResponseWriter, r *http.Request) {

	//render upload html tempalte

	fmt.Println("hello there")
	h.Renderer.Render(w, "upload.html", nil)

}
