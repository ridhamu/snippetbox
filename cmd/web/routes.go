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

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	// routes that doesn't required auth
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	requiredAuth := dynamic.Append(app.requireAuthentication)
	// routes that required auth
	mux.Handle("GET /snippet/create", requiredAuth.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", requiredAuth.ThenFunc(app.snippetCreatePost))
	mux.Handle("POST /user/logout", requiredAuth.ThenFunc(app.userLogout))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
