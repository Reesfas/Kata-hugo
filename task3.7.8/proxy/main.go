package main

import (
	"database/sql"
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/joho/godotenv"
	jsoniter "github.com/json-iterator/go"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/ptflp/godecoder"
	httpSwagger "github.com/swaggo/http-swagger"
	"gitlab.com/ptflp/goboilerplate/config"
	"golang.org/x/net/context"
	_ "hugo/task3.7.8/proxy/docs"
	"hugo/task3.7.8/proxy/internal/controller"
	"hugo/task3.7.8/proxy/internal/repository"
	"hugo/task3.7.8/proxy/internal/service"
	"log"
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
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
		log.Fatalf("Something wrong with database %v", err)
	}
	dir := "./migrations"

	if err = goose.Up(db, dir); err != nil {
		log.Fatalf("failed to apply migrations: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	proxy := NewReverseProxy("hugo", "1313")
	err = os.Setenv("HOST", proxy.host)
	if err != nil {
		return
	}

	decoder := godecoder.NewDecoder(jsoniter.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		DisallowUnknownFields:  true,
	})

	conf := config.NewAppConf()
	logger := NewLogger(conf, os.Stdout)
	resp := controller.NewResponder(decoder, logger)
	repo := repository.NewGeoRep(db)
	geo := service.NewGeoSerive(repo)
	geocode := controller.NewGeo(geo, resp)
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

		r.Mount("/mycustompath/pprof", pprofHandler())

		r.Post("/api/address/search", geocode.Search)
		r.Post("/api/address/geocode", geocode.GeocodeAddress)
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
	if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server error: %v", err)
	}
	<-stop

}

func pprofHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch strings.TrimPrefix(r.URL.Path, "/mycustompath/pprof") {
		case "/allocs", "/block", "/cmdline", "/goroutine", "/heap", "/mutex", "/profile", "/threadcreate", "/trace":
			pprof.Index(w, r)
		case "/":
			pprof.Index(w, r)
		default:
			pprof.Handler(strings.TrimPrefix(r.URL.Path, "/mycustompath/pprof")).ServeHTTP(w, r)
		}
	})
}
