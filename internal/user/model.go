package user

import "time"

type User struct {
	ID			uint	`gorm:"primaryKey"`
	Email		string	`gorm:"unique;not null"`
	Password	string	`gorm:"not null"`
	RoleID		uint
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
