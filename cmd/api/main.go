package main

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	const PORT = ":8080"
	if err := run(PORT); err != nil {
		slog.Error("Failed to marshal JSON", "error", err)
	}
}

func run(addr string) error {
	// ==== DATABASE =====
	database_url := "postgres://postgres:1234@localhost:5432/todo"
	config, err := pgxpool.ParseConfig(database_url)
	if err != nil {
		panic(err)
	}
	config.MaxConns = 25
	config.MinConns = 5

	pool, err := pgxpool.NewWithConfig(context.Background(), config)

	if err != nil {
		panic(err)
	}
	defer pool.Close()
	if err := pool.Ping(context.Background()); err != nil {
		panic(err)
	}

	api := http.NewServeMux()

	// ===== MIDDLEWARE =====

	// ===== ROUTES =====
	// ROUTER: HEALTH
	api.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"Hello, World!"}`))
	})
	slog.Info("Starting server on port" + addr)
	return http.ListenAndServe(addr, nil)
}
