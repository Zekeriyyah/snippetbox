package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// Initialize the router
	router := httprouter.New()

	//Setting custom 404 error handler to the router in other to use app custom Not Found helper function
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	//Register mux To handle static file
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	//Register the other application routes
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)

	// //Passing mux as the next parameter to the secureHeader middleware.
	// return app.recoverPanic((secureHeaders(mux)))

	//using justinas/alice to chain the middlewares
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
}
