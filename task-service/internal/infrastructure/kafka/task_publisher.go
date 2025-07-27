package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/oziev02/taskflow-microservices/task-service/internal/domain/task"
	"github.com/segmentio/kafka-go"
)

type TaskPublisher struct {
	createdWriter *kafka.Writer
	updatedWriter *kafka.Writer
	deletedWriter *kafka.Writer
}

func NewTaskPublisher(broker, createdTopic, updatedTopic, deletedTopic string) *TaskPublisher {
	return &TaskPublisher{
		createdWriter: newWriter(broker, createdTopic),
		updatedWriter: newWriter(broker, updatedTopic),
		deletedWriter: newWriter(broker, deletedTopic),
	}
}

func newWriter(broker, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func (p *TaskPublisher) PublishTaskCreated(t *task.Task) error {
	return p.publish(p.createdWriter, "task_created", t)
}

func (p *TaskPublisher) PublishTaskUpdated(t *task.Task) error {
	return p.publish(p.updatedWriter, "task_updated", t)
}

func (p *TaskPublisher) PublishTaskDeleted(taskID int64) error {
	return p.publish(p.deletedWriter, "task_deleted", map[string]interface{}{"id": taskID})
}

func (p *TaskPublisher) publish(writer *kafka.Writer, event string, payload interface{}) error {
	data := map[string]interface{}{
		"event":   event,
		"payload": payload,
	}
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}
	msg := kafka.Message{
		Key:   []byte(event),
		Value: b,
	}
	return writer.WriteMessages(context.Background(), msg)
}
