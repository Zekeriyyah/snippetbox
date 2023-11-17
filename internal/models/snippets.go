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

// Insert new snippet into the database
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	// Statement to insert new instance of snippet in the database
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?,?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// Executing the database statement
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	//Access the index of newly created snippets from its table
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Return a specific snippet from the database by ID
func (m *SnippetModel) Get(id int) (*Snippet, error) {

	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id=?`

	//Query row having ID=id
	row := m.DB.QueryRow(stmt, id)

	//Initialize snippet struct to hold the data for the return snippet
	s := &Snippet{}

	//Scanning the return row into the Snippet object instance s
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	// if everything is fine
	return s, nil
}

// Return the 10 most recent snippets
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	//Write SQL statement
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	//Query the database using Query()
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	//Create list of instances of Snippet
	snippets := []*Snippet{}

	//iterate through and scan each snippet into the instance
	for rows.Next() {
		s := &Snippet{}

		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	// If the loop is not successful, return the error
	if err = rows.Err(); err != nil {
		return nil, err
	}

	//return the list and log error
	return snippets, nil
}
