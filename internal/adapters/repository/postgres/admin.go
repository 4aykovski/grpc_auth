package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/4aykovski/grpc_auth_sso/internal/adapters/repository"
	"github.com/4aykovski/grpc_auth_sso/internal/entity"
	"github.com/4aykovski/grpc_auth_sso/pkg/database/postgres"
)

type AdminRepository struct {
	db *postgres.Db
}

func NewAdminRepository(db *postgres.Db) *AdminRepository {
	return &AdminRepository{
		db: db,
	}
}

// GetAdmin returns admin by user id
func (r *AdminRepository) GetAdmin(ctx context.Context, userID int) (entity.Admin, error) {
	stmt, err := r.db.Prepare("SELECT * FROM admins WHERE id = $1")
	if err != nil {
		return entity.Admin{}, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var admin entity.Admin
	err = stmt.QueryRowContext(ctx, userID).Scan(&admin.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Admin{}, fmt.Errorf("failed to get admin: %w", repository.ErrUserNotFound)
		}

		return entity.Admin{}, fmt.Errorf("failed to get admin: %w", err)
	}

	return admin, nil
}
