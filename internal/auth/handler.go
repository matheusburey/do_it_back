package auth

import (
	"do_it_back/internal/config"
	"do_it_back/internal/pkg"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type HandlerResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Token string    `json:"token"`
}

type Handler struct {
	service *Service
	cfg     *config.Config
}

func NewHandler(cfg *config.Config, service *Service) *Handler {
	return &Handler{
		service: service,
		cfg:     cfg,
	}
}

func (h *Handler) Register(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.EncodeJSON(
			w,
			pkg.Response{Error: "invalid body"},
			http.StatusBadRequest,
		)
		return
	}

	user, err := h.service.Register(
		r.Context(),
		req.Name,
		req.Email,
		req.Password,
	)

	if err != nil {
		pkg.EncodeJSON(
			w,
			pkg.Response{Error: err.Error()},
			http.StatusBadRequest,
		)
		return
	}

	token, err := pkg.GenerateAccessToken(user.ID, h.cfg.JWTSecret)
	if err != nil {
		pkg.EncodeJSON(w, pkg.Response{Error: err.Error()}, http.StatusInternalServerError)
		return
	}

	resp := HandlerResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Token: token,
	}

	pkg.EncodeJSON(w, pkg.Response{Data: resp}, http.StatusCreated)
}

func (h *Handler) Login(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.EncodeJSON(w, pkg.Response{Error: "invalid body"}, http.StatusBadRequest)
		return
	}

	user, err := h.service.Login(
		r.Context(),
		req.Email,
		req.Password,
	)

	if err != nil {
		pkg.EncodeJSON(w, pkg.Response{Error: "password or email is incorrect"}, http.StatusNotFound)
		return
	}

	token, err := pkg.GenerateAccessToken(user.ID, h.cfg.JWTSecret)
	if err != nil {
		pkg.EncodeJSON(w, pkg.Response{Error: err.Error()}, http.StatusInternalServerError)
		return
	}

	resp := HandlerResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Token: token,
	}
	pkg.EncodeJSON(w, pkg.Response{Data: resp}, http.StatusCreated)
}
