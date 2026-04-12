package auth

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTConfig struct {
	Secret []byte
	TTL time.Duration
}

func GenerateAccessToken(userID uuid.UUID, role string, cfg JWTConfig) (string, error) {
	claims, err := NewClaims(userID, role, cfg.TTL)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(cfg.Secret)
}

func ParseToken(tokenString string, cfg JWTConfig) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}

		return cfg.Secret, nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	return  claims, nil
}
