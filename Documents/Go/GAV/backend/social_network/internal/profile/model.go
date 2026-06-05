package profile

import "github.com/google/uuid"

type UserProfile struct {
	UserID          uuid.UUID `gorm:"primaryKey" json:"user_id"`
	Name            string    `json:"name"`
	Surname         string    `json:"surname"`
	Username        string    `gorm:"uniqueIndex" json:"username"`
	ProfilePhotoUrl string    `json:"profile_photo_url"`
	Bio             string    `json:"bio"`
	Address         string    `json:"address"`
	BirthDate       string    `json:"birth_date"`
}
