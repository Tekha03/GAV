package handlers

import "gav/internal/comment"

type CommentHandler struct {
	service comment.CommentService
}

func  NewCommentHandler(service comment.CommentService) *CommentHandler {
	return &CommentHandler{service: service}
}

type CreateCommentRequest struct {
	PostID	uint	`json:"post_id" validate:"required"`
	Content	string	`json:"content" validate:"required,min=1,max=500"`
}
