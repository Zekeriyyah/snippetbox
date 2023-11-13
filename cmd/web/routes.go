package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	//Register mux To handle static file
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	//Register the other application routes
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	return mux
}
