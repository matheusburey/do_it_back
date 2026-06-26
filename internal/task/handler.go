package task

import (
	"context"
	"do_it_back/internal/pkg"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type CreateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	IsCompleted bool   `json:"is_completed"`
}

type UpdateRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	IsCompleted *bool   `json:"is_completed"`
}

func (cr CreateRequest) Valid(ctx context.Context) pkg.Evaluator {
	var eval pkg.Evaluator

	eval.CheckField(pkg.NotBlank(cr.Title), "title", "title is required")
	eval.CheckField(pkg.MinLength(cr.Title, 5) && pkg.MaxLength(cr.Title, 100), "title", "min length is 5 and max length is 100")
	eval.CheckField(pkg.NotBlank(cr.Description), "description", "description is required")
	eval.CheckField(pkg.MinLength(cr.Description, 5) && pkg.MaxLength(cr.Description, 100), "description", "min length is 5 and max length is 100")

	return eval
}

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) ListTasks(
	w http.ResponseWriter,
	r *http.Request,
) {
	user_id, ok := pkg.UserIDFromContext(r.Context())

	if !ok {
		pkg.EncodeJSON(w, pkg.Response{Error: "invalid user"}, http.StatusNotFound)
		return
	}
	tasks, err := h.service.GetTasks(r.Context(), user_id)

	if err != nil {
		pkg.EncodeJSON(w, pkg.Response{Error: "invalid user"}, http.StatusNotFound)
		return
	}
	pkg.EncodeJSON(w, pkg.Response{Data: tasks}, http.StatusOK)
}

func (h *Handler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	data, problems, err := pkg.DecodeValidJSON[CreateRequest](r)
	if err != nil {
		pkg.EncodeJSON(
			w,
			pkg.Response{Error: "One or more fields are invalid.", Fields: problems},
			http.StatusBadRequest,
		)
		return
	}

	user_id, ok := pkg.UserIDFromContext(r.Context())

	if !ok {
		pkg.EncodeJSON(w, pkg.Response{Error: "invalid user"}, http.StatusNotFound)
		return
	}

	task, err := h.service.Create(
		r.Context(),
		user_id,
		data.Title,
		data.Description,
		data.IsCompleted,
	)

	if err != nil {
		pkg.EncodeJSON(
			w,
			pkg.Response{Error: err.Error()},
			http.StatusBadRequest,
		)
		return
	}

	pkg.EncodeJSON(w, pkg.Response{Data: task}, http.StatusCreated)
}

func (h *Handler) GetTaskById(
	w http.ResponseWriter,
	r *http.Request,
) {
	user_id, ok := pkg.UserIDFromContext(r.Context())

	if !ok {
		pkg.EncodeJSON(w, pkg.Response{Error: "invalid user"}, http.StatusNotFound)
		return
	}
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		pkg.EncodeJSON(w, pkg.Response{Error: "invalid id"}, http.StatusBadRequest)
		return
	}

	task, err := h.service.GetTaskById(r.Context(), id, user_id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			pkg.EncodeJSON(w, pkg.Response{Error: "task not found"}, http.StatusNotFound)
			return

		}
		pkg.EncodeJSON(w, pkg.Response{Error: "internal server error"}, http.StatusInternalServerError)
		return
	}
	pkg.EncodeJSON(w, pkg.Response{Data: task}, http.StatusOK)
}

func (h *Handler) Update(
	w http.ResponseWriter,
	r *http.Request,
) {
	user_id, ok := pkg.UserIDFromContext(r.Context())

	if !ok {
		pkg.EncodeJSON(w, pkg.Response{Error: "invalid user"}, http.StatusNotFound)
		return
	}
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		pkg.EncodeJSON(w, pkg.Response{Error: "invalid id"}, http.StatusBadRequest)
		return
	}

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.EncodeJSON(
			w,
			pkg.Response{Error: "invalid body"},
			http.StatusBadRequest,
		)
		return
	}

	task, err := h.service.Update(r.Context(), id, user_id, req.IsCompleted, req.Title, req.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			pkg.EncodeJSON(w, pkg.Response{Error: "task not found"}, http.StatusNotFound)
			return

		}
		pkg.EncodeJSON(w, pkg.Response{Error: "internal server error"}, http.StatusInternalServerError)
		return
	}
	pkg.EncodeJSON(w, pkg.Response{Data: task}, http.StatusOK)
}

func (h *Handler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {
	user_id, ok := pkg.UserIDFromContext(r.Context())

	if !ok {
		pkg.EncodeJSON(w, pkg.Response{Error: "invalid user"}, http.StatusNotFound)
		return
	}
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		pkg.EncodeJSON(w, pkg.Response{Error: "invalid id"}, http.StatusBadRequest)
		return
	}

	err = h.service.Delete(r.Context(), id, user_id)
	if err != nil {
		if err.Error() == "task not found" {
			pkg.EncodeJSON(w, pkg.Response{Error: "task not found"}, http.StatusNotFound)
			return

		}
		pkg.EncodeJSON(w, pkg.Response{Error: "internal server error"}, http.StatusInternalServerError)
		return
	}
	pkg.EncodeJSON(w, pkg.Response{}, http.StatusNoContent)
}
