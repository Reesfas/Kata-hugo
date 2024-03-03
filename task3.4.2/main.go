package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"hugoproxy-main/task3.4.2/internal/controller"
	"hugoproxy-main/task3.4.2/internal/repository"
	"hugoproxy-main/task3.4.2/internal/service"
	"log"
	"net/http"
	"os"
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
		log.Fatalf("Some shit with database %v", err)
	}
	dir := "./migrations"

	if err = goose.Up(db, dir); err != nil {
		log.Fatalf("failed to apply migrations: %v", err)
	}
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService)
	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	r.HandleFunc("/api/users", userController.CreateUser).Methods("POST")
	r.HandleFunc("/api/users/{id}", userController.GetUser).Methods("GET")
	r.HandleFunc("/api/users/{id}", userController.UpdateUser).Methods("PUT")
	r.HandleFunc("/api/users/{id}", userController.DeleteUser).Methods("DELETE")
	r.HandleFunc("/api/users", userController.ListUsers).Methods("GET")

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server started")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
