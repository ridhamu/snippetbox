package main

import "github.com/ridhamu/snippetbox/internal/models"

type templateData struct {
	Snippet     models.Snippet
	SnippetList []models.Snippet
}
