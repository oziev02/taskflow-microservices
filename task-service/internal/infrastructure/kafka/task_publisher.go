package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/oziev02/taskflow-microservices/task-service/internal/domain/task"
	"github.com/segmentio/kafka-go"
)

type TaskPublisher struct {
	writer *kafka.Writer
	topic  string
}

func NewTaskPublisher(brokerAddr, topic string) *TaskPublisher {
	return &TaskPublisher{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokerAddr),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
		topic: topic,
	}
}

func (p *TaskPublisher) PublishTaskCreated(t *task.Task) error {
	return p.publish("task.created", t)
}

func (p *TaskPublisher) PublishTaskUpdated(t *task.Task) error {
	return p.publish("task.updated", t)
}

func (p *TaskPublisher) PublishTaskDeleted(taskID int64) error {
	return p.publish("task.deleted", map[string]interface{}{"id": taskID})
}

func (p *TaskPublisher) publish(event string, payload interface{}) error {
	data := map[string]interface{}{
		"event":   event,
		"payload": payload,
	}

	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(event),
		Value: b,
	}

	return p.writer.WriteMessages(context.Background(), msg)
}
