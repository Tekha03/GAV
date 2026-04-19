package kafka

import "shared/events"

type EventProducer interface {
	PublishEvent(event events.Event) error
}
