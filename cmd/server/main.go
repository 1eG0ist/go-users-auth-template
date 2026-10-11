package main

import (
	"context"
	"errors"
	"go-users-auth-template/internal/users"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Api working correctly"))
	})

	userHandlers := users.NewUserHandlers(users.NewUserService(users.NewUserInMemoryRepository()))

	r.Mount("/users", userHandlers.Routes())

	server := &http.Server{
		Addr:    ":8021",
		Handler: r,
	}

	// channel for waiting exit
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("server starting on port 8021 [http://127.0.0.1:8021/health]", "address", server.Addr)

		if err := http.ListenAndServe(":8021", r); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error starting server: %s", err)
			os.Exit(1)
		}
	}()

	/*
		shutdown section
	*/
	<-signalChan

	slog.Info("server shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown failer", "error", err)
		os.Exit(1)
	}
	slog.Info("server stopped")
}
