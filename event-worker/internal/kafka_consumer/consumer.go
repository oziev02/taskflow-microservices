package kafka_consumer

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type Writer interface {
	HandleEvent(topic string, message []byte) error
}

type Consumer struct {
	reader *kafka.Reader
	writer Writer
	topics []string
}

func NewConsumer(broker, groupID string, topics []string, writer Writer) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     []string{broker},
			GroupID:     groupID,
			GroupTopics: topics,
		}),
		writer: writer,
		topics: topics,
	}
}

func (c *Consumer) Start() error {
	defer c.reader.Close()

	for {
		msg, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("failed to read message: %v", err)
			continue
		}

		log.Printf("received message on topic %s: %s", msg.Topic, string(msg.Value))

		if err := c.writer.HandleEvent(msg.Topic, msg.Value); err != nil {
			log.Printf("failed to handle event: %v", err)
		}
	}
}
