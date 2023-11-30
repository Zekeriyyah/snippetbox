package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/Zekeriyyah/snippetbox/internal/models"
)

type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	Form        any
}

func newTemplateCache() (map[string]*template.Template, error) {
	//Initialize a map to act as cache
	cache := map[string]*template.Template{}

	//Getting the all the file path for the page
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	//Loop through the pages and parse the file path
	for _, page := range pages {
		//Store the main path name i.e homt.tmpl.html in var name
		name := filepath.Base(page)

		//Create a slice to hold the static files to be parse

		// //To parse the base template
		// ts, err := template.ParseFiles("./ui/html/base.tmpl.html")
		// if err != nil {
		// 	return nil, err
		// }

		//To make use of custom template function humanDate, register the FuncMap with the template set and call ParseFile
		//to parse the templates

		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		//To parse all template files in the partials directory
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
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
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}
