package gateway

import (
	"context"
	"net/http"
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
				writeError(w, http.StatusUnauthorized, err)
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
		return nil, errMissingToken
	}

	parts := strings.Fields(header)
	if len(parts) != 2 || parts[0] != "Bearer" || strings.TrimSpace(parts[1]) == "" {
		return nil, errInvalidToken
	}

	claims := &authClaims{}
	token, err := jwt.ParseWithClaims(parts[1], claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errUnexpectedSigningMethod
		}

		return secret, nil
	})

	if err != nil || !token.Valid || claims.UserID == uuid.Nil {
		return nil, errInvalidToken
	}

	return claims, nil
}

func currentUserID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(userIDContextKey).(uuid.UUID)
	return id, ok
}
