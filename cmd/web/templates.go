package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/Zekeriyyah/snippetbox/internal/models"
	"github.com/Zekeriyyah/snippetbox/ui"
)

type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any
	Flash           string //to store flash msg for rendering
	IsAuthenticated bool   //To check if a session has AuthenticatedUserID field
	CSRFToken       string //To store the csrf token
}

func newTemplateCache() (map[string]*template.Template, error) {
	//Initialize a map to act as cache
	cache := map[string]*template.Template{}

	// Parsing file from the embedded files
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	//Loop through the pages and parse the file path
	for _, page := range pages {
		//Store the main path name i.e home.tmpl.html in var name
		name := filepath.Base(page)

		//Create a slice to hold the static files to be parse
		patterns := []string{
			"html/base.tmpl.html",
			"html/partials/*.tmpl.html",
			page,
		}

		//parse the templates using ParseFS instead of ParseFiles since from embedded files

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		//Add the template set to the map, using the name of the page
		cache[name] = ts
	}

	//Return the map
	return cache, nil

}

func humanDate(t time.Time) string {

	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}
