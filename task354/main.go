package main

import (
	"1.19/internal/service"
	"database/sql"
	"github.com/brianvoe/gofakeit/v7"
	"log"
)

func main() {
	// Подключение к базе данных
	db, err := sql.Open("postgres", "user=username dbname=dbname sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Инициализация таблицы авторов
	err = initializeAuthorsTable(db)
	if err != nil {
		log.Fatal(err)
	}
	err = initializeBooksTable(db)
	if err != nil {
		log.Fatal(err)
	}
	err = initializeUsersTable(db)
	if err != nil {
		log.Fatal(err)
	}
}

func initializeAuthorsTable(db *sql.DB) error {
	// Проверяем, пуста ли таблица авторов
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Authors").Scan(&count)
	if err != nil {
		return err
	}

	// Если таблица пуста, добавляем 10 авторов
	if count == 0 {
		// Генерируем и добавляем 10 случайных авторов
		for i := 0; i < 10; i++ {
			author := service.Authors{Name: gofakeit.Name()}
			_, err = db.Exec("INSERT INTO Authors (Name) VALUES ($1)", author.Name)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func initializeBooksTable(db *sql.DB) error {
	// Проверяем, пуста ли таблица книг
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Books").Scan(&count)
	if err != nil {
		return err
	}

	// Если таблица пуста, добавляем 100 книг
	if count == 0 {
		// Получаем существующих авторов
		authors, err := getExistingAuthors(db)
		if err != nil {
			return err
		}

		// Добавляем 10 книг для каждого автора
		for _, author := range authors {
			for i := 0; i < 10; i++ {
				book := service.Book{
					Title:  gofakeit.Book().Title,
					Author: author,
				}
				_, err = db.Exec("INSERT INTO Books (Title, AuthorID) VALUES ($1, $2)", book.Title, book.Author.Id)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func getExistingAuthors(db *sql.DB) ([]service.Authors, error) {
	rows, err := db.Query("SELECT ID FROM Authors")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []service.Authors
	for rows.Next() {
		var author service.Authors
		if err = rows.Scan(&author.Id); err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return authors, nil
}

func initializeUsersTable(db *sql.DB) error {
	// Проверяем, пуста ли таблица пользователей
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Users").Scan(&count)
	if err != nil {
		return err
	}

	// Если таблица пуста, добавляем больше 50 пользователей
	if count == 0 {
		// Генерируем и добавляем пользователей
		for i := 0; i < 50; i++ {
			user := service.User{Name: gofakeit.Name()}
			_, err := db.Exec("INSERT INTO Users (Name) VALUES ($1)", user.Name)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
