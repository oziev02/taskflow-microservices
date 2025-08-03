package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	clickhouse_writer "github.com/oziev02/taskflow-microservices/event-worker/internal/clickhouse_writer"
	kafka_consumer "github.com/oziev02/taskflow-microservices/event-worker/internal/kafka_consumer"
)

func main() {
	_ = godotenv.Load()

	// ClickHouse setup
	clickhouseAddr := os.Getenv("CLICKHOUSE_ADDR")
	writer := clickhouse_writer.NewClickHouseWriter(clickhouseAddr)

	// Kafka setup
	broker := os.Getenv("KAFKA_BROKER")
	groupID := os.Getenv("KAFKA_GROUP_ID")
	topics := []string{
		os.Getenv("KAFKA_TOPIC_TASK_CREATED"),
		os.Getenv("KAFKA_TOPIC_TASK_UPDATED"),
		os.Getenv("KAFKA_TOPIC_TASK_DELETED"),
	}

	// Start Kafka Consumer
	consumer := kafka_consumer.NewConsumer(broker, groupID, topics, writer)
	log.Println("Starting event-worker...")
	if err := consumer.Start(); err != nil {
		log.Fatalf("event-worker error: %v", err)
	}
}
