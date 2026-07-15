package handlers

import (
	"net/http"

	"rag/internal/render"
	"rag/internal/services"
)

type Handlers struct {
	uploadService *services.UploadService
	ragService    *services.RagService
	Renderer      *render.Renderer
}

func New(uploadService *services.UploadService, ragService *services.RagService, renderer *render.Renderer) *Handlers {
	return &Handlers{
		uploadService: uploadService,
		ragService:    ragService,
		Renderer:      renderer,
	}
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	if err := h.Renderer.Render(w, "home.html", nil); err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) GetAllUploadFiles(w http.ResponseWriter, r *http.Request) {
	files, err := h.uploadService.GetAllUploadedFile()
	if err != nil {
		http.Error(w, "unable to retrieve all uploaded files", http.StatusInternalServerError)
		return
	}

	if err := h.Renderer.Render(w, "all_uploaded_file.html", files); err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) UploadDoc(w http.ResponseWriter, r *http.Request) {
	if err := h.Renderer.Render(w, "upload.html", nil); err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) ProcessDoc(w http.ResponseWriter, r *http.Request) {
	// Limit request size to 10 MB.
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "file too large or invalid form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "missing uploaded file", http.StatusBadRequest)
		return
	}
	defer func() {
		_ = file.Close()
	}()

	description := r.FormValue("description")

	_, err = h.uploadService.SaveUploadedFile(r.Context(), file, header, description)
	if err != nil {
		http.Error(w, "unable to save uploaded file", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/docs", http.StatusSeeOther)
}
