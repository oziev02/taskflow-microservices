package task

type EventPublisher interface {
	PublishTaskCreated(task *Task) error
	PublishTaskUpdated(task *Task) error
	PublishTaskDeleted(taskID int64) error
}
