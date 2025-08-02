package task

import (
	"context"
)

type Cache interface {
	// Одна задача
	GetTask(ctx context.Context, id int64) (*Task, error)
	SetTask(ctx context.Context, t *Task) error
	DeleteTask(ctx context.Context, id int64) error

	// Список задач
	GetTasks(ctx context.Context) ([]*Task, error)
	SetTasks(ctx context.Context, list []*Task) error
	InvalidateTasksList(ctx context.Context) error
}
