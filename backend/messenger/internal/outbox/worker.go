package outbox

import (
	"context"
	"encoding/json"
	"log/slog"
	"messenger/internal/kafka"
	"messenger/internal/repository"
	"shared/events"
	"time"

	"github.com/google/uuid"
)

const (
	defaultBatchSize = 100
	defaultInterval  = time.Second
	defaultLease     = 30 * time.Second
	maxRetryDelay    = 5 * time.Minute
)

type Worker struct {
	repo      repository.OutboxRepository
	producer  kafka.EventProducer
	batchSize int
	interval  time.Duration
	lease     time.Duration
}

func NewWorker(repo repository.OutboxRepository, producer kafka.EventProducer) *Worker {
	return &Worker{
		repo:      repo,
		producer:  producer,
		batchSize: defaultBatchSize,
		interval:  defaultInterval,
		lease:     defaultLease,
	}
}

func (w *Worker) Run(ctx context.Context) {
	if w == nil || w.repo == nil || w.producer == nil {
		return
	}

	w.process(ctx)
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.process(ctx)
		}
	}
}

func (w *Worker) process(ctx context.Context) {
	pending, err := w.repo.ClaimPending(ctx, w.batchSize, w.lease)
	if err != nil {
		slog.Error("failed to claim outbox events", "error", err)
		return
	}

	for _, item := range pending {
		if ctx.Err() != nil {
			return
		}

		var event events.Event
		if err := json.Unmarshal(item.Payload, &event); err != nil {
			w.markFailed(ctx, item.ID, item.Attempts, err)
			continue
		}

		if err := w.producer.PublishEvent(event); err != nil {
			w.markFailed(ctx, item.ID, item.Attempts, err)
			continue
		}

		if err := w.repo.MarkPublished(ctx, item.ID); err != nil {
			slog.Error("failed to mark outbox event as published", "event_id", item.ID, "error", err)
		}
	}
}

func (w *Worker) markFailed(ctx context.Context, id uuid.UUID, attempts int, cause error) {
	nextAttemptAt := time.Now().Add(retryDelay(attempts))
	if err := w.repo.MarkFailed(ctx, id, cause, nextAttemptAt); err != nil {
		slog.Error("failed to schedule outbox retry", "event_id", id, "error", err)
	}
}

func retryDelay(attempts int) time.Duration {
	if attempts < 0 {
		attempts = 0
	}
	if attempts > 8 {
		attempts = 8
	}

	delay := time.Second * time.Duration(1<<attempts)
	if delay > maxRetryDelay {
		return maxRetryDelay
	}
	return delay
}
