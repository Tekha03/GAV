package dto

import (
	"gav/internal/chat/model"

	"github.com/google/uuid"
)

type AttachmentRequest struct {
	URL      string               `json:"url"`
	Type     model.AttachmentType `json:"type"`
	FileName string               `json:"file_name"`
	FileSize int64                `json:"file_size"`
}

type AttachmentResponse struct {
	ID        uuid.UUID           	`json:"id"`
	MessageID uuid.UUID           	`json:"message_id"`
	URL       string              	`json:"url"`
	Type      model.AttachmentType 	`json:"type"`
	FileName  string              	`json:"file_name"`
	FileSize  int64               	`json:"file_size"`
}