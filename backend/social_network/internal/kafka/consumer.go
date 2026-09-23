package kafka

import (
	"context"
	"log"
	"shared/events"
	"time"

	"github.com/IBM/sarama"
)

type Consumer struct {
	client  sarama.ConsumerGroup
	handler sarama.ConsumerGroupHandler
	topics  []string
}

func (c *Consumer) Close() error {
	if c == nil || c.client == nil {
		return nil
	}
	return c.client.Close()
}

func NewConsumer(brokers []string, groupID string, handler sarama.ConsumerGroupHandler) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	client, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		client:  client,
		handler: handler,
		topics: []string{
			events.ChatTopic,
			events.MessageTopic,
			events.ReactionTopic,
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
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second):
			}
		}
	}
}
