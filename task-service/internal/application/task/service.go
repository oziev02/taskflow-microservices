package task

import (
	"context"
	//"errors"
	"fmt"

	"github.com/oziev02/taskflow-microservices/task-service/internal/domain/task"
)

type Service struct {
	repo      task.Repository     // PostgreSQL (основной источник)
	publisher task.EventPublisher // Kafka
	cache     task.Cache          // Redis-кэш
}

func NewService(repo task.Repository, publisher task.EventPublisher, cache task.Cache) *Service {
	return &Service{repo: repo, publisher: publisher, cache: cache}
}

func (s *Service) Create(title, description string) (*task.Task, error) {
	newTask := task.NewTask(title, description)

	id, err := s.repo.Save(newTask)
	if err != nil {
		return nil, fmt.Errorf("failed to save task: %w", err)
	}
	newTask.ID = id

	// Кэш: записываем отдельный объект, инвалидируем список
	if s.cache != nil {
		_ = s.cache.SetTask(context.Background(), newTask)
		_ = s.cache.InvalidateTasksList(context.Background())
	}

	// Публикация события
	if err := s.publisher.PublishTaskCreated(newTask); err != nil {
		fmt.Println("failed to publish event:", err)
	}

	return newTask, nil
}

func (s *Service) GetByID(id int64) (*task.Task, error) {
	// Пытаемся достать из кэша
	if s.cache != nil {
		if t, err := s.cache.GetTask(context.Background(), id); err == nil && t != nil {
			return t, nil
		}
	}
	// Иначе — из БД
	t, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}
	// Пишем в кэш
	if s.cache != nil {
		_ = s.cache.SetTask(context.Background(), t)
	}
	return t, nil
}

func (s *Service) GetAll() ([]*task.Task, error) {
	// Пытаемся из кэша
	if s.cache != nil {
		if items, err := s.cache.GetTasks(context.Background()); err == nil && items != nil {
			return items, nil
		}
	}
	// Иначе — из БД
	list, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to load tasks: %w", err)
	}
	// Пишем в кэш
	if s.cache != nil {
		_ = s.cache.SetTasks(context.Background(), list)
	}
	return list, nil
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

	// Кэш: обновляем объект и инвалидируем список
	if s.cache != nil {
		_ = s.cache.SetTask(context.Background(), t)
		_ = s.cache.InvalidateTasksList(context.Background())
	}

	// Публикация события
	_ = s.publisher.PublishTaskUpdated(t)

	return t, nil
}

func (s *Service) Delete(id int64) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	// Кэш: удаляем объект и инвалидируем список
	if s.cache != nil {
		_ = s.cache.DeleteTask(context.Background(), id)
		_ = s.cache.InvalidateTasksList(context.Background())
	}

	// Публикация события
	_ = s.publisher.PublishTaskDeleted(id)
	return nil
}
