package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Zekeriyyah/snippetbox/internal/models"
	"github.com/Zekeriyyah/snippetbox/internal/validator"
	"github.com/julienschmidt/httprouter"
)

type SnippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

// struct to hold signup form data
type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

// struct to hold user login data
type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {

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

	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
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

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// Declare a new empty instance of the snippetCreateForm struct.
	var form SnippetCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Using CheckFied() to write appropriate error message after validating the form
	form.CheckField(validator.NotBlank(form.Title), "title", "*This field cannot be blank!")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "*This field cannot be morethan 100 characters")
	form.CheckField(validator.NotBlank(form.Content), "content", "*This field cannot be blank!")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "*This field must not take value other than 1, 7 or 365")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	//Passing the data to SnippetModel
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Add a flash message to the session data on successful creation of snippet
	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	//Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	//Initialize the Form field in the templateData and you can set default Expires value
	data.Form = SnippetCreateForm{
		Expires: 1,
	}

	app.render(w, http.StatusOK, "create.tmpl.html", data)
}

//New handlers to handle usser authentication

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, http.StatusOK, "signup.tmpl.html", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	//Create Instance of userSignupForm
	var form userSignupForm

	//Parse the form into the userSignupForm struct
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	//Validate the form content
	form.CheckField(validator.NotBlank(form.Name), "name", "**This field cannot be blank!")
	form.CheckField(validator.NotBlank(form.Email), "email", "**This field cannot be blank!")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "**This field must be a valid email address!")
	form.CheckField(validator.NotBlank(form.Password), "password", "**This field cannot be blank!")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "**This field must be at least 8 characters long!")

	// Re-render form with error msg if any
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		return
	}

	//Add user to the database or return error msg if email already exist
	err = app.user.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "**email address is already in use!")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		} else {
			app.serverError(w, err)
		}

		return
	}

	//Otherwise, send a placeholder response to confirm successful signup
	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful, please log in.")

	//Redirect user to login page
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.tmpl.html", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	//Decode form data into userLoginForm struct
	var form userLoginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	//Perform validation check on the email and password
	form.CheckField(validator.NotBlank(form.Email), "email", "**This field cannot be blank!")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "**This field must be a valid email")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	//Re-render the login form with error msg if any
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		return
	}

	//Authenticate user for login
	id, err := app.user.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("**email or password is incorrect!")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.tmpl.html", data)
			return
		} else {
			app.serverError(w, err)
			return
		}
		return
	}

	//Renew the session ID for fresh login of user
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	//Add the already authenticated user's ID to the session data
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	//Redirect the user to the create snippet page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	//Renew user Token
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Remove authentication ID of user
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	// Deem user of successful logout
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	//Redirect user to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
