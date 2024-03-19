package controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"task3.4.3/internal/repository"
	"task3.4.3/internal/service"
)

type UserController interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
	ListUsers(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

type UserContr struct {
	serv service.UserService
}

func NewUserController(userService service.UserService) *UserContr {
	return &UserContr{userService}
}

// CreateUser
// @Summary Create user
// @Description Create a new user
// @Accept json
// @Produce json
// @Param user body repository.User true "User object that needs to be added"
// @Success 200 {object} string "Successful operation"
// @Failure 400 {string} string "Invalid user data"
// @Failure 500 {string} string "Internal server error"
// @Router /user [post]
func (u *UserContr) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user repository.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := u.serv.Create(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// GetUser
// @Summary Get user by username
// @Description Get user information by username
// @Produce json
// @Param username path string true "Username of the user to get"
// @Success 200 {object} repository.User "Successful operation"
// @Failure 404 {string} string "User not found"
// @Router /user/{username} [get]
func (u *UserContr) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	user, err := u.serv.GetByUsername(r.Context(), username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		return
	}
}

// UpdateUser
// @Summary Update user
// @Description Update an existing user
// @Accept json
// @Produce json
// @Param username path string true "Username of the user to update"
// @Param user body repository.User true "User object with updated information"
// @Success 200 {object} repository.User "Successful operation"
// @Failure 400 {string} string "Invalid user data"
// @Failure 500 {string} string "Internal server error"
// @Router /user/{username} [put]
func (u *UserContr) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	var user repository.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.Username = username

	err := u.serv.Update(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		return
	}
}

// DeleteUser
// @Summary Delete user
// @Description Delete user by username
// @Param username path string true "Username of the user to delete"
// @Success 204 "No Content"
// @Failure 500 {string} string "Internal server error"
// @Router /user/{username} [delete]
func (u *UserContr) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := u.serv.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListUsers
// @Summary List users with given conditions
// @Description List users based on provided conditions
// @Accept json
// @Produce json
// @Param conditions body repository.Conditions true "Conditions for filtering users"
// @Success 200 {array} repository.User "Successful operation"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal server error"
// @Router /user/createWithList [post]
func (u *UserContr) ListUsers(w http.ResponseWriter, r *http.Request) {
	var c repository.Conditions

	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		return
	}

	users, err1 := u.serv.List(r.Context(), c)
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		return
	}
}

// Login
// @Summary User login
// @Description Logs user into the system
// @Accept json
// @Produce json
// @Param username body string true "Username for login"
// @Param password body string true "Password for login"
// @Success 200 {object} string "Successful operation"
// @Failure 400 {string} string "Invalid request format"
// @Failure 500 {string} string "Internal server error"
// @Router /user/login [get]
func (u *UserContr) Login(w http.ResponseWriter, r *http.Request) {
	var user repository.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	token, err := u.serv.Login(r.Context(), user.Username, user.Password)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{"token": token})
	if err != nil {
		return
	}
}

// Logout
// @Summary User logout
// @Description Logs out current logged in user session
// @Success 200 {string} string "Old cookie deleted. Logged out!"
// @Router /user/logout [get]
func (u *UserContr) Logout(w http.ResponseWriter, r *http.Request) {
	c := http.Cookie{
		Name:   "token",
		MaxAge: -1}
	http.SetCookie(w, &c)

	_, err := w.Write([]byte("Logged out!\n"))
	if err != nil {
		return
	}
}

// CreateUserWithArray
// @Summary Create users with array input
// @Description Creates new users with an array of user objects
// @Accept json
// @Produce json
// @Param users body []repository.User true "Array of user objects to create"
// @Success 201 "Created"
// @Failure 400 {string} string "Failed to decode request body"
// @Failure 500 {string} string "Failed to create users"
// @Router /user/createWithArray [post]
func (u *UserContr) CreateUserWithArray(w http.ResponseWriter, r *http.Request) {
	var users []repository.User
	if err := json.NewDecoder(r.Body).Decode(&users); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	if err := u.serv.CreateUserWithArray(users); err != nil {
		http.Error(w, "Failed to create users", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// CreateUserWithList
// @Summary Create users with list input
// @Description Creates new users with a list of user objects
// @Accept json
// @Produce json
// @Param users body []repository.User true "List of user objects to create"
// @Success 201 "Created"
// @Failure 400 {string} string "Failed to decode request body"
// @Failure 500 {string} string "Failed to create users"
// @Router /user/createWithList [post]
func (u *UserContr) CreateUserWithList(w http.ResponseWriter, r *http.Request) {
	var users []repository.User
	if err := json.NewDecoder(r.Body).Decode(&users); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	if err := u.serv.CreateUserWithList(users); err != nil {
		http.Error(w, "Failed to create users", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
