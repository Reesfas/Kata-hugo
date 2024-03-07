package repository

import (
	"1.19/internal/service"
	"database/sql"
)

type AuthorRepository struct {
	*sql.DB
}

func (r *AuthorRepository) AddAuthor(name string) (int, error) {
	stmt, err := r.DB.Prepare("INSERT INTO Authors (Name) VALUES ($1) RETURNING ID")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var authorID int
	err = stmt.QueryRow(name).Scan(&authorID)
	if err != nil {
		return 0, err
	}

	return authorID, nil
}

func (r *AuthorRepository) GetAuthorsWithBooks() ([]service.Authors, error) {
	query := `
		SELECT a.ID AS AuthorID, a.Name AS AuthorName, b.ID AS BookID, b.Title AS BookTitle
		FROM Authors a
		LEFT JOIN Books b ON a.ID = b.AuthorID
		ORDER BY a.ID, b.ID;
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []service.Authors
	var currentAuthorID int
	var currentAuthor service.Authors

	for rows.Next() {
		var authorID int
		var authorName string
		var bookID sql.NullInt64
		var bookTitle sql.NullString

		if err = rows.Scan(&authorID, &authorName, &bookID, &bookTitle); err != nil {
			return nil, err
		}

		if authorID != currentAuthorID {
			if currentAuthorID != 0 {
				authors = append(authors, currentAuthor)
			}
			currentAuthorID = authorID
			currentAuthor = service.Authors{
				Id:   authorID,
				Name: authorName,
			}
		}

		if bookID.Valid && bookTitle.Valid {
			currentAuthor.Books = append(currentAuthor.Books, service.Book{
				ID:    int(bookID.Int64),
				Title: bookTitle.String,
			})
		}
	}
	if currentAuthorID != 0 {
		authors = append(authors, currentAuthor)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return authors, nil
}

func (r *AuthorRepository) GetTopAuthors(limit int) ([]service.Authors, error) {
	// Запрос для получения топа читаемых авторов
	query := `
		SELECT a.ID, a.Name, COUNT(b.ID) AS BookCount
		FROM Authors a
		LEFT JOIN Books b ON a.ID = b.AuthorID
		LEFT JOIN Rentals r ON b.ID = r.BookID
		GROUP BY a.ID, a.Name
		ORDER BY BookCount DESC
		LIMIT $1;
	`

	rows, err := r.DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topAuthors []service.Authors
	for rows.Next() {
		var author service.Authors
		err = rows.Scan(&author.Id, &author.Name)
		if err != nil {
			return nil, err
		}
		topAuthors = append(topAuthors, author)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return topAuthors, nil
}
