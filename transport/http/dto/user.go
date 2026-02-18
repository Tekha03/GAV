package dto

import (
	"gav/internal/user"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID		uuid.UUID	`json:"id"`
	Email	string	`json:"email"`
}

func NewUserResponse(user *user.User) UserResponse {
	return UserResponse{ID: user.ID, Email: user.Email,}
}
