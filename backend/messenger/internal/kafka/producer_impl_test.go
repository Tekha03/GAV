package kafka

import (
	"shared/events"
	"testing"
)

func TestEventTopicsMatchSharedContract(t *testing.T) {
	cases := []struct {
		event events.EventType
		topic string
	}{
		{events.EventTypeChatCreated, events.ChatTopic},
		{events.EventTypeMessageSent, events.MessageTopic},
		{events.EventTypeReactionAdded, events.ReactionTopic},
	}
	for _, item := range cases {
		got, err := resolveTopic(item.event)
		if err != nil || got != item.topic {
			t.Fatalf("resolveTopic(%q) = %q, %v; want %q", item.event, got, err, item.topic)
		}
	}
}
