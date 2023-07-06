package database

import (
	"database/sql"

	"github.com/google/uuid"
)

type Author struct {
	db   *sql.DB
	ID   string
	Name string
}

func NewAuthor(db *sql.DB) *Author {
	return &Author{db: db}
}

func (c *Author) Create(name string) (Author, error) {
	id := uuid.New().String()
	_, err := c.db.Exec("INSERT INTO authors (id, name) VALUES ($1, $2)",
		id, name)
	if err != nil {
		return Author{}, err
	}
	return Author{ID: id, Name: name}, nil
}

func (c *Author) FindAll() ([]Author, error) {
	rows, err := c.db.Query("SELECT id, name FROM authors")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	authors := []Author{}
	for rows.Next() {
		var id, name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}
		authors = append(authors, Author{ID: id, Name: name})
	}
	return authors, nil
}

func (c *Author) FindByBookID(bookID string) (Author, error) {
	var id, name string
	err := c.db.QueryRow("SELECT a.id, a.name FROM authors a JOIN books bo ON a.id = bo.author_id WHERE bo.id = $1", bookID).
		Scan(&id, &name)
	if err != nil {
		return Author{}, err
	}
	return Author{ID: id, Name: name}, nil
}

func (c *Author) Find(id string) (Author, error) {
	var name string
	err := c.db.QueryRow("SELECT name FROM authors WHERE id = $1", id).
		Scan(&name)
	if err != nil {
		return Author{}, err
	}
	return Author{ID: id, Name: name}, nil
}
