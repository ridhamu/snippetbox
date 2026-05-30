package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "GO")
	// w.WriteHeader(http.StatusOK)
	// _, _ = fmt.Fprintf(w, "welcomd to home\n")

	files := []string{
		"./ui/html/base.html",
		"./ui/html/pages/home.html",
		"./ui/html/partials/nav.html",
	}

	template, err := template.ParseFiles(files...)
	if err != nil {
		// log.Println(err.Error())
		// app.logger.Error(err.Error(), "method", r.Method, "URI", r.URL.RequestURI())
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		app.serverError(w, r, err)
		return
	}

	// err = template.Execute(w, nil)
	err = template.ExecuteTemplate(w, "base", nil)
	if err != nil {
		// log.Println(err.Error())
		// app.logger.Error(err.Error(), "method", r.Method, "URI", r.URL.RequestURI())
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		app.serverError(w, r, err)
		return
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	_, _ = fmt.Fprintf(w, "Displaying snippetview with id %d\n", id)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintln(w, "Displaying form for creating post")
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	_, _ = fmt.Fprintln(w, "Creating new snippet...")
}
