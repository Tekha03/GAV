package user

import "github.com/google/uuid"

type UpdateUserInput struct {
	Email		*string
	Password	*string
	RoleID		*uuid.UUID
}
