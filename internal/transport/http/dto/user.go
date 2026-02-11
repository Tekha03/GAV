package dto

import "gav/internal/user"

type UserResponse struct {
	ID		uint	`json:"id"`
	Email	string	`json:"email"`
}

func NewUserResponse(user *user.User) UserResponse {
	return UserResponse{ID: user.ID, Email: user.Email,}
}
