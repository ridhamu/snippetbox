// Package models is internal utilities for accessing various data source
package models

import (
	"database/sql"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModels struct {
	DB *sql.DB
}

func (sm *SnippetModels) Insert(title string, content string, expires int) (int, error) {
	return 1, nil
}

func (sm *SnippetModels) Get(id int) (Snippet, error) {
	return Snippet{}, nil
}

func (sm *SnippetModels) Latest() ([]Snippet, error) {
	return nil, nil
}
