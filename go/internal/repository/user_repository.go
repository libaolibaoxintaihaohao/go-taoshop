package repository

import (
	"context"
	"database/sql"

	"taoshop/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO users (username, email, password_hash)
		VALUES (?, ?, ?)`,
		user.Username, user.Email, user.PasswordHash,
	)
	if err != nil {
		return err
	}

	user.ID, _ = result.LastInsertId()
	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, username, email, password_hash, created_at
		FROM users
		WHERE email = ?`, email)

	var user models.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*models.User, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, username, email, password_hash, created_at
		FROM users
		WHERE id = ?`, id)

	var user models.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt); err != nil {
		return nil, err
	}
	return &user, nil
}
