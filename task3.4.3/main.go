package main

import (
	"database/sql"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
	"task3.4.3/internal/controller"
	"task3.4.3/internal/repository"
	"task3.4.3/internal/service"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
}

func main() {
	db, err := sql.Open("postgres", "postgres://username:password@localhost/dbname?sslmode=disable")
	if err != nil {
		log.Fatal("Could not connect to the database:", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	petRepo := repository.NewPetRep(db)
	orderRepo := repository.NewOrderRep(db)

	userService := service.NewUserService(userRepo)
	petService := service.NewPetServ(petRepo)
	orderService := service.NewOrderServ(orderRepo)

	userController := controller.NewUserController(userService)
	petController := controller.NewPetRep(petService)
	orderController := controller.NewOrderRep(orderService)

	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Post("/users", userController.CreateUser)
		r.Get("/users/{username}", userController.GetUser)
		r.Put("/users/{username}", userController.UpdateUser)
		r.Delete("/users/{id}", userController.DeleteUser)
		r.Get("/users", userController.ListUsers)
		r.Post("/users/login", userController.Login)
		r.Post("/users/logout", userController.Logout)
	})

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/pets", petController.Create)
		r.Get("/pets/{id}", petController.GetByID)
		r.Get("/pets/status/{status}", petController.GetByStatus)
		r.Post("/pets/uploadImages", petController.UploadImages)
		r.Put("/pets/{id}", petController.FullUpdate)
		r.Patch("/pets/{id}", petController.PartialUpdate)
		r.Delete("/pets/{id}", petController.Delete)
	})

	r.Group(func(r chi.Router) {
		r.Post("/orders", orderController.Create)
		r.Get("/orders/{id}", orderController.GetByID)
		r.Delete("/orders/{id}", orderController.Delete)

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator)
			r.Get("/orders/inventory", orderController.GetInventory)
		})
	})

	r.HandleFunc("/swagger/*", httpSwagger.Handler())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
