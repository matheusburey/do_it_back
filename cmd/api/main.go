package main

import (
	"context"
	"do_it_back/internal/auth"
	"do_it_back/internal/config"
	"do_it_back/internal/middleware"
	"do_it_back/internal/pkg"
	"log/slog"
	"net/http"

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
		pkg.EncodeJSON(w, pkg.Response{Data: "ok"}, http.StatusOK)
	})

	// ROUTER: AUTH
	auth.AuthNewModule(api, cfg, pool)

	addr := ":" + cfg.Port
	slog.Info("Starting server on port" + addr)
	return http.ListenAndServe(addr, nil)
}
