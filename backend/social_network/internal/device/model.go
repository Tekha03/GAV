package device

import "github.com/google/uuid"


type DeviceToken struct {
	UserID uuid.UUID `json:"-" gorm:"type:uuid;index"`
	Token  string    `json:"token" gorm:"type:text;index"`
}