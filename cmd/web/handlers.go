package main

import (
    "errors"
    "fmt"
    "net/http"
    "strconv"

    "github.com/go-zoo/bone"
    "jackson.software/snippetbox/pkg/forms"
    "jackson.software/snippetbox/pkg/models"
)

// Home handler shows the home page with the latest snippets.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
    s, err := app.snippets.Latest()
    if err != nil {
        app.serverError(w, err)
        return
    }

    app.render(w, r, "home.page.tmpl", &templateData{Snippets: s})
}

// ShowSnippet handler shows a specific snippet.
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(bone.GetValue(r, "id"))
    if err != nil || id < 1 {
        app.notFound(w)
        return
    }

    s, err := app.snippets.Get(id)
    if err != nil {
        if errors.Is(err, models.ErrNoRecord) {
            app.notFound(w)
        } else {
            app.serverError(w, err)
        }
        return
    }

    app.render(w, r, "show.page.tmpl", &templateData{Snippet: s})
}

// ShowSnippetForm handler shows a form with fields to create a new snippet.
func (app *application) showSnippetForm(w http.ResponseWriter, r *http.Request) {
    app.render(w, r, "create.page.tmpl", &templateData{
        Form: forms.New(nil),
    })
}

// CreateSnippet handler creates a new snippet.
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        app.clientError(w, http.StatusBadRequest)
        return
    }

    form := forms.New(r.PostForm)
    form.Required("title", "content", "expires")
    form.MaxLength("title", 100)
    form.PermittedValues("expires", "365", "7", "1")

    if !form.Valid() {
        app.render(w, r, "create.page.tmpl", &templateData{Form: form})
        return
    }

    id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
    if err != nil {
        app.serverError(w, err)
        return
    }

    app.session.Put(r, "flash", "Your Snippet was successfully created")

    http.Redirect(w, r, fmt.Sprintf("/snippets/%d", id), http.StatusSeeOther)
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
    app.render(w, r, "signup.page.tmpl", &templateData{
        Form: forms.New(nil),
    })
}

// SignupUser signs the user up if the given data if the given data is valid,
// by creating a new account.
func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        app.clientError(w, http.StatusBadRequest)
        return
    }

    form := forms.New(r.PostForm)
    form.Required("name", "email", "password")
    form.MaxLength("name", 255)
    form.MaxLength("email", 255)
    form.MatchesPattern("email", forms.EmailRX)
    form.MinLength("password", 10)

    if !form.Valid() {
        app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
        return
    }

    err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
    if err != nil {
        if errors.Is(err, models.ErrDuplicateEmail) {
            form.Errors.Add("email", "Address is already in use")
            app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
        } else {
            app.serverError(w, err)
        }
        return
    }

    app.session.Put(r, "flash", "You signup was successful. Please log in")
    http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}

// LoginUserForm renders the login page for the user to login.
func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
    app.render(w, r, "login.page.tmpl", &templateData{
        Form: forms.New(nil),
    })
}

// LoginUser logs the user in by validating the given email and password.
// If the credentials are invalid, an error will be returned.
func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        app.clientError(w, http.StatusBadRequest)
        return
    }

    form := forms.New(r.PostForm)
    id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
    if err != nil {
        if errors.Is(err, models.ErrInvalidCredentials) {
            form.Errors.Add("generic", "Email or Password is incorrect")
            app.render(w, r, "login.page.tmpl", &templateData{Form: form})
        } else {
            app.serverError(w, err)
        }
        return
    }

    app.session.Put(r, "authenticatedUserID", id)
    http.Redirect(w, r, "/snippets/create", http.StatusSeeOther)
}

// LogoutUser logs the user out by deleting the authentication cookie.
func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
    app.session.Remove(r, "authenticatedUserID")
    app.session.Put(r, "flash", "You've been logged out successfully!")

    http.Redirect(w, r, "/", http.StatusSeeOther)
}
