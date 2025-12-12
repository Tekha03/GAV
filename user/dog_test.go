package user_test

import (
	"testing"

	"gav/user"
)

func TestNewDog(t *testing.T) {
	d := user.NewDog("Sharik", 5, "Labrador", "Friendly", "Male")

	if d.DogName != "Sharik" {
		t.Fatalf("DogName: expected %q, got %q", "Sharik", d.DogName)
	}
	if d.Age != 5 {
		t.Fatalf("Age: expected %d, got %d", 5, d.Age)
	}
	if d.Breed != "Labrador" {
		t.Fatalf("Breed: expected %q, got %q", "Labrador", d.Breed)
	}
	if d.Character != "Friendly" {
		t.Fatalf("Character: expected %q, got %q", "Friendly", d.Character)
	}
	if d.Gender != "Male" {
		t.Fatalf("Gender: expected %q, got %q", "Male", d.Gender)
	}
}
