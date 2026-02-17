package dto

import "github.com/google/uuid"


type FollowRequest struct {
	UserID uuid.UUID	`json:"user_id" validate:"required"`
}
