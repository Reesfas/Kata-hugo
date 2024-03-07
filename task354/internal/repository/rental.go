package repository

import (
	"database/sql"
	"errors"
	"time"
)

type RentalRepository struct {
	*sql.DB
}

func (r *RentalRepository) ReturnBook(bookID, userID int) error {
	// Проверяем, арендована ли книга пользователем
	isRented, err := r.IsBookRentedByUser(bookID, userID)
	if err != nil {
		return err
	}
	if !isRented {
		return errors.New("book is not rented by user")
	}

	// Обновляем запись об аренде, устанавливая дату возврата
	_, err = r.DB.Exec("UPDATE Rentals SET ReturnDate = NOW() WHERE BookID = $1 AND UserID = $2", bookID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *RentalRepository) IsBookRentedByUser(bookID, userID int) (bool, error) {
	var count int
	err := r.DB.QueryRow("SELECT COUNT(*) FROM Rentals WHERE BookID = $1 AND UserID = $2 AND ReturnDate IS NULL", bookID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *RentalRepository) RentBook(bookID, userID int) error {
	isRented, err := r.IsBookRented(bookID)
	if err != nil {
		return err
	}
	if isRented {
		return errors.New("book already rented")
	}

	// Создаем запись об аренде книги пользователем
	_, err = r.DB.Exec("INSERT INTO rentals (UserID, BookID, RentDate) VALUES ($1, $2, $3)", userID, bookID, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (r *RentalRepository) IsBookRented(bookID int) (bool, error) {
	var rented bool
	err := r.DB.QueryRow("SELECT IsRented FROM Books WHERE ID = $1 AND return_date IS NULL", bookID).Scan(&rented)
	if err != nil {
		return false, err
	}
	return rented, nil
}
