package service

import (
	"errors"
	"task354/internal/repository"
)

type User struct {
	ID          int
	Name        string
	RentedBooks []Book
}

type Book struct {
	ID     int
	Title  string
	Author Authors
	Rented bool
}

type Authors struct {
	Id    int
	Name  string
	Books []Book
}

type LibraryService struct {
	UserRepository   repository.UserRepository
	BookRepository   repository.BookRepository
	AuthorRepository repository.AuthorRepository
	RentalRepository repository.RentalRepository
}

func (s *LibraryService) RentBook(userID, bookID int) error {
	// Проверка существования пользователя
	_, err := s.UserRepository.GetUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Проверка существования книги
	_, err = s.BookRepository.GetBookByID(bookID)
	if err != nil {
		return errors.New("book not found")
	}

	// Если все проверки прошли успешно, арендуем книгу
	return s.RentalRepository.RentBook(bookID, userID)
}

func (s *LibraryService) ReturnBook(bookID, userID int) error {
	err := s.RentalRepository.ReturnBook(bookID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *LibraryService) AddBook(title string, authorId int) error {
	err := s.BookRepository.AddBook(title, authorId)
	if err != nil {
		return err
	}
	return nil
}

func (s *LibraryService) GetAllBooks() ([]Book, error) {
	books, err := s.BookRepository.GetAllBooks()
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (s *LibraryService) AddAuthor(name string) {
	s.AuthorRepository.AddAuthor(name)
}

func (s *LibraryService) GetAuthorsWithBooks() ([]Authors, error) {
	authors, err := s.AuthorRepository.GetAuthorsWithBooks()
	if err != nil {
		return nil, err
	}
	return authors, err
}

func (s *LibraryService) GetTopAuthors(limit int) ([]Authors, error) {
	return s.AuthorRepository.GetTopAuthors(limit)
}

func (s *LibraryService) CreateUser(name string) error {
	err := s.UserRepository.CreateUser(name)
	if err != nil {
		return err
	}
	return nil
}

func (s *LibraryService) GetUsers() ([]User, error) {
	books, err := s.UserRepository.GetUsersWithRentedBooks()
	if err != nil {
		return nil, err
	}
	return books, nil
}
