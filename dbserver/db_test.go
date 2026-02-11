package dbserver

import (
	"gav/internal/user"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestInitDB(t *testing.T) {
	var err error

	DB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory DB: %v", err)
	}

	err = DB.AutoMigrate(&user.User{})
	if err != nil {
		t.Fatalf("failed to auto-migrate user: %v", err)
	}

	u := user.User{ ID: 1 }
	if err := DB.Create(&u).Error; err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	var out user.User
	if err := DB.First(&out, 1).Error; err != nil {
		t.Fatalf("failed to read user: %v", err)
	}

	if out.ID != 1 {
		t.Errorf("expected ID 1, got %d", out.ID)
	}
}
