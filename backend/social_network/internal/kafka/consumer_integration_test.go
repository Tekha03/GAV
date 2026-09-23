package kafka

import (
	"context"
	"encoding/json"
	"os"
	"shared/events"
	"social_network/internal/chat"
	"social_network/internal/message"
	messengersync "social_network/internal/messenger_sync"
	"social_network/internal/reaction"
	"strings"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

func TestKafkaConsumerReceivesSharedChatTopic(t *testing.T) {
	brokers := os.Getenv("GAV_TEST_KAFKA_BROKERS")
	if brokers == "" {
		t.Skip("GAV_TEST_KAFKA_BROKERS is not set")
	}
	store := messengersync.NewStore()
	handler, err := NewHandler(message.NewUseCase(store, nil), chat.NewUseCase(store, nil), reaction.NewUseCase(store, nil))
	if err != nil {
		t.Fatal(err)
	}
	consumer, err := NewConsumer(strings.Split(brokers, ","), "gav-kafka-test-"+uuid.NewString(), handler)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan struct{})
	go func() { consumer.Start(ctx); close(done) }()
	t.Cleanup(func() { cancel(); _ = consumer.Close(); <-done })

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(strings.Split(brokers, ","), config)
	if err != nil {
		t.Fatal(err)
	}
	defer producer.Close()
	chatID, memberID := uuid.New(), uuid.New()
	payload, err := json.Marshal(events.ChatCreatedData{ChatID: chatID, Members: []uuid.UUID{memberID}})
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := json.Marshal(events.Event{EventID: uuid.New(), EventType: events.EventTypeChatCreated, Timestamp: time.Now(), Data: payload})
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := producer.SendMessage(&sarama.ProducerMessage{Topic: events.ChatTopic, Value: sarama.ByteEncoder(encoded)}); err != nil {
		t.Fatal(err)
	}
	deadline := time.After(15 * time.Second)
	for {
		for _, id := range store.ChatMembers(chatID) {
			if id == memberID {
				return
			}
		}
		select {
		case <-deadline:
			t.Fatal("Kafka chat event did not reach social consumer")
		case <-time.After(50 * time.Millisecond):
		}
	}
}
