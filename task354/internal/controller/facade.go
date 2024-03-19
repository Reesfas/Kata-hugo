package controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"

	"task354/internal/service"
)

type FacadeController struct {
	facade *service.Facade
}

func NewFacadeController(facade *service.Facade) *FacadeController {
	return &FacadeController{
		facade: facade,
	}
}

// UsersList
// @Summary Get list of users
// @Description Returns a list of users.
// @Tags Users
// @Produce json
// @Success 200 {array} service.User
// @Failure 500 {string} string "Internal Server Error"
// @Router /users [get]
func (c *FacadeController) UsersList(w http.ResponseWriter, r *http.Request) {
	users, err := c.facade.UsersList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

// UserAdd
// @Summary Add a new user
// @Description Adds a new user to the system.
// @Tags Users
// @Accept json
// @Param user body service.User true "User object that needs to be added"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Router /users [post]
func (c *FacadeController) UserAdd(w http.ResponseWriter, r *http.Request) {
	var newUser service.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c.facade.UserAdd(newUser.Name)
	w.WriteHeader(http.StatusOK)
}

// AuthorsTop
// @Summary Get top authors
// @Description Returns the top authors based on the given limit.
// @Tags Authors
// @Param limit path int true "Limit of authors to be returned"
// @Success 200 {string} string "OK"
// @Router /authors/{limit} [get]
func (c *FacadeController) AuthorsTop(w http.ResponseWriter, r *http.Request) {
	limitStr := mux.Vars(r)["limit"]
	limit, _ := strconv.Atoi(limitStr)
	c.facade.AuthorsTop(limit)
	w.WriteHeader(http.StatusOK)
}

// AuthorsList
// @Summary Get list of authors
// @Description Returns a list of authors with their associated books.
// @Tags Authors
// @Produce json
// @Success 200 {array} service.Authors
// @Failure 500 {string} string "Internal Server Error"
// @Router /authors [get]
func (c *FacadeController) AuthorsList(w http.ResponseWriter, r *http.Request) {
	authors, err := c.facade.AuthorsList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(authors)
}

// AuthorAdd
// @Summary Add a new author
// @Description Adds a new author to the system.
// @Tags Authors
// @Accept json
// @Param author body service.Authors true "Author object that needs to be added"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Router /authors [post]
func (c *FacadeController) AuthorAdd(w http.ResponseWriter, r *http.Request) {
	var newAuthor service.Authors
	err := json.NewDecoder(r.Body).Decode(&newAuthor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c.facade.AuthorAdd(newAuthor.Name)
	w.WriteHeader(http.StatusOK)
}

// BookRent
// @Summary Rent a book
// @Description Allows a user to rent a book by providing the user ID and book ID.
// @Tags Books
// @Accept json
// @Param UserID body int true "User ID for renting a book"
// @Param BookID body int true "Book ID for renting a book"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Router /books/rent [post]
func (c *FacadeController) BookRent(w http.ResponseWriter, r *http.Request) {
	var requestData map[string]int
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, ok := requestData["UserID"]
	if !ok {
		http.Error(w, "UserID is missing", http.StatusBadRequest)
		return
	}

	bookID, ok := requestData["BookID"]
	if !ok {
		http.Error(w, "BookID is missing", http.StatusBadRequest)
		return
	}

	c.facade.BookRent(userID, bookID)
	w.WriteHeader(http.StatusOK)
}

// BookReturn
// @Summary Return a book
// @Description Allows a user to return a book by providing the user ID and book ID.
// @Tags Books
// @Accept json
// @Param UserID body int true "User ID for returning a book"
// @Param BookID body int true "Book ID for returning a book"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /books/return [post]
func (c *FacadeController) BookReturn(w http.ResponseWriter, r *http.Request) {
	var requestData map[string]int
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, ok := requestData["UserID"]
	if !ok {
		http.Error(w, "UserID is missing", http.StatusBadRequest)
		return
	}

	bookID, ok := requestData["BookID"]
	if !ok {
		http.Error(w, "BookID is missing", http.StatusBadRequest)
		return
	}

	err = c.facade.BookReturn(userID, bookID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// BookList
// @Summary Get list of books
// @Description Returns a list of books.
// @Tags Books
// @Produce json
// @Success 200 {array} service.Book
// @Failure 500 {string} string "Internal Server Error"
// @Router /books [get]
func (c *FacadeController) BookList(w http.ResponseWriter, r *http.Request) {
	books, err := c.facade.BookList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(books)
}

// BookAdd
// @Summary Add a new book
// @Description Adds a new book to the system.
// @Tags Books
// @Accept json
// @Param book body service.Book true "Book object that needs to be added"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Router /books [post]
func (c *FacadeController) BookAdd(w http.ResponseWriter, r *http.Request) {
	var newBook service.Book
	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c.facade.BookAdd(newBook.Title, newBook.Author.Id)
	w.WriteHeader(http.StatusOK)
}
