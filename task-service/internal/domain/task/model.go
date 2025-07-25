package task

import "time"

type Task struct {
	ID          int64
	Title       string
	Description string
	IsDone      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewTask(title, description string) *Task {
	now := time.Now()
	return &Task{
		Title:       title,
		Description: description,
		IsDone:      false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (t *Task) MarkDone() {
	t.IsDone = true
	t.UpdatedAt = time.Now()
}

func (t *Task) Update(title, description string) {
	t.Title = title
	t.Description = description
	t.UpdatedAt = time.Now()
}
