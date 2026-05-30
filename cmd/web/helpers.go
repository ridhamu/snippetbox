package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	// app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace) // since it will printout the \t and \n instead of enter of tab, i want to manipulate it so it able to format using tab or enter
	app.logger.Error(err.Error(), "method", method, "uri", uri) // since it will printout the \t and \n instead of enter of tab, i want to manipulate it so it able to format using tab or enter
	fmt.Printf("%s\n", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
