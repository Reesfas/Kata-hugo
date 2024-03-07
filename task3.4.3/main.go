package main

import (
	"database/sql"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"log"
	"net/http"
	"os"
	"strings"
	"task3.4.3/internal/controller"
	"task3.4.3/internal/repository"
	"task3.4.3/internal/service"
)

func main() {
	// Подключение к базе данных PostgreSQL
	db, err := sql.Open("postgres", "postgres://username:password@localhost/dbname?sslmode=disable")
	if err != nil {
		log.Fatal("Could not connect to the database:", err)
	}
	defer db.Close()

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(db)
	petRepo := repository.NewPetRep(db)
	orderRepo := repository.NewOrderRep(db)

	// Инициализация сервисов
	userService := service.NewUserService(userRepo)
	petService := service.NewPetServ(petRepo)
	orderService := service.NewOrderServ(orderRepo)

	// Инициализация контроллеров
	userController := controller.NewUserController(userService)
	petController := controller.NewPetRep(petService)
	orderController := controller.NewOrderRep(orderService)

	// Инициализация маршрутизатора Gorilla Mux
	r := mux.NewRouter()

	// Регистрация обработчиков маршрутов
	r.HandleFunc("/users", userController.CreateUser).Methods("POST")
	r.HandleFunc("/users/{username}", userController.GetUser).Methods("GET")
	r.HandleFunc("/users/{username}", userController.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", userController.DeleteUser).Methods("DELETE")
	r.HandleFunc("/users", userController.ListUsers).Methods("GET")
	r.HandleFunc("/users/login", userController.Login).Methods("POST")
	r.HandleFunc("/users/logout", userController.Logout).Methods("POST")

	r.HandleFunc("/pets", petController.Create).Methods("POST")
	r.HandleFunc("/pets/{id}", petController.GetByID).Methods("GET")
	r.HandleFunc("/pets/status/{status}", petController.GetByStatus).Methods("GET")
	r.HandleFunc("/pets/uploadImages", petController.UploadImages).Methods("POST")
	r.HandleFunc("/pets/{id}", petController.FullUpdate).Methods("PUT")
	r.HandleFunc("/pets/{id}", petController.PartialUpdate).Methods("PATCH")
	r.HandleFunc("/pets/{id}", petController.Delete).Methods("DELETE")

	r.HandleFunc("/orders", orderController.Create).Methods("POST")
	r.HandleFunc("/orders/{id}", orderController.GetByID).Methods("GET")
	r.HandleFunc("/orders/{id}", orderController.Delete).Methods("DELETE")
	r.HandleFunc("/orders/inventory", orderController.GetInventory).Methods("GET")

	n := negroni.New()

	n.Use(negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		authMiddleware(next).ServeHTTP(w, r)
	}))

	n.UseHandler(r)

	// Запуск HTTP сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, n))
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization token required", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			http.Error(w, "Invalid token format", http.StatusBadRequest)
			return
		}

		tokenString := bearerToken[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Здесь вы должны вернуть ключ проверки подписи (signing key)
			return []byte("your-secret-key"), nil
		})
		if err != nil {
			http.Error(w, "Failed to parse token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
