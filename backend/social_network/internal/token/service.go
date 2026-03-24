package token

import (
	"context"

	"github.com/google/uuid"
)

type TokenService interface {
	CreateRefresh(ctx context.Context, userID uuid.UUID) (string, error)
	ValidateAndRotate(ctx context.Context, refresh string) (uuid.UUID, string, error)
	Revoke(ctx context.Context, refresh string) error
	RevokeAllForUser(ctx context.Context, userID uuid.UUID) error
}
