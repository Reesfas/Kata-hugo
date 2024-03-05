package main

import (
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	jsoniter "github.com/json-iterator/go"
	"github.com/ptflp/godecoder"
	httpSwagger "github.com/swaggo/http-swagger"
	"gitlab.com/ptflp/goboilerplate/config"
	"golang.org/x/net/context"
	_ "hugoproxy-main/task3.3.3/proxy/docs"
	"hugoproxy-main/task3.3.3/proxy/internal/controller"
	"hugoproxy-main/task3.3.3/proxy/internal/service"
	"log"
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
	decoder := godecoder.NewDecoder(jsoniter.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		DisallowUnknownFields:  true,
	})
	conf := config.NewAppConf()
	logger := NewLogger(conf, os.Stdout)
	resp := controller.NewResponder(decoder, logger)
	geo := service.NewGeoSerive()
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
		//r.Use(jwtauth.Verifier(tokenAuth))
		//r.Use(jwtauth.Authenticator)

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
