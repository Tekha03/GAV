package handlers

import "gav/internal/comment"

type CommentHandler struct {
	service comment.Service
}

func NewCommentHandler(service comment.Service) *CommentHandler {
	return &CommentHandler{service: service}
}
