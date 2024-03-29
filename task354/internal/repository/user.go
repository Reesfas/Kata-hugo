package repository

import (
	"database/sql"
	"errors"
	"task354/internal/service"
)

type UserRepositoryInterface interface {
	CreateUser(name string) error
	GetUserByID(userID int) (*service.User, error)
	GetUsersWithRentedBooks() ([]service.User, error)
}

type UserRepository struct {
	*sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(name string) error {
	query := "INSERT INTO Users (Name) VALUES ($1);"
	_, err := r.DB.Exec(query, name)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetUserByID(userID int) (*service.User, error) {
	var user service.User
	err := r.DB.QueryRow("SELECT id, name FROM users WHERE id = $1", userID).Scan(&user.ID, &user.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUsersWithRentedBooks() ([]service.User, error) {
	query := `
		SELECT u.ID, u.Name, b.ID AS BookID, b.Title, a.ID AS AuthorID, a.Name AS AuthorName
		FROM Users u
		LEFT JOIN Rentals r ON u.ID = r.UserID
		LEFT JOIN Books b ON r.BookID = b.ID
		LEFT JOIN Authors a ON b.AuthorID = a.ID;
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Мапим результаты запроса в список пользователей с их арендованными книгами
	usersMap := make(map[int]*service.User)
	for rows.Next() {
		var user service.User
		var book service.Book
		var author service.Authors
		err = rows.Scan(&user.ID, &user.Name, &book.ID, &book.Title, &author.Id, &author.Name)
		if err != nil {
			return nil, err
		}
		// Если пользователь встречается впервые, добавляем его в map
		if _, ok := usersMap[user.ID]; !ok {
			usersMap[user.ID] = &user
		}
		// Добавляем книгу в список арендованных книг пользователя
		usersMap[user.ID].RentedBooks = append(usersMap[user.ID].RentedBooks, service.Book{
			ID:     book.ID,
			Title:  book.Title,
			Author: service.Authors{Id: author.Id, Name: author.Name},
		})
	}

	// Преобразуем map в список пользователей
	var users []service.User
	for _, user := range usersMap {
		users = append(users, *user)
	}

	return users, nil
}
