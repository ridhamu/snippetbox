package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ridhamu/snippetbox/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "GO")
	// w.WriteHeader(http.StatusOK)
	// _, _ = fmt.Fprintf(w, "welcomd to home\n")

	snippetList, err := app.snippetModel.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	for _, snippet := range snippetList {
		fmt.Fprintf(w, "%+v\n", snippet)
	}

	// files := []string{
	// 	"./ui/html/base.html",
	// 	"./ui/html/pages/home.html",
	// 	"./ui/html/partials/nav.html",
	// }
	//
	// template, err := template.ParseFiles(files...)
	// if err != nil {
	// 	// log.Println(err.Error())
	// 	// app.logger.Error(err.Error(), "method", r.Method, "URI", r.URL.RequestURI())
	// 	// http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	app.serverError(w, r, err)
	// 	return
	// }
	//
	// // err = template.Execute(w, nil)
	// err = template.ExecuteTemplate(w, "base", nil)
	// if err != nil {
	// 	// log.Println(err.Error())
	// 	// app.logger.Error(err.Error(), "method", r.Method, "URI", r.URL.RequestURI())
	// 	// http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	app.serverError(w, r, err)
	// 	return
	// }
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	snippet, err := app.snippetModel.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	_, _ = fmt.Fprintf(w, "%+v", snippet)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintln(w, "Displaying form for creating post")
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := app.snippetModel.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
