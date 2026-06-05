package dto

import (
	"time"

	"github.com/google/uuid"
)

type ChatRequestDTO struct {
	IsGroup  bool   `json:"is_group"`
	Title    string `json:"title"`
	PhotoURL string `json:"photo_url"`
}

type ChatResponseDTO struct {
	ID        uuid.UUID `json:"id"`
	IsGroup   bool      `json:"is_group"`
	Title     string    `json:"title"`
	PhotoURL  string    `json:"photo_url"`
	CreatedAt time.Time `json:"created_at"`
}