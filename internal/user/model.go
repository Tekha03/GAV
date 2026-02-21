package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID			uuid.UUID	`gorm:"primaryKey"`
	Email		string	`gorm:"unique;not null"`
	Password	string	`gorm:"not null"`
	RoleID		uuid.UUID
	CreatedAt	time.Time
	UpdatedAt	time.Time
}

func NewUser(email, passwordHash string) *User {
    return &User{
        Email: email,
        Password: passwordHash,
        CreatedAt: time.Now(),
    }
}
