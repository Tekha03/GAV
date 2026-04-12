package kafka

import (
	"context"
	"log"
	"os"
	"social_network/internal/chat"
	"social_network/internal/message"
	messengersync "social_network/internal/messenger_sync"
	"social_network/internal/notification"
	"social_network/internal/reaction"
	"strings"

	"github.com/google/uuid"
)

func LaunchKafka(ctx context.Context, notifier notification.NotificationService) error {
	store := messengersync.NewStore()
	brokers := loadBrokers()

	messageUC := message.NewUseCase(store, notifier)
	chatUC := chat.NewUseCase(store, notifier)
	reactionUC := reaction.NewUseCase(store, notifier)

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

func loadBrokers() []string {
	value := strings.TrimSpace(os.Getenv("KAFKA_BROKERS"))
	if value == "" {
		return []string{"localhost:9092"}
	}

	parts := strings.Split(value, ",")
	brokers := make([]string, 0, len(parts))
	for _, part := range parts {
		broker := strings.TrimSpace(part)
		if broker != "" {
			brokers = append(brokers, broker)
		}
	}

	if len(brokers) == 0 {
		return []string{"localhost:9092"}
	}

	return brokers
}
