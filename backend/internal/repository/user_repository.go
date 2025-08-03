package repository

import (
	"database/sql"

	"github.com/xarcher/backend/internal/domain"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	query := `INSERT INTO users (username, password, created_at) VALUES ($1, $2, $3) RETURNING id`
	return r.db.QueryRow(query, user.Username, user.Password, user.CreatedAt).Scan(&user.ID)
}

func (r *userRepository) GetByUsername(username string) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT id, username, password, created_at FROM users WHERE username = $1`
	err := r.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetByID(id int) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT id, username, password, created_at FROM users WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}
