package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/4aykovski/grpc_auth_sso/internal/app"
	"github.com/4aykovski/grpc_auth_sso/internal/config"
	"github.com/4aykovski/grpc_auth_sso/pkg/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.InitLogger(cfg.Env)

	log.Info("Starting sso service", slog.String("env", cfg.Env))
	log.Debug("Tokens TTL", slog.Duration("access_token_ttl", cfg.AccessTokenTtl), slog.Duration("refresh_token_ttl", cfg.RefreshTokenTtl))
	log.Debug("GRPC Configuration", slog.String("host", cfg.GRPC.Host), slog.Int("port", cfg.GRPC.Port), slog.Duration("timeout", cfg.GRPC.Timeout))
	log.Debug("Postgres Configuration", slog.String("host", cfg.Postgres.Host), slog.Int("port", cfg.Postgres.Port), slog.String("database", cfg.Postgres.Database))

	application := app.New(
		log,
		cfg.GRPC.Port,
		cfg.Postgres.DSNTemplate,
		cfg.AccessTokenTtl,
		cfg.RefreshTokenTtl,
	)

	go application.GRPCServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop
	log.Info("Stopping sso service", slog.String("signal", sign.String()))
	application.GRPCServer.Stop()
	log.Info("sso service stopped")
}
