package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "github.com/Saidurbu/go-lang-crud/internal/config"
	"github.com/Saidurbu/go-lang-crud/internal/handlers/student"
	"github.com/Saidurbu/go-lang-crud/internal/middleware"
	"github.com/Saidurbu/go-lang-crud/internal/storage/sqlite"
)

func main() {

	cfg := config.MustLoad()

	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Database connection established", "Environment", slog.String("env", cfg.Env))

	router := http.NewServeMux()

	router.HandleFunc("POST /api/registration", student.Registration(storage))
	router.HandleFunc("POST /api/login", student.Login(storage))

	router.HandleFunc("GET /api/profile", middleware.JWTAuth(student.GetProfile(storage)))

	router.HandleFunc("GET /api/students", middleware.JWTAuth(student.GetList(storage)))
	router.HandleFunc("POST /api/students", middleware.JWTAuth(student.New(storage)))
	router.HandleFunc("GET /api/students/{id}", middleware.JWTAuth(student.GetById(storage)))
	router.HandleFunc("PUT /api/students/{id}", middleware.JWTAuth(student.Update(storage)))
	router.HandleFunc("DELETE /api/students/{id}", middleware.JWTAuth(student.Delete(storage)))

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	fmt.Printf("Starting server on %s\n", cfg.Addr)

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}

	}()

	<-done

	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Error shutting down server:", slog.String("error", err.Error()))
	}

	slog.Info("Server gracefully stopped")
}
