package redisCache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/oziev02/taskflow-microservices/task-service/internal/domain/task"
)

type TaskCache struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewTaskCache(addr string, db int, ttl time.Duration) *TaskCache {
	return &TaskCache{
		rdb: redis.NewClient(&redis.Options{
			Addr: addr,
			DB:   db,
		}),
		ttl: ttl,
	}
}

func keyTask(id int64) string { return fmt.Sprintf("task:%d", id) }
func keyTasksAll() string     { return "tasks:all" }

func (c *TaskCache) GetTask(ctx context.Context, id int64) (*task.Task, error) {
	b, err := c.rdb.Get(ctx, keyTask(id)).Bytes()
	if err != nil {
		return nil, err
	}
	var t task.Task
	if err := json.Unmarshal(b, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

func (c *TaskCache) SetTask(ctx context.Context, t *task.Task) error {
	b, err := json.Marshal(t)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, keyTask(t.ID), b, c.ttl).Err()
}

func (c *TaskCache) DeleteTask(ctx context.Context, id int64) error {
	return c.rdb.Del(ctx, keyTask(id)).Err()
}

func (c *TaskCache) GetTasks(ctx context.Context) ([]*task.Task, error) {
	b, err := c.rdb.Get(ctx, keyTasksAll()).Bytes()
	if err != nil {
		return nil, err
	}
	var list []*task.Task
	if err := json.Unmarshal(b, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (c *TaskCache) SetTasks(ctx context.Context, list []*task.Task) error {
	b, err := json.Marshal(list)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, keyTasksAll(), b, c.ttl).Err()
}

func (c *TaskCache) InvalidateTasksList(ctx context.Context) error {
	return c.rdb.Del(ctx, keyTasksAll()).Err()
}
