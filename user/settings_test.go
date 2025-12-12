package user_test

import (
	"testing"
	"time"

	"gav/user"
)

func TestNewUserSettings(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	s := user.NewUserSettings(true, false, "StrongPass123", now)

	if !s.Private {
		t.Fatalf("Private expected true")
	}
	if s.EmailNotification {
		t.Fatalf("EmailNotification expected false")
	}
	if s.Password != "StrongPass123" {
		t.Fatalf("Password expected %q, got %q", "StrongPass123", s.Password)
	}
	if !s.CreatedAt.Equal(now) {
		t.Fatalf("CreatedAt expected %v, got %v", now, s.CreatedAt)
	}
}
