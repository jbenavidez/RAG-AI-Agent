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

func (h *RagHandler) GetAllUploadFiles(w http.ResponseWriter, r *http.Request) {

	//get all files
	files, err := h.service.GetAllfile()
	if err != nil {
		http.Error(w, "unable to retrieve all uploaded files", http.StatusBadRequest)
		return
	}
	fmt.Println("all the ifles,", files)
	h.Renderer.Render(w, "upload.html", nil)

}

func (h *RagHandler) UploadDoc(w http.ResponseWriter, r *http.Request) {

	h.Renderer.Render(w, "upload.html", nil)

}

func (h *RagHandler) ProcessDoc(w http.ResponseWriter, r *http.Request) {

	// limit max uplaod
	// Limit request size to 10 MB
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "file too large or invalid form", http.StatusBadRequest)
		return
	}
	//get file
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "missing uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	description := r.FormValue("description")
	//save file with pending status
	storedfile, err := h.service.SaveUploadedFile(r.Context(), file, header, description)
	if err != nil {
		fmt.Println("something break", err)
		http.Error(w, "unable to save uploaded file", http.StatusInternalServerError)
		return
	}
	// TODO:spin a gorutine to process
	fmt.Println("the file", storedfile)
	// render back
	http.Redirect(w, r, "/upload", http.StatusSeeOther)

}
