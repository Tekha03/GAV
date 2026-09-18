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
)

func LaunchKafka(ctx context.Context, notifier notification.NotificationService) error {
	store := messengersync.NewStore()
	brokers := loadBrokers()

	// New-message notifications are delivered directly by messenger over gRPC.
	// Repeating them from Kafka would create a second notification for each message.
	messageUC := message.NewUseCase(store, nil)
	// The current projection is in-memory and cannot deduplicate replayed events.
	// Do not emit user notifications from it until a durable idempotent handler exists.
	chatUC := chat.NewUseCase(store, nil)
	reactionUC := reaction.NewUseCase(store, nil)

	handler, err := NewHandler(messageUC, chatUC, reactionUC)
	if err != nil {
		return err
	}

	groupID := strings.TrimSpace(os.Getenv("KAFKA_CONSUMER_GROUP"))
	if groupID == "" {
		groupID = "gav-social-network"
	}
	consumer, err := NewConsumer(brokers, groupID, handler)
	if err != nil {
		return err
	}

	go consumer.Start(ctx)
	go func() {
		<-ctx.Done()
		if err := consumer.Close(); err != nil {
			log.Printf("Kafka consumer close error: %v", err)
		}
	}()

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
