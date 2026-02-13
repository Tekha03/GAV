package settings

import "github.com/google/uuid"

type UserSettings struct {
    UserID         uuid.UUID
    ProfilePrivacy  bool
    ShowLocation   bool
    AllowMessages  bool
}
