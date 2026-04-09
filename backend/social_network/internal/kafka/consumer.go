package kafka

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type Consumer struct {
}

func NewConsumer(brokers []string, groupID uuid.UUID) error {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	client, err := sarama.NewConsumerGroup(brokers, groupID.String(), config)
	if err != nil {
		return err
	}

	handler := &Handler{}

	for {
		err := client.Consume(nil, []string{"messages"}, handler)
		if err != nil {
			log.Println(err)
		}
	}
}
