package user_test

import (
	"testing"

	"gav/user"
)

func TestNewUserStats(t *testing.T) {
	s := user.NewUserStats()

	if s == nil {
		t.Fatalf("NewUserStats returned nil")
	}
	if len(s.Posts) != 0 {
		t.Fatalf("Posts expected length 0, got %d", len(s.Posts))
	}
	if s.PostsCnt != 0 {
		t.Fatalf("PostsCnt expected 0, got %d", s.PostsCnt)
	}
	if s.FollowersCnt != 0 {
		t.Fatalf("FollowersCnt expected 0, got %d", s.FollowersCnt)
	}
}
