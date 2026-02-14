package middleware

import (
	"context"
	"net/http"
	"strings"

	"gav/internal/auth"
	"gav/internal/errors"

	"github.com/google/uuid"
)

type ctxKey struct{}

var userIDKey = ctxKey{}

func JWTAuth(cfg auth.JWTConfig) func(http.Handler) http.Handler {
	return func (next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			header := strings.TrimSpace(r.Header.Get("Authorization"))
			if header == "" {
				ctx := context.WithValue(r.Context(), errors.CodeAuthError, errors.ErrMissingToken.Error())
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			parts := strings.Split(header, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				ctx := context.WithValue(r.Context(), errors.CodeAuthError, errors.ErrInvalidToken.Error())
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			token := strings.TrimSpace(parts[1])
			if token == "" {
				ctx := context.WithValue(r.Context(), errors.CodeAuthError, errors.ErrInvalidToken.Error())
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			userID, err := auth.ParseToken(token, cfg)
			if err != nil {
				ctx := context.WithValue(r.Context(), errors.CodeAuthError, err.Error())
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(userIDKey).(uuid.UUID)
	return id, ok
}
