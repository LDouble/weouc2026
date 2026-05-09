package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/liangluo/weouc2026/services/api-server/internal/modules/system"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/auth"
	appconfig "github.com/liangluo/weouc2026/services/api-server/internal/platform/config"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/httpx"
)

type App struct {
	config appconfig.AppConfig
	logger *slog.Logger
	router *gin.Engine
}

func New(cfg appconfig.AppConfig, logger *slog.Logger) *App {
	return &App{
		config: cfg,
		logger: logger,
		router: NewRouter(cfg, logger),
	}
}

func NewRouter(cfg appconfig.AppConfig, logger *slog.Logger) *gin.Engine {
	gin.SetMode(resolveGinMode(cfg.Service.Environment))

	engine := gin.New()
	engine.Use(
		httpx.RequestIDMiddleware(),
		httpx.AccessLogMiddleware(logger),
		httpx.RecoveryMiddleware(logger),
		auth.ContextMiddleware(cfg),
	)

	system.NewModule(cfg).RegisterRoutes(engine)
	return engine
}

func (a *App) Run(ctx context.Context) error {
	server := &http.Server{
		Addr:              a.config.Server.Address(),
		Handler:           a.router,
		ReadHeaderTimeout: a.config.Server.ReadTimeout,
		ReadTimeout:       a.config.Server.ReadTimeout,
		WriteTimeout:      a.config.Server.WriteTimeout,
	}

	a.logger.Info("api server starting", "address", server.Addr)

	errCh := make(chan error, 1)
	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		a.logger.Info("api server shutting down")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), a.config.Server.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return err
		}

		return <-errCh
	}
}

func resolveGinMode(environment string) string {
	switch environment {
	case "prod", "production":
		return gin.ReleaseMode
	case "test":
		return gin.TestMode
	default:
		return gin.DebugMode
	}
}
