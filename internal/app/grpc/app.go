package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	authGRPC "github.com/4aykovski/grpc_auth_sso/internal/adapters/grpc/auth"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(
	log *slog.Logger,
	authService authGRPC.AuthService,
	port int,
) *App {

	gRPCServer := grpc.NewServer()

	authGRPC.Register(gRPCServer, log, authService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic("failed to run gRPC server: " + err.Error())
	}
}

func (a *App) Run() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("failed to listen gRPC: %w", err)
	}

	a.log.Info("starting gRPC server", slog.String("address", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("failed to serve gRPC: %w", err)
	}

	return nil
}

func (a *App) Stop() {
	a.log.Info("stopping gRPC server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}
