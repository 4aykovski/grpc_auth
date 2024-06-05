package app

import (
	"log/slog"
	"time"

	"github.com/4aykovski/grpc_auth_sso/internal/adapters/repository/postgres"
	grpcapp "github.com/4aykovski/grpc_auth_sso/internal/app/grpc"
	"github.com/4aykovski/grpc_auth_sso/internal/service/auth"
	pgDatabase "github.com/4aykovski/grpc_auth_sso/pkg/database/postgres"
	"github.com/4aykovski/grpc_auth_sso/pkg/hasher"
	"github.com/4aykovski/grpc_auth_sso/pkg/manager/secret"
	"github.com/4aykovski/grpc_auth_sso/pkg/manager/token"
)

type App struct {
	GRPCApp *grpcapp.App
}

func New(
	log *slog.Logger,
	dSNTemplate string,
	port int,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
) (*App, error) {

	pgdb, err := pgDatabase.New(dSNTemplate)
	if err != nil {
		return nil, err
	}

	adminRepo := postgres.NewAdminRepository(pgdb)
	userRepo := postgres.NewUserRepository(pgdb)
	appRepo := postgres.NewAppRepository(pgdb)

	tokenManager := &token.Manager{}
	secretManager := &secret.Manager{}
	bcrypt := &hasher.BCrypt{}

	authService := auth.New(log, userRepo, appRepo, adminRepo, tokenManager, secretManager, bcrypt, accessTokenTTL, refreshTokenTTL)

	gRPCApp := grpcapp.New(
		log,
		authService,
		port,
	)

	return &App{
		GRPCApp: gRPCApp,
	}, nil
}
