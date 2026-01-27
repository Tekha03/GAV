package middleware

import (
	"context"
	"net/http"
	"strings"

	"gav/internal/auth"
)

type contextKey string

const userIDKey contextKey = "user_id"

func Auth(authService *auth.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, reader *http.Request) {
			header := reader.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				http.Error(writer, "unauthorized", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(header, "Bearer ")
			userID, err := auth.ParseToken(token)
			if err != nil {
				http.Error(writer, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(reader.Context(), userIDKey, userID)
			next.ServeHTTP(writer, reader.WithContext(ctx))
		})
	}
}

func UserID(ctx context.Context) (uint, bool) {
	id, ok := ctx.Value(userIDKey).(uint)
	return id, ok
}
