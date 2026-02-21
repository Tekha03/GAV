package dto

import (
	"gav/internal/post"
	"time"

	"github.com/google/uuid"
)

type PostRequest struct {
	Content	string	`json:"content" validate:"required,max=5000"`
}

type PostResponse struct {
	ID			uuid.UUID		`json:"id"`
	AuthorID	uuid.UUID		`json:"author_id"`
	Content		string			`json:"content"`
	CreatedAt	time.Time		`json:"created_at"`
}

func NewPostResponse(post *post.Post) PostResponse {
	return PostResponse{
		ID:			post.ID,
		AuthorID:	post.UserID,
		Content: post.Content,
		CreatedAt: post.CreatedAt,
	}
}
