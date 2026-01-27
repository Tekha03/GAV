package auth

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

var secret = []byte("CHANGE_ME_SECRET")

func GenerateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":	   time.Now().Add(24 * time.Hour).Unix(),
		"iat":	   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ParseToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return 0, err
	}

	claims := token.Claims.(jwt.MapClaims)
	return uint(claims["user_id"].(float64)), nil
}
