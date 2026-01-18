package user_test

import (
	"testing"
	"time"

	"gav/user"
)

func TestNewUserSettings(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	s := user.NewUserSettings(true, false, "gav@mail.com", "StrongPass123", now)

	if !s.Private {
		t.Fatalf("Private expected true")
	}
	if s.EmailNotification {
		t.Fatalf("EmailNotification expected false")
	}
	if s.Email != "gav@mail.com" {
		t.Fatalf("Email expected %q, got %q", "gav@mail.com", s.Email)
	}
	if s.PasswordHash != "StrongPass123" {
		t.Fatalf("Password expected %q, got %q", "StrongPass123", s.PasswordHash)
	}
	if !s.CreatedAt.Equal(now) {
		t.Fatalf("CreatedAt expected %v, got %v", now, s.CreatedAt)
	}
}
