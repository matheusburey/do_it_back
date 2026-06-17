package auth

import (
	"do_it_back/internal/config"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func AuthNewModule(mux *http.ServeMux, cfg *config.Config, db *pgxpool.Pool) {
	r := NewRepository(db)
	s := NewService(r)
	h := NewHandler(cfg, s)

	mux.HandleFunc("POST /auth/login", h.Login)
	mux.HandleFunc("POST /auth/register", h.Register)
}
