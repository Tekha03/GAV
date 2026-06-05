package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"primaryKey"`
	Email    string    `gorm:"unique;not null"`
	Password string    `gorm:"not null"`
	Role     string    `gorm:"type:text;not null;default:'user'"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Lat            *float64
	Lon            *float64
	LocationStatus LocationStatus
	Visibility     LocationVisibility
}

func NewUser(email, passwordHash string) (*User, error) {
	if email == "" {
		return nil, ErrEmailEmpty
	}
	if passwordHash == "" {
		return nil, ErrPasswordHashEmpty
	}

	return &User{
		ID:        uuid.New(),
		Email:     email,
		Password:  passwordHash,
		Role:      "user",
		CreatedAt: time.Now(),
	}, nil
}
