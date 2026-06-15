package main

import (
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/justinas/alice"
	"github.com/ridhamu/snippetbox/ui"
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

	// staticFileHandler := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})
	// mux.Handle("GET /static/", http.StripPrefix("/static", staticFileHandler))

	// mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	staticFiles, err := fs.Sub(ui.Files, "static")
	if err != nil {
		panic(err)
	}

	staticHandler := http.FileServer(neuteredFileSystem{http.FS(staticFiles)})
	mux.Handle("GET /static/", staticHandler)

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	// routes that doesn't required auth
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))
	mux.Handle("GET /ping", dynamic.ThenFunc(ping))

	requiredAuth := dynamic.Append(app.requireAuthentication)
	// routes that required auth
	mux.Handle("GET /snippet/create", requiredAuth.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", requiredAuth.ThenFunc(app.snippetCreatePost))
	mux.Handle("POST /user/logout", requiredAuth.ThenFunc(app.userLogout))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}
