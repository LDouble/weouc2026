package logger

import (
	"log/slog"
	"os"

	appconfig "github.com/liangluo/weouc2026/services/api-server/internal/platform/config"
)

func New(cfg appconfig.AppConfig) *slog.Logger {
	level := slog.LevelInfo
	if cfg.Service.Environment == "local" || cfg.Service.Environment == "dev" || cfg.Service.Environment == "test" {
		level = slog.LevelDebug
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	return slog.New(handler).With(
		"service", cfg.Service.Name,
		"environment", cfg.Service.Environment,
		"version", cfg.Service.Version,
	)
}
