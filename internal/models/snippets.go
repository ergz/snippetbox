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

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	statement := `INSERT INTO snippets (title, content, created, expires) VALUES (?, ?, CURRENT_TIMESTAMP, datetime(CURRENT_TIMESTAMP, '+' || ? || ' day'))`
	result, err := m.DB.Exec(statement, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (Snippet, error) {
	statment := `SELECT id, title, content, created, expires FROM snippets WHERE expires > CURRENT_TIMESTAMP AND id = ?`
	row := m.DB.QueryRow(statment, id)
	var snip Snippet

	err := row.Scan(&snip.ID, &snip.Title, &snip.Content, &snip.Created, &snip.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}

	return snip, nil
}

func (m *SnippetModel) Latest(n int) ([]Snippet, error) {
	statement := `SELECT * FROM snippets where expires > CURRENT_TIMESTAMP ORDER BY id DESC LIMIT ?`
	rows, err := m.DB.Query(statement, n)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var snippets []Snippet

	for rows.Next() {
		var s Snippet
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	return snippets, nil
}
