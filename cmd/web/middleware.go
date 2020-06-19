package main

import (
    "context"
    "errors"
    "fmt"
    "jackson.software/snippetbox/pkg/models"
    "net/http"

	"github.com/justinas/nosurf"
)

// SecureHeaders adds extra header to every request in order to add security.
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

// LogRequest logs the IP, protocol, request method and file of every request .
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

// RecoverPanic checks if a panic occurred and if so, it will close the current
// connection and log the panic's argument and respond with a server error
// instead of just closing the request without providing an answer.
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// RequireAuthentication checks if the request comes from a user who is logged in.
// If the request comes from a non-logged in user, the user will be redirected
// to the login page.
func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/users/login", http.StatusSeeOther)
			return
		}

		// Do not cache pages that require authentication
		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

// CSRF middleware which checks for a valid CSRF token in order
// for a form submit to be handled correctly.
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}

// Checks the request is coming from an authenticated and active user.
// If the user is authenticated and active, a new copy of the request
// will be created and a key will be added to the request context,
// confirming that the request is coming from a valid user.
func (app *application) authenticate(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        exists := app.session.Exists(r, "authenticatedUserID")
        if !exists {
            next.ServeHTTP(w, r)
            return
        }

        // user is not valid (might have been removed etc.)
        user, err := app.users.Get(app.session.GetInt(r, "authenticatedUserID"))
        if errors.Is(err, models.ErrNoRecord) || !user.Active {
            app.session.Remove(r, "authenticatedUserID")
            next.ServeHTTP(w, r)
            return
        } else if err != nil {
            app.serverError(w, err)
            return
        }

        ctx := context.WithValue(r.Context(), contextKeyIsAuthenticated, true)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
