package task

import (
	"fmt"

	"github.com/oziev02/taskflow-microservices/task-service/internal/domain/task"
)

type Service struct {
	repo      task.Repository     // интерфейс доступа к данным (PostgreSQL/Redis)
	publisher task.EventPublisher // интерфейс публикации событий (Kafka)
}

func NewService(repo task.Repository, publisher task.EventPublisher) *Service {
	return &Service{repo: repo, publisher: publisher}
}

func (s *Service) Create(title, description string) (*task.Task, error) {
	newTask := task.NewTask(title, description)

	id, err := s.repo.Save(newTask)
	if err != nil {
		return nil, fmt.Errorf("failed to save task: %w", err)
	}
	newTask.ID = id

	// Публикуем событие в Kafka
	err = s.publisher.PublishTaskCreated(newTask)
	if err != nil {
		fmt.Println("failed to publish event:", err)
	}

	return newTask, nil
}

func (s *Service) GetByID(id int64) (*task.Task, error) {
	return s.repo.FindByID(id)
}

func (s *Service) GetAll() ([]*task.Task, error) {
	return s.repo.FindAll()
}

func (s *Service) Update(id int64, title, description string) (*task.Task, error) {
	t, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	t.Update(title, description)
	if err := s.repo.Update(t); err != nil {
		return nil, err
	}

	// Публикация события
	_ = s.publisher.PublishTaskUpdated(t)

	return t, nil
}

func (s *Service) Delete(id int64) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	// Публикация события
	_ = s.publisher.PublishTaskDeleted(id)
	return nil
}
