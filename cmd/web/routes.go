package main

import (
	"net/http"

	"github.com/Zekeriyyah/snippetbox/ui"
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

	/*	====== NO LONGER SERVING STATIC FILES FROM THE DISK

		//Register mux To handle static file
		fileServer := http.FileServer(http.Dir("./ui/static/"))
		router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))
	*/

	//SERVING STATIC FILES FROM EMBEDDED FILE SYSTEM IN GLOBAL VARIABLE ui.Files
	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	router.HandlerFunc(http.MethodGet, "/ping", ping)

	// Create a dynamic middleware chain for session manager to wrap the routes except /static/*filepath
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	// middleware for some protected routes such as Get /snippet/create, Post /snippet/create, and Post /user/logout

	protected := dynamic.Append(app.requireAuthentication)

	//Register the other application routes
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))

	//Add five routes for authentication of user
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	//using justinas/alice to chain the middlewares
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
}
