package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID	uuid.UUID	`json:"sub"`
	Role	string		`json:"role"`
	jwt.RegisteredClaims
}

func NewClaims(userID uuid.UUID, role string, ttl time.Duration) (*Claims, error) {
	if userID == uuid.Nil {
		return nil, ErrUserIDNil
	}
	if role == "" {
		return nil, ErrEmptyRole
	}

	return &Claims{
		UserID: userID,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}, nil
}
