package profile

import "github.com/google/uuid"

type UserProfile struct {
    UserID          uuid.UUID   `gorm:"primaryKey"`
    Name            string
    Surname         string
    Username        string      `gorm:"unique"`
    ProfilePhotoUrl  string
    Bio             string
    Address         string
    BirthDate       string
}
