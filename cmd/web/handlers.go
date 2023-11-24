package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Zekeriyyah/snippetbox/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	//Check if the url path is strictly '/'
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	//panic("oops! something went wrong") //Deliberate panic to test my recoverPanic middleware use in stack tracing the error to server error log to return error 500 for better user experience

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	//Initialize templateDate with to store current year value
	data := app.newTemplateData(r)
	data.Snippets = snippets

	//Render the template using render helper function
	app.render(w, http.StatusOK, "home.tmpl.html", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}

		return
	}

	//Initialize new template data with current year stamp
	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.tmpl.html", data)
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
	// w.Write([]byte("Snippet is created"))

	title := "0 snail"
	content := "0 snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
	expires := 7

	//Passing the data to SnippetModel
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	//Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)

}
