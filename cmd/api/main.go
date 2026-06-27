package main

import (
	"context"
	"do_it_back/internal/auth"
	"do_it_back/internal/config"
	"do_it_back/internal/middleware"
	"do_it_back/internal/pkg"
	"do_it_back/internal/task"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		slog.Error("Error loading .env file", "error", err)
	}
	cfg := config.Load()
	if err := run(cfg); err != nil {
		slog.Error("Server error", "error", err)
	}
}

func run(cfg *config.Config) error {
	// ==== DATABASE =====
	config, err := pgxpool.ParseConfig(cfg.DatabaseRrl)
	if err != nil {
		return err
	}
	config.MaxConns = 25
	config.MinConns = 5

	ctx := context.Background()
	pool, err := pgxpool.NewWithConfig(ctx, config)

	if err != nil {
		return err
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		return err
	}

	api := http.NewServeMux()

	// ===== MIDDLEWARE =====
	server := middleware.LoggerMiddleware(api)
	server = middleware.AuthMiddleware(cfg, server)
	http.Handle("/api/", http.StripPrefix("/api", server))

	// ===== ROUTES =====
	// ROUTER: HEALTH
	api.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := pool.Ping(ctx); err != nil {
			slog.Error("Database error", "error", err)
			pkg.EncodeJSON(w, pkg.Response{Error: "internal server error"}, http.StatusInternalServerError)
			return
		}
		pkg.EncodeJSON(w, pkg.Response{Data: "ok"}, http.StatusOK)
	})

	// ROUTER: AUTH
	auth.NewAuthModule(api, cfg, pool)

	// ROUTER: TASK
	task.NewTaskModule(api, pool)

	addr := ":" + cfg.Port
	srv := &http.Server{
		Addr:              addr,
		Handler:           http.DefaultServeMux,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	// Channel to signal server shutdown
	shutdownCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	serverErr := make(chan error, 1)
	go func() {
		slog.Info("starting server", "addr", addr)
		serverErr <- srv.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		return err
	case <-shutdownCtx.Done():
		slog.Info("shutting down server")
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(timeoutCtx); err != nil {
			return err // o Shutdown estourou o prazo de 10s
		}
		return nil
	}
}
