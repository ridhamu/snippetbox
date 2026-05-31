// Package models is internal utilities for accessing various data source
package models

import (
	"database/sql"
	"errors"
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
	stmt := `INSERT INTO snippets (title, content, created, expires)
VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := sm.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (sm *SnippetModels) Get(id int) (Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
WHERE expires > UTC_TIMESTAMP() AND id = ?`

	var oneRowSnippet Snippet

	row := sm.DB.QueryRow(stmt, id)
	err := row.Scan(&oneRowSnippet.ID, &oneRowSnippet.Title, &oneRowSnippet.Content, &oneRowSnippet.Created, &oneRowSnippet.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}

	return oneRowSnippet, nil
}

func (sm *SnippetModels) Latest() ([]Snippet, error) {
	return nil, nil
}
