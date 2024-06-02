package app

import (
	"log/slog"

	grpcapp "github.com/4aykovski/grpc_auth_sso/internal/app/grpc"
)

type App struct {
	gRPCApp *grpcapp.App
}

func New(
	log *slog.Logger,
	dSNTemplate string,
	port int,
) *App {

	// TODO: Initialize repositories

	// TODO: Initialize services (usecases)

	gRPCApp := grpcapp.New(
		log,
		port,
	)

	return &App{
		gRPCApp: gRPCApp,
	}
}
