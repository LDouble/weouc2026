package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/liangluo/weouc2026/services/api-server/internal/modules/analytics"
	"github.com/liangluo/weouc2026/services/api-server/internal/modules/campus_life"
	clrepo "github.com/liangluo/weouc2026/services/api-server/internal/modules/campus_life/repo"
	"github.com/liangluo/weouc2026/services/api-server/internal/modules/file_center"
	"github.com/liangluo/weouc2026/services/api-server/internal/modules/iam"
	iamrepo "github.com/liangluo/weouc2026/services/api-server/internal/modules/iam/repo"
	"github.com/liangluo/weouc2026/services/api-server/internal/modules/notification"
	notificationrepo "github.com/liangluo/weouc2026/services/api-server/internal/modules/notification/repo"
	"github.com/liangluo/weouc2026/services/api-server/internal/modules/portal"
	portalrepo "github.com/liangluo/weouc2026/services/api-server/internal/modules/portal/repo"
	"github.com/liangluo/weouc2026/services/api-server/internal/modules/system"
	systemrepo "github.com/liangluo/weouc2026/services/api-server/internal/modules/system/repo"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/audit"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/auth"
	appconfig "github.com/liangluo/weouc2026/services/api-server/internal/platform/config"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/httpx"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/migrate"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/persistence"
	"github.com/liangluo/weouc2026/services/api-server/internal/providers/academic_provider"
	"github.com/liangluo/weouc2026/services/api-server/internal/providers/storage_provider"
	"github.com/liangluo/weouc2026/services/api-server/internal/providers/wechat_provider"
)

type App struct {
	config  appconfig.AppConfig
	logger  *slog.Logger
	router  *gin.Engine
	closers []io.Closer
}

func New(cfg appconfig.AppConfig, logger *slog.Logger) (*App, error) {
	router, closers, err := buildRouter(cfg, logger)
	if err != nil {
		return nil, err
	}

	return &App{
		config:  cfg,
		logger:  logger,
		router:  router,
		closers: closers,
	}, nil
}

func NewRouter(cfg appconfig.AppConfig, logger *slog.Logger) (*gin.Engine, error) {
	router, _, err := buildRouter(cfg, logger)
	return router, err
}

