package task

type Repository interface {
	Save(task *Task) (int64, error)
	FindByID(id int64) (*Task, error)
	FindAll() ([]*Task, error)
	Update(task *Task) error
	Delete(id int64) error
}
