package render

import (
	"bytes"
	"clipService/internal/config"
	"clipService/internal/models"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

var functions = template.FuncMap{
	"dateFromLayout": func(d string) string {
		layout := time.RFC3339

		t, err := time.Parse(layout, d)
		if err != nil {
			fmt.Println("Error parsing timestamp:", err)

		}
		fmt.Println("the_time2", t)
		return t.Format("2006-01-02 15:04") // apply format
	},
}

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func NewRenderer(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	return td
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template

	if app.UseCache {
		// get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)
	td = AddDefaultData(td, r)
	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser", err)
	}

}

func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	// find template by  using regez
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	// loop  html pages
	for _, page := range pages {
		name := filepath.Base(page)
		// create template
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		// find layout
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {

			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}
		// add template to cache
		myCache[name] = ts
	}

	return myCache, nil
}
