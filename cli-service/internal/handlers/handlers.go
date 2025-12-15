package handlers

import (
	"clipService/internal/config"
	"clipService/internal/models"
	"clipService/internal/render"
	"net/http"
)

type Repository struct {
	App *config.AppConfig
}

var Repo *Repository

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

// ChatRoom is the handler for the chatroom page
func (m *Repository) ChatRoom(w http.ResponseWriter, r *http.Request) {

	// Render template
	render.RenderTemplate(w, r, "chatroom.page.tmpl", &models.TemplateData{
		Data: nil,
	})
}
