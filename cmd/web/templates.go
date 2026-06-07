package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/ridhamu/snippetbox/internal/models"
)

type templateData struct {
	CurrentYear  int
	Snippet      models.Snippet
	SnippetList  []models.Snippet
	Form         any
	FlashMessage string
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var funcs = template.FuncMap{
	"HumanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	fileList, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range fileList {
		name := filepath.Base(page)

		// ts, err := template.ParseFiles("./ui/html/base.html")
		ts, err := template.New(name).Funcs(funcs).ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
