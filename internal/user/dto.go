package user

import "github.com/google/uuid"

type UpdateuserInput struct {
	Email		*string
	Password	*string
	RoleID		*uuid.UUID
}
