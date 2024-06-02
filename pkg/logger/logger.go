package logger

import (
	"log/slog"
	"os"

	"github.com/natefinch/lumberjack"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func InitLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level:     slog.LevelDebug,
				AddSource: true,
			}),
		)
	case envDev:
		lumber := &lumberjack.Logger{
			Filename:   "logs/app.log",
			MaxSize:    10,
			MaxBackups: 3,
			MaxAge:     7,
		}
		log = slog.New(
			slog.NewJSONHandler(lumber, &slog.HandlerOptions{
				Level:     slog.LevelDebug,
				AddSource: true,
			}),
		)
	case envProd:
		lumber := &lumberjack.Logger{
			Filename:   "logs/app.log",
			MaxSize:    10,
			MaxBackups: 3,
			MaxAge:     7,
		}
		log = slog.New(
			slog.NewJSONHandler(lumber, &slog.HandlerOptions{
				Level:     slog.LevelInfo,
				AddSource: true,
			}),
		)
	}
	return log
}
