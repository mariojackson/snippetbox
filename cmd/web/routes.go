package main

import (
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate)

	mux := bone.New()

	mux.Get("/", dynamicMiddleware.ThenFunc(app.home)).Options()
	mux.Get("/snippets/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.showSnippetForm))
	mux.Get("/snippets/:id", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.showSnippet))
	mux.Post("/snippets", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippet))

	mux.Get("/users/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/users/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/users/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/users/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/users/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
