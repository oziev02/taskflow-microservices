package user

type Repository interface {
	Create(email, passwordHash string) (int64, error)
	FindByEmail(email string) (*User, error)
	FindByID(id int64) (*User, error)
}
