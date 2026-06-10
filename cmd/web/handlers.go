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

	app.sessionManager.Put(r.Context(), "flash", "Snippet succesfully created")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

type UserSignupFormData struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = UserSignupFormData{}
	app.render(w, r, http.StatusOK, data, "signup.html")
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var formData UserSignupFormData

	err := app.decodePostForm(r, &formData)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// not blank name
	formData.CheckField(validator.NotBlankString(formData.Name), "name", "This field cannot be blank")
	// not blank email
	formData.CheckField(validator.NotBlankString(formData.Email), "email", "This field cannot be blank")
	//	valid email
	formData.CheckField(validator.Matches(formData.Email, validator.EmailRX), "email", "This field must be a valid email address")
	// password cannot be empty
	formData.CheckField(validator.NotBlankString(formData.Password), "password", "This field cannot be blank")
	// password minimum 8 characters length
	formData.CheckField(validator.MinChars(formData.Password, 8), "password", "This field must be at least 8 characters long")

	if !formData.Valid() { // i want to add break point here and see the errorFields
		data := app.newTemplateData(r)
		data.Form = formData
		app.render(w, r, http.StatusUnprocessableEntity, data, "signup.html")
		return
	}

	err = app.userModel.Insert(formData.Name, formData.Email, formData.Password)
	if err != nil {
		if errors.Is(err, models.ErrDupEmail) {
			formData.AddField("email", "Email already taken")
			data := app.newTemplateData(r)
			data.Form = formData
			app.render(w, r, http.StatusUnprocessableEntity, data, "signup.html")
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your Signup was succesful. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Displaying form for user signup")
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Logging in the user")
}

func (app *application) userLogout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Logging out the user")
}
