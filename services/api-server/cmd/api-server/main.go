package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/liangluo/weouc2026/services/api-server/internal/platform/app"
	appconfig "github.com/liangluo/weouc2026/services/api-server/internal/platform/config"
	applogger "github.com/liangluo/weouc2026/services/api-server/internal/platform/logger"
)

func main() {
	cfg, err := appconfig.Load()
	if err != nil {
		slog.Error("load config failed", "error", err)
		os.Exit(1)
	}

	logger := applogger.New(cfg)
	application, err := app.New(cfg, logger)
	if err != nil {
		logger.Error("build api server failed", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := application.Run(ctx); err != nil {
		logger.Error("api server exited with error", "error", err)
		os.Exit(1)
	}
}
