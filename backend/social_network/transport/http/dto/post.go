package dto

import (
	"time"

	"github.com/google/uuid"
)

type PostRequest struct {
	Content  string `json:"content" validate:"required,max=5000"`
	ImageUrl string `json:"image_url,omitempty"`
}

type PostResponse struct {
	ID        uuid.UUID `json:"id"`
	AuthorID  uuid.UUID `json:"author_id"`
	Content   string    `json:"content"`
	ImageUrl  string    `json:"image_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
