package outbox

import (
	"context"
	"encoding/json"
	"errors"
	"messenger/internal/model"
	"shared/events"
	"testing"
	"time"

	"github.com/google/uuid"
)

type fakeRepository struct {
	pending       []*model.OutboxEvent
	publishedIDs  []uuid.UUID
	failedIDs     []uuid.UUID
	nextAttemptAt time.Time
}

func (r *fakeRepository) Create(context.Context, *model.OutboxEvent) error {
	return nil
}

func (r *fakeRepository) ClaimPending(context.Context, int, time.Duration) ([]*model.OutboxEvent, error) {
	return r.pending, nil
}

func (r *fakeRepository) MarkPublished(_ context.Context, id uuid.UUID) error {
	r.publishedIDs = append(r.publishedIDs, id)
	return nil
}

func (r *fakeRepository) MarkFailed(_ context.Context, id uuid.UUID, _ error, nextAttemptAt time.Time) error {
	r.failedIDs = append(r.failedIDs, id)
	r.nextAttemptAt = nextAttemptAt
	return nil
}

type fakeProducer struct {
	err       error
	published []events.Event
}

func (p *fakeProducer) PublishEvent(event events.Event) error {
	if p.err != nil {
		return p.err
	}
	p.published = append(p.published, event)
	return nil
}

func TestWorkerMarksPublishedEvent(t *testing.T) {
	event := events.Event{
		EventID:   uuid.New(),
		EventType: events.EventTypeMessageSent,
		Timestamp: time.Now(),
		Data:      json.RawMessage(`{"message_id":"00000000-0000-0000-0000-000000000001"}`),
	}
	payload, err := json.Marshal(event)
	if err != nil {
		t.Fatal(err)
	}

	repo := &fakeRepository{pending: []*model.OutboxEvent{{ID: event.EventID, Payload: payload}}}
	producer := &fakeProducer{}
	worker := NewWorker(repo, producer)

	worker.process(context.Background())

	if len(producer.published) != 1 || producer.published[0].EventID != event.EventID {
		t.Fatalf("published events = %#v, want event %s", producer.published, event.EventID)
	}
	if len(repo.publishedIDs) != 1 || repo.publishedIDs[0] != event.EventID {
		t.Fatalf("published IDs = %v, want %s", repo.publishedIDs, event.EventID)
	}
	if len(repo.failedIDs) != 0 {
		t.Fatalf("failed IDs = %v, want none", repo.failedIDs)
	}
}

func TestWorkerSchedulesRetryWhenPublishFails(t *testing.T) {
	event := events.Event{
		EventID:   uuid.New(),
		EventType: events.EventTypeMessageSent,
		Timestamp: time.Now(),
		Data:      json.RawMessage(`{}`),
	}
	payload, err := json.Marshal(event)
	if err != nil {
		t.Fatal(err)
	}

	repo := &fakeRepository{pending: []*model.OutboxEvent{{ID: event.EventID, Payload: payload}}}
	producer := &fakeProducer{err: errors.New("kafka unavailable")}
	worker := NewWorker(repo, producer)
	startedAt := time.Now()

	worker.process(context.Background())

	if len(repo.publishedIDs) != 0 {
		t.Fatalf("published IDs = %v, want none", repo.publishedIDs)
	}
	if len(repo.failedIDs) != 1 || repo.failedIDs[0] != event.EventID {
		t.Fatalf("failed IDs = %v, want %s", repo.failedIDs, event.EventID)
	}
	if !repo.nextAttemptAt.After(startedAt) {
		t.Fatalf("next attempt = %s, want after %s", repo.nextAttemptAt, startedAt)
	}
}
