package kafka

import (
	"testing"
	"time"

	"shared/events"

	"github.com/google/uuid"
)

func TestResolveTopic(t *testing.T) {
	cases := map[events.EventType]string{
		events.EventTypeChatCreated:       chatTopic,
		events.EventTypeChatMemberAdded:   chatTopic,
		events.EventTypeChatMemberRemoved: chatTopic,
		events.EventTypeChatUpdated:       chatTopic,
		events.EventTypeChatDeleted:       chatTopic,
		events.EventTypeMessageSent:       messageTopic,
		events.EventTypeMessageEdited:     messageTopic,
		events.EventTypeMessageDeleted:    messageTopic,
		events.EventTypeReactionAdded:     reactionTopic,
		events.EventTypeReactionRemoved:   reactionTopic,
	}

	for eventType, expectedTopic := range cases {
		topic, err := resolveTopic(eventType)
		if err != nil {
			t.Fatalf("resolveTopic returned error for %s: %v", eventType, err)
		}
		if topic != expectedTopic {
			t.Fatalf("unexpected topic for %s: %s", eventType, topic)
		}
	}
}

func TestResolveTopicUnknown(t *testing.T) {
	if _, err := resolveTopic(events.EventType("unknown")); err == nil {
		t.Fatal("expected error for unknown event type")
	}
}

func TestProducerHelpers(t *testing.T) {
	var nilProducer *Producer
	if err := nilProducer.Close(); err != nil {
		t.Fatalf("Close on nil producer returned error: %v", err)
	}

	producer := &Producer{}
	if err := producer.Close(); err != nil {
		t.Fatalf("Close on producer with nil inner producer returned error: %v", err)
	}

	err := producer.PublishEvent(events.Event{
		EventID:   uuid.New(),
		EventType: events.EventType("missing"),
		Timestamp: time.Now(),
		Data:      []byte(`{"ok":true}`),
	})
	if err == nil {
		t.Fatal("expected error for unsupported event type")
	}
}
