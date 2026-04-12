package token

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, t *RefreshToken) error
	GetByHash(ctx context.Context, hash string) (*RefreshToken, error)
	Revoke(ctx context.Context, hash string) error
	RevokeAllForUser(ctx context.Context, userID uuid.UUID) error
}
