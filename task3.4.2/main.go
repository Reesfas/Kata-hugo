package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	"hugoproxy-main/task3.4.2/internal/controller"
	"hugoproxy-main/task3.4.2/internal/repository"
	"hugoproxy-main/task3.4.2/internal/service"
	"log"
	"net/http"
	"os"
)

func main() {
	db, err := sql.Open("postgres", "user="+os.Getenv("DB_USER")+" password="+os.Getenv("DB_PASSWORD")+" dbname="+os.Getenv("DB_NAME")+" sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService)

	r := mux.NewRouter()

	r.HandleFunc("/api/users", userController.CreateUser).Methods("POST")
	r.HandleFunc("/api/users/{id}", userController.GetUser).Methods("GET")
	r.HandleFunc("/api/users/{id}", userController.UpdateUser).Methods("PUT")
	r.HandleFunc("/api/users/{id}", userController.DeleteUser).Methods("DELETE")
	r.HandleFunc("/api/users", userController.ListUsers).Methods("GET")

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		return
	}
}
