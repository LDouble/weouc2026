package config

import (
	"testing"
	"time"
)

func TestLoadUsesDefaults(t *testing.T) {
	t.Setenv("API_SERVER_NAME", "")
	t.Setenv("API_SERVER_ENV", "")
	t.Setenv("API_SERVER_VERSION", "")
	t.Setenv("API_SERVER_HOST", "")
	t.Setenv("API_SERVER_PORT", "")
	t.Setenv("API_SERVER_AUTH_ACCESS_TOKEN_TTL", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.Service.Name != "api-server" {
		t.Fatalf("expected default service name, got %q", cfg.Service.Name)
	}
	if cfg.Server.Port != 8080 {
		t.Fatalf("expected default port 8080, got %d", cfg.Server.Port)
	}
	if cfg.Server.ReadTimeout != 5*time.Second {
		t.Fatalf("expected default read timeout, got %s", cfg.Server.ReadTimeout)
	}
	if cfg.Auth.AccessTokenTTL != 24*time.Hour {
		t.Fatalf("expected default auth ttl, got %s", cfg.Auth.AccessTokenTTL)
	}
	if cfg.Dependencies.Postgres.Enabled {
		t.Fatal("expected postgres health probe to be disabled by default")
	}
	if cfg.Dependencies.Redis.Enabled {
		t.Fatal("expected redis health probe to be disabled by default")
	}
}

func TestLoadUsesEnvOverrides(t *testing.T) {
	t.Setenv("API_SERVER_NAME", "weouc-api")
	t.Setenv("API_SERVER_ENV", "test")
	t.Setenv("API_SERVER_VERSION", "1.2.3")
	t.Setenv("API_SERVER_HOST", "127.0.0.1")
	t.Setenv("API_SERVER_PORT", "9090")
	t.Setenv("API_SERVER_READ_TIMEOUT", "7s")
	t.Setenv("API_SERVER_WRITE_TIMEOUT", "12s")
	t.Setenv("API_SERVER_SHUTDOWN_TIMEOUT", "15s")
	t.Setenv("API_SERVER_AUTH_USER_ID_HEADER", "X-Test-User")
	t.Setenv("API_SERVER_AUTH_ACCESS_TOKEN_TTL", "48h")
	t.Setenv("API_SERVER_POSTGRES_ENABLED", "true")
	t.Setenv("API_SERVER_POSTGRES_HOST", "postgres")
	t.Setenv("API_SERVER_POSTGRES_DATABASE", "weouc_dev")
	t.Setenv("API_SERVER_REDIS_ENABLED", "true")
	t.Setenv("API_SERVER_REDIS_HOST", "redis")
	t.Setenv("API_SERVER_REDIS_DB", "2")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.Service.Name != "weouc-api" {
		t.Fatalf("expected overridden service name, got %q", cfg.Service.Name)
	}
	if cfg.Server.Address() != "127.0.0.1:9090" {
		t.Fatalf("expected overridden address, got %q", cfg.Server.Address())
	}
	if cfg.Server.WriteTimeout != 12*time.Second {
		t.Fatalf("expected overridden write timeout, got %s", cfg.Server.WriteTimeout)
	}
	if cfg.Auth.UserIDHeader != "X-Test-User" {
		t.Fatalf("expected overridden auth header, got %q", cfg.Auth.UserIDHeader)
	}
	if cfg.Auth.AccessTokenTTL != 48*time.Hour {
		t.Fatalf("expected overridden auth ttl, got %s", cfg.Auth.AccessTokenTTL)
	}
	if !cfg.Dependencies.Postgres.Enabled || cfg.Dependencies.Postgres.Database != "weouc_dev" {
		t.Fatalf("expected overridden postgres config, got %+v", cfg.Dependencies.Postgres)
	}
	if !cfg.Dependencies.Redis.Enabled || cfg.Dependencies.Redis.Address() != "redis:6379" || cfg.Dependencies.Redis.Database != 2 {
		t.Fatalf("expected overridden redis config, got %+v", cfg.Dependencies.Redis)
	}
}

func TestValidateRejectsInvalidPort(t *testing.T) {
	cfg := AppConfig{
		Service: ServiceConfig{Name: "api-server"},
		Server: ServerConfig{
			Host:            "127.0.0.1",
			Port:            0,
			ReadTimeout:     time.Second,
			WriteTimeout:    time.Second,
			ShutdownTimeout: time.Second,
		},
		Auth: AuthConfig{
			UserIDHeader:        "X-User-ID",
			RolesHeader:         "X-User-Roles",
			PermissionsHeader:   "X-User-Permissions",
			AcademicBoundHeader: "X-Academic-Bound",
			AccessTokenTTL:      time.Hour,
		},
	}

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected Validate() to reject invalid port")
	}
}

func TestValidateRejectsEnabledDependencyWithInvalidConfig(t *testing.T) {
	cfg := AppConfig{
		Service: ServiceConfig{Name: "api-server"},
		Server: ServerConfig{
			Host:            "127.0.0.1",
			Port:            8080,
			ReadTimeout:     time.Second,
			WriteTimeout:    time.Second,
			ShutdownTimeout: time.Second,
		},
		Auth: AuthConfig{
			UserIDHeader:        "X-User-ID",
			RolesHeader:         "X-User-Roles",
			PermissionsHeader:   "X-User-Permissions",
			AcademicBoundHeader: "X-Academic-Bound",
			AccessTokenTTL:      time.Hour,
		},
		Dependencies: DependenciesConfig{
			Postgres: PostgresConfig{
				Enabled:            true,
				Host:               "",
				Port:               5432,
				Database:           "weouc",
				User:               "weouc",
				SSLMode:            "disable",
				HealthCheckTimeout: time.Second,
			},
		},
	}

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected Validate() to reject invalid enabled dependency config")
	}
}
