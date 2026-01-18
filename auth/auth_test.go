package auth_test

import (
	"gav/auth"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "my_password"
	hashedPassword, err := auth.HashPassword(password)

	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if hashedPassword == password {
		t.Fatalf("HashPassword returned raw password instead of hash")
	}
}

func TestCheckPassword(t *testing.T) {
	password := "secret123"
	hashedPassword, err := auth.HashPassword(password)

	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if !auth.CheckPassword(hashedPassword, password) {
		t.Fatalf("CheckPassword returned false for valid password")
	}
	if auth.CheckPassword("wrongpassword", hashedPassword) {
		t.Fatalf("CheckPassword returned true for invalid password")
	}
}
