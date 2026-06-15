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

type SnippetModel struct {
	DB *sql.DB
}

type SnippetModelInterface interface {
	Insert(title string, content string, expires int) (int, error)
	Get(id int) (Snippet, error)
	Latest() ([]Snippet, error)
}

func (sm *SnippetModel) Insert(title string, content string, expires int) (int, error) {
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

func (sm *SnippetModel) Get(id int) (Snippet, error) {
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

func (sm *SnippetModel) Latest() ([]Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	var snippetList []Snippet

	rows, err := sm.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var oneRowSnippet Snippet

		err = rows.Scan(&oneRowSnippet.ID, &oneRowSnippet.Title, &oneRowSnippet.Content, &oneRowSnippet.Created, &oneRowSnippet.Expires)
		if err != nil {
			return nil, err
		}

		snippetList = append(snippetList, oneRowSnippet)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippetList, nil
}
