package repository

import (
	"context"
	"messenger/internal/model"
	"time"

	"github.com/google/uuid"
)

type OutboxRepository interface {
	Create(ctx context.Context, event *model.OutboxEvent) error
	ClaimPending(ctx context.Context, limit int, lease time.Duration) ([]*model.OutboxEvent, error)
	MarkPublished(ctx context.Context, id uuid.UUID) error
	MarkFailed(ctx context.Context, id uuid.UUID, cause error, nextAttemptAt time.Time) error
}
