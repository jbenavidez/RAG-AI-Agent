package render

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"rag/internal/config"
)

type Renderer struct {
	App *config.AppConfig
}

func New(app *config.AppConfig) *Renderer {
	return &Renderer{
		App: app,
	}
}

func (r *Renderer) Render(w http.ResponseWriter, name string, data any) error {
	t, ok := r.App.TemplateCache[name]
	if !ok {
		return fmt.Errorf("template %q not found", name)
	}

	return t.ExecuteTemplate(w, "base", data)
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.ParseFiles(
			"./templates/base.html",
			page,
		)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
		fmt.Println("loaded:", name)
	}

	return cache, nil
}
