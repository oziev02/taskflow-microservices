package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/oziev02/taskflow-microservices/auth-service/internal/domain/user"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(email, passwordHash string) (int64, error) {
	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id`
	var id int64
	if err := r.db.QueryRow(query, email, passwordHash).Scan(&id); err != nil {
		return 0, fmt.Errorf("insert user: %w", err)
	}
	return id, nil
}

func (r *UserRepository) FindByEmail(email string) (*user.User, error) {
	query := `SELECT id, email, password_hash, created_at FROM users WHERE email = $1`
	var u user.User
	if err := r.db.Get(&u, query, email); err != nil {
		return nil, fmt.Errorf("find by email: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) FindByID(id int64) (*user.User, error) {
	query := `SELECT id, email, password_hash, created_at FROM users WHERE id = $1`
	var u user.User
	if err := r.db.Get(&u, query, id); err != nil {
		return nil, fmt.Errorf("find by id: %w", err)
	}
	return &u, nil
}
