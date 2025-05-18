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

	router.HandleFunc("POST /api/students", student.New(storage))

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
