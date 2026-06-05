package settings

import "github.com/google/uuid"

type UserSettings struct {
	UserID         uuid.UUID `json:"user_id"`
	ProfilePrivacy bool      `json:"profile_privacy"`
	ShowLocation   bool      `json:"show_location"`
	AllowMessages  bool      `json:"allow_messages"`
}
