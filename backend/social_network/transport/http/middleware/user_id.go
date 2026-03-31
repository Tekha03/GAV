package middleware

import (
	"context"

	"github.com/google/uuid"
)

type ctxKey struct{}

var UserIDKey = ctxKey{}

func UserID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(UserIDKey).(uuid.UUID)
	return id, ok
}
