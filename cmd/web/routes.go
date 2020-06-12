package main

import (
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable)

	mux := bone.New()

	mux.Get("/", dynamicMiddleware.ThenFunc(app.home)).Options()
	mux.Get("/snippets/create", dynamicMiddleware.ThenFunc(app.showSnippetForm))
	mux.Get("/snippets/:id", dynamicMiddleware.ThenFunc(app.showSnippet))
	mux.Post("/snippets", dynamicMiddleware.ThenFunc(app.createSnippet))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
