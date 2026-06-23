package task

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewTaskModule(mux *http.ServeMux, db *pgxpool.Pool) {
	r := NewRepository(db)
	s := NewService(r)
	h := NewHandler(s)

	mux.HandleFunc("GET /tasks", h.ListTasks)
	mux.HandleFunc("POST /task", h.Create)
	mux.HandleFunc("GET /task/{id}", h.GetTaskById)
	mux.HandleFunc("PUT /task/{id}", h.Update)
	mux.HandleFunc("DELETE /task/{id}", h.Delete)
}
