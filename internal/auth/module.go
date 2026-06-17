package auth

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func AuthNewModule(mux *http.ServeMux, db *pgxpool.Pool) {
	r := NewRepository(db)
	s := NewService(r)
	h := NewHandler(s)

	mux.HandleFunc("POST /auth/login", h.Login)
	mux.HandleFunc("POST /auth/register", h.Register)
}
