package kafka

import (
	"context"
	"log"
	"social_network/internal/chat"
	"social_network/internal/message"
	"social_network/internal/reaction"

	"github.com/google/uuid"
)

func LaunchKafka(ctx context.Context, cancel context.CancelFunc) error {
	brokers := []string{"localhost:9092"}

	messageUC := message.NewUseCase()
	chatUC := chat.NewUseCase()
	reactionUC := reaction.NewUseCase()

	handler, err := NewHandler(messageUC, chatUC, reactionUC)
	if err != nil {
		return err
	}

	consumer, err := NewConsumer(brokers, uuid.New(), handler)
	if err != nil {
		return err
	}

	go consumer.Start(ctx)

	log.Println("Kafka consumer started")
	return nil
}
