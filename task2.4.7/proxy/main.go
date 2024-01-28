package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
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
	go WorkerTest()
}

func WorkerTest() {
	t := time.NewTicker(5 * time.Second)
	var b int = 0
	tree := GenerateTree(5)
	for {
		select {
		case tt := <-t.C:
			err := updateFile("/app/static/tasks/_index.md", tt, b)
			if err != nil {
				log.Println(err)
			}
			graph := generateRandomGraph(5, 30)
			err = updateGraphFile("/app/static/tasks/graph.md", graph)
			if err != nil {
				return
			}
			if treeSize(tree.Root) >= 100 {
				tree = GenerateTree(5)
			} else {
				key := rand.Intn(100)
				tree.Insert(key)
			}
			err = updateTree("/app/static/tasks/binary.md", tree)
			b++
		}
	}
}
