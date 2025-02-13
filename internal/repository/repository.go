package repository

import (
	"context"
	"database/sql"
	"errors"
	"merch-shop/internal/models"
)

type Repository interface {
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	Add(ctx context.Context, u *models.User) error
}

type PostgresRepository struct {
	DB *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		DB: db,
	}
}

func (r *PostgresRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
	    SELECT id, password_hash, created_at
	    FROM users
	    WHERE username = $1`

	user := &models.User{
		Username: username,
	}

	err := r.DB.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return user, nil
}

// TODO: add Transaction with adding coin account

func (r *PostgresRepository) Add(ctx context.Context, u *models.User) error {
	query := `
	    INSERT INTO users(username, password_hash)
	    VALUES ($1, $2)
	    RETURNING id, created_at`

	args := []any{u.Username, u.PasswordHash}

	err := r.DB.QueryRowContext(ctx, query, args...).Scan(&u.ID, &u.CreatedAt)
	return err
}
