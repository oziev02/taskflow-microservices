package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/oziev02/taskflow-microservices/task-service/internal/domain/task"
)

// TaskRepository реализует интерфейс task.Repository
type TaskRepository struct {
	db *sqlx.DB
}

func NewTaskRepository(db *sqlx.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Save(t *task.Task) (int64, error) {
	query := `
		INSERT INTO tasks (title, description, is_done, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	var id int64
	err := r.db.QueryRow(
		query,
		t.Title,
		t.Description,
		t.IsDone,
		t.CreatedAt,
		t.UpdatedAt,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert task: %w", err)
	}
	return id, nil
}

func (r *TaskRepository) FindByID(id int64) (*task.Task, error) {
	query := `SELECT id, title, description, is_done, created_at, updated_at FROM tasks WHERE id = $1`

	var t task.Task
	err := r.db.Get(&t, query, id)
	if err != nil {
		return nil, fmt.Errorf("find task: %w", err)
	}
	return &t, nil
}

func (r *TaskRepository) FindAll() ([]*task.Task, error) {
	query := `SELECT id, title, description, is_done, created_at, updated_at FROM tasks ORDER BY created_at DESC`

	var tasks []*task.Task
	err := r.db.Select(&tasks, query)
	if err != nil {
		return nil, fmt.Errorf("find all tasks: %w", err)
	}
	return tasks, nil
}

func (r *TaskRepository) Update(t *task.Task) error {
	query := `
		UPDATE tasks
		SET title = $1, description = $2, is_done = $3, updated_at = $4
		WHERE id = $5
	`
	_, err := r.db.Exec(query, t.Title, t.Description, t.IsDone, t.UpdatedAt, t.ID)
	return err
}

func (r *TaskRepository) Delete(id int64) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
