package app

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/liangluo/weouc2026/services/api-server/internal/modules/campus_life"
	clrepo "github.com/liangluo/weouc2026/services/api-server/internal/modules/campus_life/repo"
	"github.com/liangluo/weouc2026/services/api-server/internal/modules/iam"
	iamrepo "github.com/liangluo/weouc2026/services/api-server/internal/modules/iam/repo"
	"github.com/liangluo/weouc2026/services/api-server/internal/modules/system"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/auth"
	appconfig "github.com/liangluo/weouc2026/services/api-server/internal/platform/config"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/httpx"
	"github.com/liangluo/weouc2026/services/api-server/internal/providers/academic_provider"
	"github.com/liangluo/weouc2026/services/api-server/internal/providers/wechat_provider"
)

type App struct {
	config  appconfig.AppConfig
	logger  *slog.Logger
	router  *gin.Engine
	closers []io.Closer
}

func New(cfg appconfig.AppConfig, logger *slog.Logger) *App {
	router, closers := buildRouter(cfg, logger)
	return &App{
		config:  cfg,
		logger:  logger,
		router:  router,
		closers: closers,
	}
}

func NewRouter(cfg appconfig.AppConfig, logger *slog.Logger) *gin.Engine {
	router, _ := buildRouter(cfg, logger)
	return router
}

func buildRouter(cfg appconfig.AppConfig, logger *slog.Logger) (*gin.Engine, []io.Closer) {
	gin.SetMode(resolveGinMode(cfg.Service.Environment))

	userRepository := iamrepo.NewInMemoryUserRepository()
	sessionRepository := iamrepo.NewInMemorySessionRepository()
	captchaRepository := iamrepo.NewInMemoryCaptchaRepository()
	wechatProvider := wechat_provider.NewMockProvider()
	academicProvider := academic_provider.NewMockProvider()
	campusLifeRepository := clrepo.NewInMemoryRepository()

	iamModule := iam.NewModule(cfg, iam.Dependencies{
		UserRepository:    userRepository,
		SessionRepository: sessionRepository,
		CaptchaRepository: captchaRepository,
		WeChatProvider:    wechatProvider,
		AcademicProvider:  academicProvider,
	})

	engine := gin.New()
	engine.Use(
		httpx.RequestIDMiddleware(),
		httpx.AccessLogMiddleware(logger),
		httpx.RecoveryMiddleware(logger),
		auth.ContextMiddleware(cfg, iamModule.AuthResolver()),
	)

	iamModule.RegisterRoutes(engine)
	campus_life.NewModule(campus_life.Dependencies{
		Repository: campusLifeRepository,
	}).RegisterRoutes(engine)
	systemModule := system.NewModule(cfg)
	systemModule.RegisterRoutes(engine)

	return engine, systemModule.Closers()
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
		return errors.Join(err, closeClosers(a.closers))
	case <-ctx.Done():
		a.logger.Info("api server shutting down")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), a.config.Server.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return err
		}

		return errors.Join(<-errCh, closeClosers(a.closers))
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

func closeClosers(closers []io.Closer) error {
	var errs []error
	for _, closer := range closers {
		if closer == nil {
			continue
		}
		if err := closer.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
