package config

import (
	"testing"
	"time"
)

func TestLoadUsesDefaults(t *testing.T) {
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
	if !cfg.Dependencies.MySQL.Enabled {
		t.Fatal("expected mysql dependency to be enabled by default")
	}
	if !cfg.Dependencies.Mongo.Enabled {
		t.Fatal("expected mongo dependency to be enabled by default")
	}
	if !cfg.Dependencies.Redis.Enabled {
		t.Fatal("expected redis dependency to be enabled by default")
	}
	if cfg.Persistence.IAMBackendOrDefault() != "mysql_redis" {
		t.Fatalf("expected default iam backend mysql_redis, got %q", cfg.Persistence.IAMBackendOrDefault())
	}
	if cfg.Persistence.CampusLifeBackendOrDefault() != "mongo" {
		t.Fatalf("expected default campus_life backend mongo, got %q", cfg.Persistence.CampusLifeBackendOrDefault())
	}
	if cfg.Persistence.PortalBackendOrDefault() != "mongo" {
		t.Fatalf("expected default portal backend mongo, got %q", cfg.Persistence.PortalBackendOrDefault())
	}
	if cfg.Persistence.NotificationBackendOrDefault() != "mongo" {
		t.Fatalf("expected default notification backend mongo, got %q", cfg.Persistence.NotificationBackendOrDefault())
	}
	if cfg.Persistence.AnalyticsBackendOrDefault() != "mongo" {
		t.Fatalf("expected default analytics backend mongo, got %q", cfg.Persistence.AnalyticsBackendOrDefault())
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
	t.Setenv("API_SERVER_MYSQL_ENABLED", "true")
	t.Setenv("API_SERVER_MYSQL_HOST", "mysql")
	t.Setenv("API_SERVER_MYSQL_DATABASE", "weouc_mysql")
	t.Setenv("API_SERVER_MONGO_ENABLED", "true")
	t.Setenv("API_SERVER_MONGO_URI", "mongodb://mongo:27017/?directConnection=true")
	t.Setenv("API_SERVER_MONGO_DATABASE", "weouc_mongo")
	t.Setenv("API_SERVER_REDIS_ENABLED", "true")
	t.Setenv("API_SERVER_REDIS_HOST", "redis")
	t.Setenv("API_SERVER_REDIS_DB", "2")
	t.Setenv("API_SERVER_COS_ENABLED", "true")
	t.Setenv("API_SERVER_COS_SECRET_ID", "secret-id")
	t.Setenv("API_SERVER_COS_SECRET_KEY", "secret-key")
	t.Setenv("API_SERVER_COS_BUCKET", "weouc-1250000000")
	t.Setenv("API_SERVER_COS_REGION", "ap-guangzhou")
	t.Setenv("API_SERVER_IAM_BACKEND", "mysql_redis")
	t.Setenv("API_SERVER_CAMPUS_LIFE_BACKEND", "mongo")
	t.Setenv("API_SERVER_PORTAL_BACKEND", "mongo")
	t.Setenv("API_SERVER_NOTIFICATION_BACKEND", "mongo")
	t.Setenv("API_SERVER_ANALYTICS_BACKEND", "mongo")
	t.Setenv("API_SERVER_AUTO_MIGRATE", "true")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.Service.Name != "weouc-api" {
		t.Fatalf("expected overridden service name, got %q", cfg.Service.Name)
	}
	if cfg.Server.Address() != "127.0.0.1:9090" {
		t.Fatalf("expected overridden server address, got %q", cfg.Server.Address())
	}
	if !cfg.Dependencies.MySQL.Enabled || cfg.Dependencies.MySQL.Address() != "mysql:3306" {
		t.Fatalf("expected overridden mysql config, got %+v", cfg.Dependencies.MySQL)
	}
	if !cfg.Dependencies.Mongo.Enabled || cfg.Dependencies.Mongo.Database != "weouc_mongo" {
		t.Fatalf("expected overridden mongo config, got %+v", cfg.Dependencies.Mongo)
	}
	if !cfg.Dependencies.Redis.Enabled || cfg.Dependencies.Redis.Address() != "redis:6379" || cfg.Dependencies.Redis.Database != 2 {
		t.Fatalf("expected overridden redis config, got %+v", cfg.Dependencies.Redis)
	}
	if !cfg.Persistence.AutoMigrate {
		t.Fatal("expected auto migrate to be enabled")
	}
}

func TestValidateRejectsInvalidPort(t *testing.T) {
	cfg := validAppConfig()
	cfg.Server.Port = 0

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected Validate() to reject invalid port")
	}
}

func TestValidateRejectsInvalidEnabledMySQLConfig(t *testing.T) {
	cfg := validAppConfig()
	cfg.Dependencies.MySQL = MySQLConfig{
		Enabled:            true,
		Port:               3306,
		Database:           "weouc",
		User:               "weouc",
		HealthCheckTimeout: time.Second,
	}

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected Validate() to reject invalid enabled mysql config")
	}
}

func TestValidateRejectsInvalidEnabledMongoConfig(t *testing.T) {
	cfg := validAppConfig()
	cfg.Dependencies.Mongo = MongoConfig{
		Enabled:            true,
		URI:                "",
		Database:           "weouc",
		HealthCheckTimeout: time.Second,
	}

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected Validate() to reject invalid enabled mongo config")
	}
}

func TestValidateRejectsInvalidIAMBackend(t *testing.T) {
	cfg := validAppConfig()
	cfg.Persistence.IAMBackend = "unknown"

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected Validate() to reject invalid iam backend")
	}
}

func TestValidateRejectsInvalidCampusLifeBackend(t *testing.T) {
	cfg := validAppConfig()
	cfg.Persistence.CampusLifeBackend = "unknown"

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected Validate() to reject invalid campus_life backend")
	}
}

func TestValidateRejectsPersistentIAMWithoutDependencies(t *testing.T) {
	cfg := validAppConfig()
	cfg.Persistence.IAMBackend = "mysql_redis"

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected Validate() to reject mysql_redis backend without dependencies")
	}
}

func TestValidateRejectsMongoBackendsWithoutMongoDependency(t *testing.T) {
	cfg := validAppConfig()
	cfg.Persistence.CampusLifeBackend = "mongo"
	cfg.Persistence.PortalBackend = "mongo"
	cfg.Persistence.NotificationBackend = "mongo"
	cfg.Persistence.AnalyticsBackend = "mongo"

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected Validate() to reject mongo backends without mongo dependency")
	}
}

func TestValidateRejectsEnabledCOSWithoutBucketAppID(t *testing.T) {
	cfg := validAppConfig()
	cfg.Dependencies.COS = COSConfig{
		Enabled:            true,
		SecretID:           "secret-id",
		SecretKey:          "secret-key",
		Bucket:             "weouc",
		Region:             "ap-guangzhou",
		STSDuration:        time.Hour,
		PresignedGETTTL:    6 * time.Hour,
		HealthCheckTimeout: time.Second,
	}

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected Validate() to reject invalid cos bucket format")
	}
}

func validAppConfig() AppConfig {
	return AppConfig{
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
	}
}
