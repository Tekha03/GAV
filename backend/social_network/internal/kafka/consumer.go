package kafka

import (
	"context"
	"log"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type Consumer struct {
	client  sarama.ConsumerGroup
	handler sarama.ConsumerGroupHandler
	topics  []string
}

func NewConsumer(brokers []string, groupID uuid.UUID, handler sarama.ConsumerGroupHandler) (*Consumer, error) {
	config := sarama.NewConfig()

	client, err := sarama.NewConsumerGroup(brokers, groupID.String(), config)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		client:  client,
		handler: handler,
		topics: []string{
			"chat_events",
			"message_events",
			"reaction_events",
		},
	}, nil
}

func (c *Consumer) Start(ctx context.Context) {
	for {
		if ctx.Err() != nil {
			return
		}

		if err := c.client.Consume(ctx, c.topics, c.handler); err != nil {
			log.Println("Error consuming messages:", err)
		}
	}
}
