package auth

import "github.com/google/uuid"

type UserInfo struct {
	ID		uuid.UUID
	Email	string
}
