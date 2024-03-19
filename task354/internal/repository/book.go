package repository

import (
	"database/sql"
	"errors"
	"task354/internal/service"
)

type BookRentalCount struct {
	BookID      int
	RentalCount int
}

type BookRepositoryInterface interface {
	GetBookByID(bookID int) (*service.Book, error)
	AddBook(title string, authorID int) error
	GetAllBooks() ([]service.Book, error)
}

type BookRepository struct {
	DB *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{DB: db}
}

func (r *BookRepository) GetBookByID(bookID int) (*service.Book, error) {
	var book service.Book
	err := r.DB.QueryRow("SELECT id, title, author_id FROM books WHERE id = $1", bookID).Scan(&book.ID, &book.Title, &book.Author.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("book not found")
		}
		return nil, err
	}
	return &book, nil
}

func (r *BookRepository) AddBook(title string, authorID int) error {
	var authorExists bool
	err := r.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM Authors WHERE ID = $1)", authorID).Scan(&authorExists)
	if err != nil {
		return err
	}
	if !authorExists {
		return errors.New("указанного автора нет в списке")
	}

	// Если автор существует, вставляем новую книгу
	_, err = r.DB.Exec("INSERT INTO Books (Title, AuthorID) VALUES ($1, $2)", title, authorID)
	if err != nil {
		return err
	}

	return nil
}

func (r *BookRepository) GetAllBooks() ([]service.Book, error) {
	query := `
		SELECT b.ID AS BookID, b.Title AS BookTitle, a.ID AS AuthorID, a.Name AS AuthorName
		FROM Books b
		JOIN Authors a ON b.AuthorID = a.ID;
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []service.Book
	for rows.Next() {
		var book service.Book
		err = rows.Scan(&book.ID, &book.Title, &book.Author.Id, &book.Author.Name)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}
