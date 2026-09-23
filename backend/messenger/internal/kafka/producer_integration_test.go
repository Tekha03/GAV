package kafka

import (
	"encoding/json"
	"os"
	"shared/events"
	"strings"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

func TestProducerPublishesToSharedMessageTopic(t *testing.T) {
	brokers := os.Getenv("GAV_TEST_KAFKA_BROKERS")
	if brokers == "" {
		t.Skip("GAV_TEST_KAFKA_BROKERS is not set")
	}
	addresses := strings.Split(brokers, ",")
	client, err := sarama.NewClient(addresses, sarama.NewConfig())
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	offset, err := client.GetOffset(events.MessageTopic, 0, sarama.OffsetNewest)
	if err != nil {
		t.Fatal(err)
	}
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		t.Fatal(err)
	}
	defer consumer.Close()
	partition, err := consumer.ConsumePartition(events.MessageTopic, 0, offset)
	if err != nil {
		t.Fatal(err)
	}
	defer partition.Close()

	producer, err := NewProducer(addresses)
	if err != nil {
		t.Fatal(err)
	}
	defer producer.Close()
	messageID := uuid.New()
	payload, err := json.Marshal(events.MessageSentData{MessageID: messageID, ChatID: uuid.New(), SenderID: uuid.New(), Text: "integration"})
	if err != nil {
		t.Fatal(err)
	}
	if err := producer.PublishEvent(events.Event{EventID: uuid.New(), EventType: events.EventTypeMessageSent, Timestamp: time.Now(), Data: payload}); err != nil {
		t.Fatal(err)
	}
	select {
	case item := <-partition.Messages():
		var got events.Event
		if err := json.Unmarshal(item.Value, &got); err != nil {
			t.Fatal(err)
		}
		if got.EventType != events.EventTypeMessageSent {
			t.Fatalf("event type = %q", got.EventType)
		}
		var data events.MessageSentData
		if err := json.Unmarshal(got.Data, &data); err != nil {
			t.Fatal(err)
		}
		if data.MessageID != messageID {
			t.Fatalf("message ID = %s, want %s", data.MessageID, messageID)
		}
	case <-time.After(15 * time.Second):
		t.Fatal("producer event did not reach Kafka")
	}
}
