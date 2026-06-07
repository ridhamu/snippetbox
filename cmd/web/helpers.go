package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
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

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, data templateData, tsName string) {
	t, ok := app.templateCache[tsName]
	if !ok {
		app.serverError(w, r, errors.New("no matching template for given path"))
		return
	}
	buf := new(bytes.Buffer)

	err := t.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)

	_, _ = buf.WriteTo(w)
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrentYear: time.Now().Year(),
	}
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		} else {
			return err
		}
	}

	return nil
}
