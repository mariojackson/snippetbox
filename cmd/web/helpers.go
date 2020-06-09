package main

import (
    "fmt"
    "net/http"
    "runtime/debug"
)

// Server error writes an error message and stack trace to the error log
// and then sends a generic 500 server error with the according description.
func (app *application) serverError(w http.ResponseWriter, err error) {
    trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
    // Instead of logging this location as the last one, use the second last one
    // which is where the error occurred
    app.errorLog.Output(2, trace)

    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Client error sends a specific status code and the corresponding description
// to the user.
func (app *application) clientError(w http.ResponseWriter, status int) {
    http.Error(w, http.StatusText(status), status)
}

// Not found, a convenient wrapper around the client error in order
// to send a 404 not found error to the user.
func (app *application) notFound(w http.ResponseWriter) {
   app.clientError(w, http.StatusNotFound)
}
