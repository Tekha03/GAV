package auth

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher struct{}

func (h *PasswordHasher) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hashedPassword), err
}

func (h *PasswordHasher) CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
