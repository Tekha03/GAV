package dto

import (
	"gav/internal/post"
	"time"
)

type PostRequest struct {
	Content	string	`json:"content"`
}

type PostResponse struct {
	ID			uint		`json:"id"`
	AuthorID	uint		`json:"author_id"`
	Content		string		`json:"content"`
	CreatedAt	time.Time	`json:"created_at"`
}

func NewPostResponse(post *post.Post) PostResponse {
	return PostResponse{
		ID:			post.ID,
		AuthorID:	post.UserID,
		Content: post.Content,
		CreatedAt: post.CreatedAt,
	}
}
