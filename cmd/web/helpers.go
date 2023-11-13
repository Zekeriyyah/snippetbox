package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// Create a serverError helper to write error message and stack trace to errorLog
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Create a clientError helper sends a specific code and a corresponding status code
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// Create a notFound helper to sends a 404 not found error to client
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
