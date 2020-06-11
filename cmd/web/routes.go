package main

import (
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	mux := bone.New()

	mux.Get("/", http.HandlerFunc(app.home)).Options()
	mux.Get("/snippets/create", http.HandlerFunc(app.showSnippetForm))
	mux.Get("/snippets/:id", http.HandlerFunc(app.showSnippet))
	mux.Post("/snippets", http.HandlerFunc(app.createSnippet))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
