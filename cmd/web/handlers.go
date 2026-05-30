package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
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
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// err = template.Execute(w, nil)
	err = template.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	_, _ = fmt.Fprintf(w, "Displaying snippetview with id %d\n", id)
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintln(w, "Displaying form for creating post")
}

func snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	_, _ = fmt.Fprintln(w, "Creating new snippet...")
}
