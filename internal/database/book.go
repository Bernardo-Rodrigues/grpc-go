package database

import (
	"database/sql"

	"github.com/google/uuid"
)

type Book struct {
	db          *sql.DB
	ID          string
	Name        string
	Description string
	AuthorID    string
}

func NewBook(db *sql.DB) *Book {
	return &Book{db: db}
}

func (c *Book) Create(name, description, authorID string) (*Book, error) {
	id := uuid.New().String()
	_, err := c.db.Exec("INSERT INTO books (id, name, description, author_id) VALUES ($1, $2, $3, $4)",
		id, name, description, authorID)
	if err != nil {
		return nil, err
	}
	return &Book{
		ID:          id,
		Name:        name,
		Description: description,
		AuthorID:    authorID,
	}, nil
}

func (c *Book) FindAll() ([]Book, error) {
	rows, err := c.db.Query("SELECT id, name, description, author_id FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	books := []Book{}
	for rows.Next() {
		var id, name, description, authorID string
		if err := rows.Scan(&id, &name, &description, &authorID); err != nil {
			return nil, err
		}
		books = append(books, Book{ID: id, Name: name, Description: description, AuthorID: authorID})
	}
	return books, nil
}

func (c *Book) FindByAuthorID(authorID string) ([]Book, error) {
	rows, err := c.db.Query("SELECT id, name, description, author_id FROM books WHERE author_id = $1", authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	books := []Book{}
	for rows.Next() {
		var id, name, description, authorID string
		if err := rows.Scan(&id, &name, &description, &authorID); err != nil {
			return nil, err
		}
		books = append(books, Book{ID: id, Name: name, Description: description, AuthorID: authorID})
	}
	return books, nil
}

func (c *Book) Find(id string) (Book, error) {
	var name, description, authorID string
	err := c.db.QueryRow("SELECT name, description, author_id FROM books WHERE id = $1", id).
		Scan(&name, &description, &authorID)
	if err != nil {
		return Book{}, err
	}
	return Book{ID: id, Name: name, Description: description, AuthorID: authorID}, nil
}
