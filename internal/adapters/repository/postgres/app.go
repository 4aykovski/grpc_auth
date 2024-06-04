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

type AppRepository struct {
	db *postgres.Db
}

func NewAppRepository(db *postgres.Db) *AppRepository {
	return &AppRepository{
		db: db,
	}
}

func (r *AppRepository) GetApp(ctx context.Context, id int) (entity.App, error) {
	stmt, err := r.db.Prepare("SELECT * FROM apps WHERE id = $1")
	if err != nil {
		return entity.App{}, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var app entity.App
	err = stmt.QueryRowContext(ctx, id).Scan(&app.ID, &app.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.App{}, fmt.Errorf("failed to get app: %w", repository.ErrAppNotFound)
		}

		return entity.App{}, fmt.Errorf("failed to get app: %w", err)
	}

	return app, nil
}
