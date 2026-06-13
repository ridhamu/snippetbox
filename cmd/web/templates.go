package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/ridhamu/snippetbox/internal/models"
	"github.com/ridhamu/snippetbox/ui"
)

type templateData struct {
	CurrentYear     int
	Snippet         models.Snippet
	SnippetList     []models.Snippet
	Form            any
	FlashMessage    string
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var funcs = template.FuncMap{
	"HumanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	fileList, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, file := range fileList {
		name := filepath.Base(file)

		patterns := []string{
			"html/base.html",
			"html/partials/*.html",
			file,
		}

		ts, err := template.New(name).Funcs(funcs).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
