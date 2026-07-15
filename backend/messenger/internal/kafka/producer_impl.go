package kafka

import (
	"encoding/json"
	"fmt"
	"shared/events"

	"github.com/IBM/sarama"
)

const (
	chatTopic     = "chat-events"
	messageTopic  = "message-events"
	reactionTopic = "reaction-events"
)

type Producer struct {
	producer sarama.SyncProducer
}

func NewProducer(brokers []string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{producer: producer}, nil
}

func (p *Producer) PublishEvent(event events.Event) error {
	if p == nil || p.producer == nil {
		return nil
	}

	bytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	topic, err := resolveTopic(event.EventType)
	if err != nil {
		return err
	}

	message := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(event.EventID.String()),
		Value: sarama.ByteEncoder(bytes),
	}

	_, _, err = p.producer.SendMessage(message)
	return err
}

func (p *Producer) Close() error {
	if p == nil || p.producer == nil {
		return nil
	}
	return p.producer.Close()
}

func resolveTopic(eventType events.EventType) (string, error) {
	switch eventType {
	case events.EventTypeChatCreated,
		events.EventTypeChatMemberAdded,
		events.EventTypeChatMemberRemoved,
		events.EventTypeChatUpdated,
		events.EventTypeChatDeleted:
		return chatTopic, nil

	case events.EventTypeMessageSent,
		events.EventTypeMessageEdited,
		events.EventTypeMessageDeleted:
		return messageTopic, nil

	case events.EventTypeReactionAdded,
		events.EventTypeReactionRemoved:
		return reactionTopic, nil

	default:
		return "", fmt.Errorf("unknown event type: %s", eventType)
	}
}
