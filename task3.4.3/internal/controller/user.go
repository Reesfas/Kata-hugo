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

func (u *UserContr) Logout(w http.ResponseWriter, r *http.Request) {
	c := http.Cookie{
		Name:   "token",
		MaxAge: -1}
	http.SetCookie(w, &c)

	_, err := w.Write([]byte("Old cookie deleted. Logged out!\n"))
	if err != nil {
		return
	}
}

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
