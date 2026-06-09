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

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /snippet/create", dynamic.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

	// user related routing
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userLoginPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))
	mux.Handle("POST /user/logout", dynamic.ThenFunc(app.userLoginPost))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
