package user_test

import (
	"testing"
	"time"

	"gav/user"
)

func TestNewUserTest(t *testing.T) {
	profile := user.NewUserProfile("Bob", "Brown", "bb", "bob@example.com", "1990-01-01", "", "", "")
	stats := user.NewUserStats()
	settings := user.NewUserSettings(false, true, "pwd", time.Now())
	dogs := []*user.Dog{
		user.NewDog("Rex", 3, "Shepherd", "Calm", "Male"),
	}

	u := user.NewUser(42, profile, stats, settings, dogs)

	if u.ID != 42 {
		t.Fatalf("ID expected 42, got %d", u.ID)
	}
	if u.Profile.Name != "Bob" {
		t.Fatalf("Profile.Name expected Bob, got %q", u.Profile.Name)
	}
	if u.Stats.PostsCnt != 0 {
		t.Fatalf("Stats.PostsCnt expected 0, got %d", u.Stats.PostsCnt)
	}
	if u.Settings.Password != "pwd" {
		t.Fatalf("Settings.Password expected %q, got %q", "pwd", u.Settings.Password)
	}
	if len(u.UserDogs) != 1 {
		t.Fatalf("UserDogs length expected 1, got %d", len(u.UserDogs))
	}
	if u.UserDogs[0].DogName != "Rex" {
		t.Fatalf("first dog name expected Rex, got %q", u.UserDogs[0].DogName)
	}
}
