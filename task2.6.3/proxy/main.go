package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "hugoproxy-main/proxy/docs"
	"log"
	"net/http"
	"os"
	"time"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
}

// @title My Title
// @version 1.0
// @description Some useful description
// @host localhost:8000
// @BasePath /
func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	proxy := NewReverseProxy("hugo", "1313")
	err := os.Setenv("HOST", proxy.host)
	if err != nil {
		return
	}
	r.Use(proxy.ReverseProxy)
	r.Get("/swagger/*", httpSwagger.Handler())
	r.Get("/api", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from API"))
	})
	r.Post("/api/address/search", Search)
	r.Post("/api/address/geocode", geocodeAddress)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Println(err)
		return
	}
}
