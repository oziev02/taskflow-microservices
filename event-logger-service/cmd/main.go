package main

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

func main() {
	topics := []string{"task_created", "task_updated", "task_deleted"}

	for _, topic := range topics {
		go consumeTopic(topic)
	}

	select {} // Блокируем main навсегда
}

func consumeTopic(topic string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"kafka:9092"},
		Topic:   topic,
		GroupID: "event-logger-group",
	})

	fmt.Printf("Listening to topic: %s\n", topic)

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("error reading message from topic %s: %v\n", topic, err)
			continue
		}

		log.Printf("Topic: %s | Key: %s | Value: %s\n", topic, string(msg.Key), string(msg.Value))
	}
}
