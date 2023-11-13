package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	//Check if the url path is strictly '/'
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Display a specific snippet with ID = %d\n", id)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	//Check for method if POST
	w.Header()["Content-Type"] = nil
	if r.Method != "POST" {

		w.Header().Set("Allow", "POST")

		// w.WriteHeader(405)
		// w.Write([]byte("Method Not Allowed"))

		app.clientError(w, 405) // Alternative way of writing error response
		return
	}
	w.Write([]byte("Snippet is created"))
}
