package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Set header
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		//pass connection to next handler
		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Use infomation logger to log IP Address of user, URL and method requested
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {

			if err := recover(); err != nil {
				//Set Connection in the header to close
				w.Header().Set("Connection", "close")

				//Call the app.serverError helper method to return internal server error response
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//if user is not login, redirect to login page
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		//Set header to avoid storage in the user browser's cache
		w.Header().Add("Cache-Control", "no-store")
		// Serve the next handler
		next.ServeHTTP(w, r)

	})
}

func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID") // Retrieve an autnenticated user id from the session,
		if id == 0 {                                                        // it returns zero if not found, serve the normal http in such case
			next.ServeHTTP(w, r)
			return
		}

		exists, err := app.user.Exists(id) // Check if user exist in the database
		if err != nil {
			app.serverError(w, err)
			return
		}

		if exists { // if the the user exists, set isAuthenticated to true and add the data to the request context
			ctx := r.Context()
			ctx = context.WithValue(ctx, isAuthenticatedContextKey, true)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}
