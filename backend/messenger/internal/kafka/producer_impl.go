package kafka

import (
	"encoding/json"
	apperrors "shared/app_errors"
	"shared/events"

	"github.com/IBM/sarama"
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
		return nil, apperrors.Wrap(apperrors.ServiceUnavailable, "failed to create Kafka producer", err)
	}

	return &Producer{producer: producer}, nil
}

func (p *Producer) PublishEvent(event events.Event) error {
	if p == nil || p.producer == nil {
		return nil
	}

	bytes, err := json.Marshal(event)
	if err != nil {
		return apperrors.Wrap(apperrors.Internal, "failed to encode event", err)
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
	return apperrors.Wrap(apperrors.ServiceUnavailable, "failed to publish event", err)
}

func (p *Producer) Close() error {
	if p == nil || p.producer == nil {
		return nil
	}
	return apperrors.Wrap(apperrors.Internal, "failed to close Kafka producer", p.producer.Close())
}

func resolveTopic(eventType events.EventType) (string, error) {
	switch eventType {
	case events.EventTypeChatCreated,
		events.EventTypeChatMemberAdded,
		events.EventTypeChatMemberRemoved,
		events.EventTypeChatUpdated,
		events.EventTypeChatDeleted:
		return events.ChatTopic, nil

	case events.EventTypeMessageSent,
		events.EventTypeMessageEdited,
		events.EventTypeMessageDeleted:
		return events.MessageTopic, nil

	case events.EventTypeReactionAdded,
		events.EventTypeReactionRemoved:
		return events.ReactionTopic, nil

	default:
		return "", apperrors.New(apperrors.Internal, "unknown event type", apperrors.WithDetail("event_type", eventType))
	}
}
