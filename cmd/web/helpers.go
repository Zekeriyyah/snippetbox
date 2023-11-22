package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
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

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {

	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	//Initialize a new buffer
	buf := new(bytes.Buffer)

	//Execute the template set and write the response body.
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
	}

	//write out the HTTP status code
	w.WriteHeader(status)

	//Write to http.ResponseWriter from the buf
	buf.WriteTo(w)

}

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}
