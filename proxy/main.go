package main

import (
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	httpSwagger "github.com/swaggo/http-swagger"
	"golang.org/x/net/context"
	_ "hugoproxy-main/proxy/docs"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
}

// @title Your API Title
// @version 1.0
// @description Your API description. You can use Markdown here.
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
	r.Post("/api/register", register)
	r.Post("/api/login", login)

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Post("/api/address/search", Search)
		r.Post("/api/address/geocode", geocodeAddress)
	})

	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-stop
		log.Println("Получен стоп сигнал. Закрываем сервер...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Server shutdown error: %v", err)
		}
		log.Println("Shutting down server...")
		close(stop)
	}()
	log.Println("Сервер на порту :8080 открыт")
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server error: %v", err)
	}
	<-stop
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

/*gin.SetMode(gin.ReleaseMode)
	r := chi.NewRouter()
	s := gin.New()
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
	r.Post("/api/register", register)
	r.Post("/api/login", login)

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Post("/api/address/search", Search)
		r.Post("/api/address/geocode", geocodeAddress)
	})
	s.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	s.Run(":8080")
	corsHandler := cors.Default().Handler(r)
	go WorkerTest()
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Starting server...")
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()
		go func() {
		err = http.ListenAndServe(":8080	", corsHandler)
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("Server started")

	<-stopChan

	log.Println("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		log.Fatalf("Error during server shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}*/
