package model

import "github.com/google/uuid"

type Reaction struct {
    ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    MessageID uuid.UUID `gorm:"type:uuid;index"`
    UserID    uuid.UUID `gorm:"type:uuid;index"`
    Emoji     string    `gorm:"type:varchar(10);not null"`
}
