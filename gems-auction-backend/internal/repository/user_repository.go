package repository

import (
	"context"
	"time"

	"github.com/boswin/gems-auction-backend/config"
	"github.com/boswin/gems-auction-backend/internal/domain"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// Create User
func (r *UserRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (full_name, email, password, role, is_active, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING id
	`

	now := time.Now()

	return config.DB.QueryRow(context.Background(), query,
		user.FullName,
		user.Email,
		user.Password,
		user.Role,
		true,
		now,
		now,
	).Scan(&user.ID)
}

// Get by Email
func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, full_name, email, password, role, is_active, created_at, updated_at
		FROM users WHERE email=$1
	`

	var user domain.User

	err := config.DB.QueryRow(context.Background(), query, email).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
