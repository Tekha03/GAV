package stats

import "github.com/google/uuid"

type UserStats struct {
	UserID		uuid.UUID
	PostCount	uint
	Followers	uint
	Followings	uint
	DogsCount	uint
}
