package main

import (
	"database/sql"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	_ "github.com/swaggo/http-swagger"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
	_ "task354/docs"
	"task354/internal/controller"
	"task354/internal/repository"
	"task354/internal/service"
)

func main() {
	err1 := godotenv.Load("db.env")
	if err1 != nil {
		log.Fatal("Ошибка загрузки файла .env")
	}

	db, err := sql.Open("postgres", "user="+os.Getenv("DB_USER")+" password="+os.Getenv("DB_PASSWORD")+" dbname="+os.Getenv("DB_NAME")+" host="+os.Getenv("DB_HOST")+" port="+os.Getenv("DB_PORT")+" sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Something wrong with database %v", err)
	}
	dir := "./migrations"

	if err = goose.Up(db, dir); err != nil {
		log.Fatalf("failed to apply migrations: %v", err)
	}
	author := repository.NewAuthorRepository(db)
	book := repository.NewBookRepository(db)
	rent := repository.NewRentalRepository(db)
	user := repository.NewUserRepository(db)
	serv := service.NewLibraryService(user, book, author, rent)
	fac := service.NewFacade(serv)
	facadeController := controller.NewFacadeController(fac)
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

	r := chi.NewRouter()
	r.Get("/swagger/*", httpSwagger.Handler())

	r.Get("/users", facadeController.UsersList)
	r.Post("/users", facadeController.UserAdd)
	r.Get("/authors/{limit}", facadeController.AuthorsTop)
	r.Get("/authors", facadeController.AuthorsList)
	r.Post("/authors", facadeController.AuthorAdd)
	r.Post("/books/rent", facadeController.BookRent)
	r.Post("/books/return", facadeController.BookReturn)
	r.Get("/books", facadeController.BookList)
	r.Post("/books", facadeController.BookAdd)

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}

func initializeAuthorsTable(db *sql.DB) error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Authors").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
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
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Books").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		authors, err := getExistingAuthors(db)
		if err != nil {
			return err
		}

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
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Users").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
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
