package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/4aykovski/grpc_auth_sso/internal/adapters/repository"
	"github.com/4aykovski/grpc_auth_sso/internal/entity"
	"github.com/4aykovski/grpc_auth_sso/pkg/database/postgres"
	"github.com/lib/pq"
)

type UserRepository struct {
	db *postgres.Db
}

func NewUserRepository(db *postgres.Db) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// SaveUser saves user to the database
func (r *UserRepository) SaveUser(ctx context.Context, user entity.User) (int64, error) {

	stmt, err := r.db.Prepare("INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id")
	if err != nil {
		return -1, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var id int64
	err = stmt.QueryRowContext(ctx, user.Email, user.PasswordHash).Scan(&id)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code.Name() == "unique_violation" {
			return -1, fmt.Errorf("failed to save user: %w", repository.ErrUserAlreadyExists)
		}

		return -1, fmt.Errorf("failed to save user: %w", err)
	}

	return id, nil
}

// GetUser returns user by email
func (r *UserRepository) GetUser(ctx context.Context, email string) (entity.User, error) {
	stmt, err := r.db.Prepare("SELECT * FROM users WHERE email = $1")
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var user entity.User
	err = stmt.QueryRowContext(ctx, email).Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, fmt.Errorf("failed to get user: %w", repository.ErrUserNotFound)
		}

		return entity.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
