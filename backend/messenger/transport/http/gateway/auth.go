package gateway

import (
	"context"
	"net/http"
	apperrors "shared/app_errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type authContextKey struct{}

var userIDContextKey = authContextKey{}

type authClaims struct {
	UserID uuid.UUID `json:"sub"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

func authMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	secret := []byte(jwtSecret)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			claims, err := parseBearerToken(r.Header.Get("Authorization"), secret)
			if err != nil {
				writeError(w, err)
				return
			}

			ctx := context.WithValue(r.Context(), userIDContextKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func parseBearerToken(header string, secret []byte) (*authClaims, error) {
	header = strings.TrimSpace(header)
	if header == "" {
		return nil, apperrors.New(apperrors.AuthTokenMissing, "authorization token is required")
	}

	parts := strings.Fields(header)
	if len(parts) != 2 || parts[0] != "Bearer" || strings.TrimSpace(parts[1]) == "" {
		return nil, apperrors.New(apperrors.AuthTokenInvalid, "invalid authorization token")
	}

	claims := &authClaims{}
	token, err := jwt.ParseWithClaims(parts[1], claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperrors.New(apperrors.AuthTokenInvalid, "invalid authorization token")
		}

		return secret, nil
	})

	if err != nil {
		return nil, apperrors.Wrap(apperrors.AuthTokenInvalid, "invalid authorization token", err)
	}
	if token == nil || !token.Valid || claims.UserID == uuid.Nil {
		return nil, apperrors.New(apperrors.AuthTokenInvalid, "invalid authorization token")
	}

	return claims, nil
}

func currentUserID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(userIDContextKey).(uuid.UUID)
	return id, ok
}
