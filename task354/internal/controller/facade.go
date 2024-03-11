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

func (c *FacadeController) UsersList(w http.ResponseWriter, r *http.Request) {
	users, err := c.facade.UsersList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

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

func (c *FacadeController) AuthorsTop(w http.ResponseWriter, r *http.Request) {
	limitStr := mux.Vars(r)["limit"]
	limit, _ := strconv.Atoi(limitStr)
	c.facade.AuthorsTop(limit)
	w.WriteHeader(http.StatusOK)
}

func (c *FacadeController) AuthorsList(w http.ResponseWriter, r *http.Request) {
	authors, err := c.facade.AuthorsList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(authors)
}

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

func (c *FacadeController) BookList(w http.ResponseWriter, r *http.Request) {
	books, err := c.facade.BookList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(books)
}

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
