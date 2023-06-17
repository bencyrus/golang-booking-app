package render

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/bencyrus/golang-booking-app/pkg/config"
	"github.com/bencyrus/golang-booking-app/pkg/models"
)

var functions = template.FuncMap{}

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

// RenderTemplate renders templates using html/template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	templateCache := map[string]*template.Template{}
	if app.UseCache {
		// Get the template cache from the app config
		templateCache = app.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()
	}

	// Get requested template from cache
	t, ok := templateCache[tmpl]
	if !ok {
		log.Println("error getting template from cache")
		return
	}
	// Render the template
	err := t.Execute(w, td)
	if err != nil {
		log.Println("error executing template:", err)
		return
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// Get all the files named *.page.html from the templates directory
	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		return myCache, err
	}

	// Loop through the pages one-by-one
	for _, page := range pages {
		// Extract the file name (like about.page.html) from the full path
		name := filepath.Base(page)
		templateSet, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		// Add any layout files to the template set
		matches, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			templateSet, err = templateSet.ParseGlob("./templates/*.layout.html")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = templateSet
	}

	return myCache, nil
}
