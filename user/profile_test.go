package user_test

import (
	"testing"

	"gav/user"
)

func TestNewUserProfile(t *testing.T) {
	p := user.NewUserProfile(
		"Alice",
		"Smith",
		"asmith",
		"alice@mail.com",
		"2000-01-01",
		"photo.png",
		"New York",
		"Hello!",
	)

	if p.Name != "Alice" {
		t.Fatalf("Name expected %q, got %q", "Alice", p.Name)
	}
	if p.Username != "asmith" {
		t.Fatalf("Username expected %q, got %q", "asmith", p.Username)
	}
	if p.Email != "alice@mail.com" {
		t.Fatalf("Email expected %q, got %q", "alice@mail.com", p.Email)
	}
	if p.BirthDate != "2000-01-01" {
		t.Fatalf("BirthDate expected %q, got %q", "2000-01-01", p.BirthDate)
	}
}

