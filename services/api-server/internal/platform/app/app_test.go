package app

import (
	"errors"
	"io"
	"testing"
	"time"

	appconfig "github.com/liangluo/weouc2026/services/api-server/internal/platform/config"
	applogger "github.com/liangluo/weouc2026/services/api-server/internal/platform/logger"
)

func TestNewRouterFailsWhenRequiredBackendsAreUnavailable(t *testing.T) {
	cfg := appconfig.AppConfig{
		Service: appconfig.ServiceConfig{
			Name:        "api-server",
			Environment: "test",
			Version:     "test-version",
		},
		Server: appconfig.ServerConfig{
			Host:            "127.0.0.1",
			Port:            8080,
			ReadTimeout:     time.Second,
			WriteTimeout:    time.Second,
			ShutdownTimeout: time.Second,
		},
		Auth: appconfig.AuthConfig{
			UserIDHeader:        "X-User-ID",
			RolesHeader:         "X-User-Roles",
			PermissionsHeader:   "X-User-Permissions",
			AcademicBoundHeader: "X-Academic-Bound",
			AccessTokenTTL:      time.Hour,
		},
		Dependencies: appconfig.DependenciesConfig{
			MySQL: appconfig.MySQLConfig{
				Enabled:            true,
				Host:               "127.0.0.1",
				Port:               1,
				Database:           "weouc",
				User:               "weouc",
				Password:           "weouc",
				Params:             "charset=utf8mb4&parseTime=True&loc=Local",
				HealthCheckTimeout: 100 * time.Millisecond,
			},
			Mongo: appconfig.MongoConfig{
				Enabled:            true,
				URI:                "mongodb://127.0.0.1:1/?directConnection=true",
				Database:           "weouc",
				HealthCheckTimeout: 100 * time.Millisecond,
			},
			Redis: appconfig.RedisConfig{
				Enabled:            true,
				Host:               "127.0.0.1",
				Port:               1,
				HealthCheckTimeout: 100 * time.Millisecond,
			},
		},
		Persistence: appconfig.PersistenceConfig{
			IAMBackend:          "mysql_redis",
			CampusLifeBackend:   "mongo",
			PortalBackend:       "mongo",
			NotificationBackend: "mongo",
			AnalyticsBackend:    "mongo",
		},
	}

	if _, err := NewRouter(cfg, applogger.New(cfg)); err == nil {
		t.Fatal("expected NewRouter to fail when required backends are unavailable")
	}
}

func TestResolveGinMode(t *testing.T) {
	if mode := resolveGinMode("production"); mode != "release" {
		t.Fatalf("expected release mode, got %q", mode)
	}
	if mode := resolveGinMode("test"); mode != "test" {
		t.Fatalf("expected test mode, got %q", mode)
	}
	if mode := resolveGinMode("local"); mode != "debug" {
		t.Fatalf("expected debug mode, got %q", mode)
	}
}

func TestCloseClosersAggregatesErrors(t *testing.T) {
	errOne := errors.New("close one failed")
	errTwo := errors.New("close two failed")

	err := closeClosers([]io.Closer{
		closerStub{err: errOne},
		nil,
		closerStub{err: errTwo},
	})
	if err == nil {
		t.Fatal("expected closeClosers to return aggregated error")
	}
	if !errors.Is(err, errOne) || !errors.Is(err, errTwo) {
		t.Fatalf("expected aggregated error to contain both failures, got %v", err)
	}
}

type closerStub struct {
	err error
}

func (c closerStub) Close() error {
	return c.err
}
