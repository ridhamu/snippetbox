package main

import (
	"net/http"

	"github.com/justinas/alice"
)

// mux := http.NewServeMux()
// staticFileHandler := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})
//
// // mux.Handle("GET /static/", http.StripPrefix("/static", neuter(staticFileHandler)))
// mux.Handle("GET /static/", http.StripPrefix("/static", staticFileHandler))
//
// mux.HandleFunc("GET /{$}", app.home)
// mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
// mux.HandleFunc("GET /snippet/create", app.snippetCreate)
// mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	staticFileHandler := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})

	// mux.Handle("GET /static/", http.StripPrefix("/static", neuter(staticFileHandler)))
	mux.Handle("GET /static/", http.StripPrefix("/static", staticFileHandler))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	myMiddleware := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return myMiddleware.Then(mux)
}
