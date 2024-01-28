package main

import (
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"os"
)

func main() {
	r := chi.NewRouter()
	proxy := NewReverseProxy("hugo", "1313")
	err := os.Setenv("HOST", proxy.host)
	if err != nil {
		return
	}
	r.Use(proxy.ReverseProxy)
	r.Get("/api", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from API"))
	})
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Println(err)
		return
	}
}
