package handlers

import "gav/internal/comment"

type CommentHandler struct {
	service comment.CommentService
}

func NewCommentHandler(service comment.CommentService) *CommentHandler {
	return &CommentHandler{service: service}
}
