package post

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Content   string    `json:"content"`
	ImageUrl  string    `json:"image_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
