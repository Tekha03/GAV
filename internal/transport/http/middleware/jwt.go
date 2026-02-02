package middleware

import (
	"context"
	"gav/internal/auth"
	"net/http"
	"strings"
)

type ctxKey string

const userIDKey ctxKey = "userID"

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		userID, err := auth.ParseToken(token)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserID(ctx context.Context) (uint, bool) {
	id, ok := ctx.Value(userIDKey).(uint)
	return id, ok
}
