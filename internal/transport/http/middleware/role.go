package middleware

import (
	"errors"
	"gav/internal/transport/response"
	"gav/internal/user"
	"net/http"
)

func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			user, ok := r.Context().Value("user").(*user.User)
			if !ok || user.RoleID == 0 {
				response.Error(w, http.StatusForbidden, errors.New("forbidden"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
