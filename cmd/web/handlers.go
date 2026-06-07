package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ridhamu/snippetbox/internal/models"
	"github.com/ridhamu/snippetbox/internal/validator"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippetList, err := app.snippetModel.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.SnippetList = snippetList

	app.render(w, r, http.StatusOK, data, "home.html")
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

	// templateData := templateData{
	// 	Snippet: snippet,
	// }

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, r, http.StatusOK, data, "view.html")
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	formData := FormData{
		Expires: 365,
	}

	data := app.newTemplateData(r)
	data.Form = formData

	app.render(w, r, http.StatusOK, data, "create.html")
}

type FormData struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	var formData FormData

	err := app.decodePostForm(r, &formData)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// not blank title
	formData.CheckField(validator.NotBlankString(formData.Title), "title", "This field cannot be blank")
	// not more than 100 chars
	formData.CheckField(validator.MaxChars(formData.Title, 100), "title", "This field cannot be more than 100 characters long")
	// not blank content
	formData.CheckField(validator.NotBlankString(formData.Content), "content", "This field cannot be blank")
	// permitted value in expires
	formData.CheckField(validator.PermittedValue(formData.Expires, 1, 7, 365), "expires", "This field must equal 1, 7, 365")

	if !formData.Valid() { // i want to add break point here and see the errorFields
		data := app.newTemplateData(r)
		data.Form = formData
		app.render(w, r, http.StatusUnprocessableEntity, data, "create.html")
		return
	}

	id, err := app.snippetModel.Insert(formData.Title, formData.Content, formData.Expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
