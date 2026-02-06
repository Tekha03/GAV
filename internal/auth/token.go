package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID int, cfg JWTConfig) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":	   time.Now().Add(cfg.TTL).Unix(),
		"iat":	   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(cfg.Secret)
}

func ParseToken(tokenString string, cfg JWTConfig) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}

		return cfg.Secret, nil
	})

	if err != nil || !token.Valid {
		return 0, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, ErrInvalidClaims
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, ErrUserIdNotFound
	}

	return uint(userID), nil
}
