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
	if cfg.Dependencies.COS.Enabled {
		t.Fatal("expected cos health probe to be disabled by default")
	}
	if cfg.Persistence.IAMBackendOrDefault() != "memory" {
		t.Fatalf("expected default iam backend memory, got %q", cfg.Persistence.IAMBackendOrDefault())
	}
	if cfg.Persistence.CampusLifeBackendOrDefault() != "memory" {
		t.Fatalf("expected default campus_life backend memory, got %q", cfg.Persistence.CampusLifeBackendOrDefault())
	}
	if cfg.Persistence.AutoMigrate {
		t.Fatal("expected auto migrate to be disabled by default")
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
	t.Setenv("API_SERVER_COS_ENABLED", "true")
	t.Setenv("API_SERVER_COS_SECRET_ID", "secret-id")
	t.Setenv("API_SERVER_COS_SECRET_KEY", "secret-key")
	t.Setenv("API_SERVER_COS_BUCKET", "weouc-1250000000")
	t.Setenv("API_SERVER_COS_REGION", "ap-guangzhou")
	t.Setenv("API_SERVER_COS_PATH_PREFIX", "miniapp")
	t.Setenv("API_SERVER_COS_STS_DURATION", "45m")
	t.Setenv("API_SERVER_COS_PRESIGNED_GET_TTL", "8h")
	t.Setenv("API_SERVER_IAM_BACKEND", "postgres_redis")
	t.Setenv("API_SERVER_CAMPUS_LIFE_BACKEND", "postgres")
	t.Setenv("API_SERVER_AUTO_MIGRATE", "true")

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
	if !cfg.Dependencies.COS.Enabled || cfg.Dependencies.COS.BucketAppID() != "1250000000" || cfg.Dependencies.COS.STSDuration != 45*time.Minute {
		t.Fatalf("expected overridden cos config, got %+v", cfg.Dependencies.COS)
	}
	if cfg.Persistence.IAMBackendOrDefault() != "postgres_redis" ||
		cfg.Persistence.CampusLifeBackendOrDefault() != "postgres" ||
		!cfg.Persistence.AutoMigrate {
		t.Fatalf("expected overridden persistence config, got %+v", cfg.Persistence)
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

func TestValidateRejectsInvalidIAMBackend(t *testing.T) {
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
		Persistence: PersistenceConfig{
			IAMBackend: "unknown",
		},
	}

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected Validate() to reject invalid iam backend")
	}
}

func TestValidateRejectsInvalidCampusLifeBackend(t *testing.T) {
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
		Persistence: PersistenceConfig{
			CampusLifeBackend: "unknown",
		},
	}

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected Validate() to reject invalid campus_life backend")
	}
}

func TestValidateRejectsEnabledCOSWithoutBucketAppID(t *testing.T) {
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
			COS: COSConfig{
				Enabled:            true,
				SecretID:           "secret-id",
				SecretKey:          "secret-key",
				Bucket:             "weouc",
				Region:             "ap-guangzhou",
				STSDuration:        time.Hour,
				PresignedGETTTL:    6 * time.Hour,
				HealthCheckTimeout: time.Second,
			},
		},
	}

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected Validate() to reject invalid cos bucket format")
	}
}

func TestValidateRejectsPersistentIAMWithoutDependencies(t *testing.T) {
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
		Persistence: PersistenceConfig{
			IAMBackend: "postgres_redis",
		},
	}

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected Validate() to reject postgres_redis backend without dependencies")
	}
}

func TestValidateRejectsPostgresCampusLifeWithoutDependencies(t *testing.T) {
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
		Persistence: PersistenceConfig{
			CampusLifeBackend: "postgres",
		},
	}

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected Validate() to reject postgres campus_life backend without postgres dependency")
	}
}
