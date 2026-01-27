package user

type User struct {
	ID			uint	`gorm:"primaryKey"`
	Email		string	`gorm:"unique;not null"`
	Password	string	`gorm:"not null"`
	RoleID		uint
}


