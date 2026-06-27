package auth

import (
	"context"
	"do_it_back/internal/config"
	"do_it_back/internal/pkg"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (rr RegisterRequest) Valid(ctx context.Context) pkg.Evaluator {
	var eval pkg.Evaluator

	eval.CheckField(pkg.NotBlank(rr.Name), "name", "name is required")
	eval.CheckField(pkg.MinLength(rr.Name, 5) && pkg.MaxLength(rr.Name, 100), "name", "min length is 5 and max length is 100")
	eval.CheckField(pkg.NotBlank(rr.Email), "email", "email is required")
	eval.CheckField(pkg.IsEmail(rr.Email), "email", "email is invalid")
	eval.CheckField(pkg.MinLength(rr.Email, 10) && pkg.MaxLength(rr.Email, 255), "email", "min length is 10 and max length is 255")
	eval.CheckField(pkg.NotBlank(rr.Password), "password", "password is required")
	eval.CheckField(pkg.MinLength(rr.Password, 8) && pkg.MaxLength(rr.Password, 255), "password", "min length is 8 and max length is 255")
	eval.CheckField(pkg.IsPassword(rr.Password), "password", "Password must contain at least 8 characters, one uppercase letter, one lowercase letter, one special character and one number.")

	return eval
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (lr LoginRequest) Valid(ctx context.Context) pkg.Evaluator {
	var eval pkg.Evaluator

	eval.CheckField(pkg.NotBlank(lr.Email), "email", "email is required")
	eval.CheckField(pkg.IsEmail(lr.Email), "email", "email is invalid")
	eval.CheckField(pkg.NotBlank(lr.Password), "password", "password is required")

	return eval
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

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	data, problems, err := pkg.DecodeValidJSON[RegisterRequest](w, r)
	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			pkg.EncodeJSON(w, pkg.Response{Error: "request body too large"}, http.StatusRequestEntityTooLarge)
			return
		}
		pkg.EncodeJSON(
			w,
			pkg.Response{Error: "One or more fields are invalid.", Fields: problems},
			http.StatusBadRequest,
		)
		return
	}

	user, err := h.service.Register(
		r.Context(),
		data.Name,
		data.Email,
		data.Password,
	)

	if err != nil {
		if errors.Is(err, ErrEmailAlreadyExists) {
			pkg.EncodeJSON(w, pkg.Response{Error: "user already exists"}, http.StatusConflict)
			return
		}

		slog.Error("internal server error", "error", err)
		pkg.EncodeJSON(
			w,
			pkg.Response{Error: "internal server error"},
			http.StatusInternalServerError,
		)
		return
	}

	token, err := pkg.GenerateAccessToken(user.ID, h.cfg.JWTSecret)
	if err != nil {
		slog.Error("error generating token", "error", err)
		pkg.EncodeJSON(w, pkg.Response{Error: "internal server error"}, http.StatusInternalServerError)
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

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	data, problems, err := pkg.DecodeValidJSON[LoginRequest](w, r)
	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			pkg.EncodeJSON(w, pkg.Response{Error: "request body too large"}, http.StatusRequestEntityTooLarge)
			return
		}
		pkg.EncodeJSON(
			w,
			pkg.Response{Error: "One or more fields are invalid.", Fields: problems},
			http.StatusBadRequest,
		)
		return
	}

	user, err := h.service.Login(
		r.Context(),
		data.Email,
		data.Password,
	)

	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			pkg.EncodeJSON(w, pkg.Response{Error: "password or email is incorrect"}, http.StatusUnauthorized)
			return
		}
		if errors.Is(err, ErrInvalidCredentials) {
			pkg.EncodeJSON(w, pkg.Response{Error: "password or email is incorrect"}, http.StatusUnauthorized)
			return
		}

		slog.Error("internal server error", "error", err)
		pkg.EncodeJSON(w, pkg.Response{Error: "internal server error"}, http.StatusInternalServerError)
		return
	}

	token, err := pkg.GenerateAccessToken(user.ID, h.cfg.JWTSecret)
	if err != nil {
		slog.Error("error generating token", "error", err)
		pkg.EncodeJSON(w, pkg.Response{Error: "internal server error"}, http.StatusInternalServerError)
		return
	}

	resp := HandlerResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Token: token,
	}
	pkg.EncodeJSON(w, pkg.Response{Data: resp}, http.StatusOK)
}