func buildRouter(cfg appconfig.AppConfig, logger *slog.Logger) (*gin.Engine, []io.Closer, error) {
	gin.SetMode(resolveGinMode(cfg.Service.Environment))

	clients, err := persistence.Open(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("open runtime clients failed: %w", err)
	}
	closers := clients.Closers()

	if err := clients.EnsureRuntimeBackendsReady(context.Background(), cfg); err != nil {
		closeClosers(closers)
		return nil, nil, fmt.Errorf("ensure runtime backends ready failed: %w", err)
	}
	if cfg.Persistence.AutoMigrate {
		if clients.Postgres == nil {
			closeClosers(closers)
			return nil, nil, fmt.Errorf("postgres client is required when auto migrate is enabled")
		}
		if err := migrate.Run(context.Background(), clients.Postgres); err != nil {
			closeClosers(closers)
			return nil, nil, fmt.Errorf("run migrations failed: %w", err)
		}
	}

	userRepository, sessionRepository, captchaRepository, err := newIAMRepositories(cfg, clients)
	if err != nil {
		closeClosers(closers)
		return nil, nil, err
	}

	wechatProvider := wechat_provider.NewMockProvider()
	academicProvider := academic_provider.NewMockProvider()
	storageProvider, err := newStorageProvider(cfg)
	if err != nil {
		closeClosers(closers)
		return nil, nil, err
	}
	campusLifeRepository, err := newCampusLifeRepository(cfg, clients)
	if err != nil {
		closeClosers(closers)
		return nil, nil, err
	}
	portalRepository := portalrepo.NewInMemoryRepository()
	notificationRepository := notificationrepo.NewInMemoryRepository()
	auditStore := audit.NewInMemoryStore()

	iamModule := iam.NewModule(cfg, iam.Dependencies{
		UserRepository:    userRepository,
		SessionRepository: sessionRepository,
		CaptchaRepository: captchaRepository,
		WeChatProvider:    wechatProvider,
		AcademicProvider:  academicProvider,
		AuditRecorder:     auditStore,
	})

	engine := gin.New()
	engine.Use(
		httpx.RequestIDMiddleware(),
		httpx.AccessLogMiddleware(logger),
		httpx.RecoveryMiddleware(logger),
		auth.ContextMiddleware(cfg, iamModule.AuthResolver()),
	)

	iamModule.RegisterRoutes(engine)
	portal.NewModule(portal.Dependencies{
		Repository:    portalRepository,
		AuditRecorder: auditStore,
	}).RegisterRoutes(engine)
	notification.NewModule(notification.Dependencies{
		Repository:    notificationRepository,
		AuditRecorder: auditStore,
	}).RegisterRoutes(engine)
	analytics.NewModule(analytics.Dependencies{
		AuditStore: auditStore,
	}).RegisterRoutes(engine)
	campus_life.NewModule(campus_life.Dependencies{
		Repository:      campusLifeRepository,
		StorageProvider: storageProvider,
		AuditRecorder:   auditStore,
	}).RegisterRoutes(engine)
	file_center.NewModule(file_center.Dependencies{
		StorageProvider: storageProvider,
	}).RegisterRoutes(engine)
	systemModule := system.NewModule(cfg, system.Dependencies{
		StatusRepository: systemrepo.NewRuntimeStatusRepository(
			systemrepo.NewPostgresProbe(cfg.Dependencies.Postgres, clients.Postgres),
			systemrepo.NewRedisProbe(cfg.Dependencies.Redis, clients.Redis),
			systemrepo.NewObjectStorageProbe(cfg.Dependencies.COS, storageProvider),
		),
	})
	systemModule.RegisterRoutes(engine)

	return engine, closers, nil
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

func newIAMRepositories(
	cfg appconfig.AppConfig,
	clients *persistence.Clients,
) (iamrepo.UserRepository, iamrepo.SessionRepository, iamrepo.CaptchaRepository, error) {
	switch cfg.Persistence.IAMBackendOrDefault() {
	case "memory":
		return iamrepo.NewInMemoryUserRepository(), iamrepo.NewInMemorySessionRepository(), iamrepo.NewInMemoryCaptchaRepository(), nil
	case "postgres_redis":
		if clients == nil || clients.Postgres == nil || clients.Redis == nil {
			return nil, nil, nil, fmt.Errorf("postgres_redis backend requires postgres and redis clients")
		}
		return iamrepo.NewPostgresUserRepository(clients.Postgres), iamrepo.NewRedisSessionRepository(clients.Redis), iamrepo.NewRedisCaptchaRepository(clients.Redis), nil
	default:
		return nil, nil, nil, fmt.Errorf("unsupported iam backend %q", cfg.Persistence.IAMBackendOrDefault())
	}
}

func newCampusLifeRepository(cfg appconfig.AppConfig, clients *persistence.Clients) (clrepo.Repository, error) {
	switch cfg.Persistence.CampusLifeBackendOrDefault() {
	case "memory":
		return clrepo.NewInMemoryRepository(), nil
	case "postgres":
		if clients == nil || clients.Postgres == nil {
			return nil, fmt.Errorf("postgres campus_life backend requires postgres client")
		}
		return clrepo.NewPostgresRepository(clients.Postgres), nil
	default:
		return nil, fmt.Errorf("unsupported campus_life backend %q", cfg.Persistence.CampusLifeBackendOrDefault())
	}
}

func newStorageProvider(cfg appconfig.AppConfig) (storage_provider.Provider, error) {
	if !cfg.Dependencies.COS.Enabled {
		return nil, nil
	}

	provider, err := storage_provider.NewCOSProvider(cfg.Dependencies.COS)
	if err != nil {
		return nil, fmt.Errorf("create cos storage provider failed: %w", err)
	}

	return provider, nil
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
