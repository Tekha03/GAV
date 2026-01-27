package dto

type CreateCommentRequest struct {
	PostID	uint   `json:"post_id" validate:"required"`
	Content string `json:"content" validate:"required,min=1,max=500"`
}
