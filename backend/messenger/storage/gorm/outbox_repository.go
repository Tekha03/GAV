package gorm

import (
	"context"
	"messenger/internal/model"
	"messenger/internal/repository"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OutboxRepository struct {
	repo *Repository
}

func NewOutboxRepository(repo *Repository) repository.OutboxRepository {
	return &OutboxRepository{repo: repo}
}

func (r *OutboxRepository) Create(ctx context.Context, event *model.OutboxEvent) error {
	return internalError("failed to create outbox event", r.repo.WithContext(ctx).Create(event).Error)
}

func (r *OutboxRepository) ClaimPending(ctx context.Context, limit int, lease time.Duration) ([]*model.OutboxEvent, error) {
	if limit <= 0 {
		return []*model.OutboxEvent{}, nil
	}

	now := time.Now()
	lockedUntil := now.Add(lease)
	var events []*model.OutboxEvent

	err := r.repo.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Where("published_at IS NULL").
			Where("next_attempt_at <= ?", now).
			Where("locked_until IS NULL OR locked_until < ?", now).
			Order("created_at ASC").
			Limit(limit).
			Find(&events).Error; err != nil {
			return err
		}

		if len(events) == 0 {
			return nil
		}

		ids := make([]uuid.UUID, 0, len(events))
		for _, event := range events {
			ids = append(ids, event.ID)
			event.LockedUntil = &lockedUntil
		}

		return tx.Model(&model.OutboxEvent{}).
			Where("id IN ?", ids).
			Update("locked_until", lockedUntil).Error
	})
	if err != nil {
		return nil, internalError("failed to claim pending outbox events", err)
	}

	return events, nil
}

func (r *OutboxRepository) MarkPublished(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	result := r.repo.WithContext(ctx).
		Model(&model.OutboxEvent{}).
		Where("id = ? AND published_at IS NULL", id).
		Updates(map[string]any{
			"published_at": now,
			"locked_until": nil,
			"last_error":   nil,
		})
	return internalError("failed to mark outbox event as published", result.Error)
}

func (r *OutboxRepository) MarkFailed(ctx context.Context, id uuid.UUID, cause error, nextAttemptAt time.Time) error {
	lastError := "unknown error"
	if cause != nil {
		lastError = cause.Error()
	}

	result := r.repo.WithContext(ctx).
		Model(&model.OutboxEvent{}).
		Where("id = ? AND published_at IS NULL", id).
		Updates(map[string]any{
			"attempts":        gorm.Expr("attempts + 1"),
			"next_attempt_at": nextAttemptAt,
			"locked_until":    nil,
			"last_error":      lastError,
		})
	return internalError("failed to mark outbox event as failed", result.Error)
}
