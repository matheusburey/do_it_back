package main

import (
	"log/slog"
	"net/http"
)

func main() {
	const PORT = ":8080"
	if err := run(PORT); err != nil {
		slog.Error("Failed to marshal JSON", "error", err)
	}
}

func run(addr string) error {
	// ==== DATABASE =====

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
