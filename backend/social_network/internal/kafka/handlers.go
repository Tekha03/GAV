package kafka

import (
	"log"
	"encoding/json"
	"shared/events"

	"github.com/IBM/sarama"
)

type Handler struct {
}

func (h *Handler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var event events.Event
		if err := json.Unmarshal(message.Value, &event); err != nil {
			continue
		}

		switch event.EventType {
		case events.EventTypeMessageSent:
			log.Println("received message_sent event")
		}

		session.MarkMessage(message, "")
	}

	return nil
}
