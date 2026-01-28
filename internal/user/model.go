package user

import "time"

type User struct {
	ID			uint	`gorm:"primaryKey"`
	Email		string	`gorm:"unique;not null"`
	Password	string	`gorm:"not null"`
	RoleID		uint
	CreateAt	time.Time
	UpdatedAt	time.Time
}