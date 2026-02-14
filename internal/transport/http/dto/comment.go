package dto

import (
	"github.com/google/uuid"
)

type CreateCommentRequest struct {
	PostID	uuid.UUID	`json:"post_id" validate:"required"`
	Content	string		`json:"content" validate:"required,max=500"`
}
