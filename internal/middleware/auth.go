package middleware

import (
	"net/http"
	"strings"

	"do_it_back/internal/config"
	"do_it_back/internal/pkg"
)

func AuthMiddleware(
	cfg *config.Config,
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		if r.URL.Path == "/health" || r.URL.Path == "/auth/login" || r.URL.Path == "/auth/register" {
			next.ServeHTTP(w, r)
			return
		}

		auth := r.Header.Get("Authorization")

		if auth == "" {
			pkg.EncodeJSON(
				w,
				pkg.Response{Error: "token not found"},
				http.StatusUnauthorized,
			)

			return
		}

		token := strings.TrimPrefix(
			auth,
			"Bearer ",
		)

		claims, err := pkg.Validate(token, cfg.JWTSecret)

		if err != nil {
			pkg.EncodeJSON(
				w,
				pkg.Response{Error: "invalid token"},
				http.StatusUnauthorized,
			)

			return
		}

		ctx := pkg.ContextWithUserID(r.Context(), claims.UserID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
