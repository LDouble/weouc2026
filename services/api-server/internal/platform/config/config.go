package config

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type AppConfig struct {
	Service      ServiceConfig
	Server       ServerConfig
	Auth         AuthConfig
	Dependencies DependenciesConfig
	Persistence  PersistenceConfig
}

type ServiceConfig struct {
	Name        string
	Environment string
	Version     string
}

type ServerConfig struct {
	Host            string
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

type AuthConfig struct {
	UserIDHeader        string
	RolesHeader         string
	PermissionsHeader   string
	AcademicBoundHeader string
	AccessTokenTTL      time.Duration
}

type DependenciesConfig struct {
	Postgres PostgresConfig
	Redis    RedisConfig
	COS      COSConfig
}

type PersistenceConfig struct {
	IAMBackend  string
	AutoMigrate bool
}

type PostgresConfig struct {
	Enabled            bool
	Host               string
	Port               int
	Database           string
	User               string
	Password           string
	SSLMode            string
	HealthCheckTimeout time.Duration
}

type RedisConfig struct {
	Enabled            bool
	Host               string
	Port               int
	Username           string
	Password           string
	Database           int
	HealthCheckTimeout time.Duration
}

type COSConfig struct {
	Enabled            bool
	SecretID           string
	SecretKey          string
	Bucket             string
	Region             string
	PathPrefix         string
	STSDuration        time.Duration
	PresignedGETTTL    time.Duration
	HealthCheckTimeout time.Duration
}

func Load() (AppConfig, error) {
	cfg := AppConfig{
		Service: ServiceConfig{
			Name:        getenv("API_SERVER_NAME", "api-server"),
			Environment: getenv("API_SERVER_ENV", "local"),
			Version:     getenv("API_SERVER_VERSION", "dev"),
		},
		Server: ServerConfig{
			Host:            getenv("API_SERVER_HOST", "0.0.0.0"),
			Port:            getenvInt("API_SERVER_PORT", 8080),
			ReadTimeout:     getenvDuration("API_SERVER_READ_TIMEOUT", 5*time.Second),
			WriteTimeout:    getenvDuration("API_SERVER_WRITE_TIMEOUT", 10*time.Second),
			ShutdownTimeout: getenvDuration("API_SERVER_SHUTDOWN_TIMEOUT", 10*time.Second),
		},
		Auth: AuthConfig{
			UserIDHeader:        getenv("API_SERVER_AUTH_USER_ID_HEADER", "X-User-ID"),
			RolesHeader:         getenv("API_SERVER_AUTH_ROLES_HEADER", "X-User-Roles"),
			PermissionsHeader:   getenv("API_SERVER_AUTH_PERMISSIONS_HEADER", "X-User-Permissions"),
			AcademicBoundHeader: getenv("API_SERVER_AUTH_ACADEMIC_BOUND_HEADER", "X-Academic-Bound"),
			AccessTokenTTL:      getenvDuration("API_SERVER_AUTH_ACCESS_TOKEN_TTL", 24*time.Hour),
		},
		Dependencies: DependenciesConfig{
			Postgres: PostgresConfig{
				Enabled:            getenvBool("API_SERVER_POSTGRES_ENABLED", false),
				Host:               getenv("API_SERVER_POSTGRES_HOST", "127.0.0.1"),
				Port:               getenvInt("API_SERVER_POSTGRES_PORT", 5432),
				Database:           getenv("API_SERVER_POSTGRES_DATABASE", "weouc"),
				User:               getenv("API_SERVER_POSTGRES_USER", "weouc"),
				Password:           getenv("API_SERVER_POSTGRES_PASSWORD", "weouc"),
				SSLMode:            getenv("API_SERVER_POSTGRES_SSL_MODE", "disable"),
				HealthCheckTimeout: getenvDuration("API_SERVER_POSTGRES_HEALTHCHECK_TIMEOUT", 2*time.Second),
			},
			Redis: RedisConfig{
				Enabled:            getenvBool("API_SERVER_REDIS_ENABLED", false),
				Host:               getenv("API_SERVER_REDIS_HOST", "127.0.0.1"),
				Port:               getenvInt("API_SERVER_REDIS_PORT", 6379),
				Username:           getenv("API_SERVER_REDIS_USERNAME", ""),
				Password:           getenv("API_SERVER_REDIS_PASSWORD", ""),
				Database:           getenvInt("API_SERVER_REDIS_DB", 0),
				HealthCheckTimeout: getenvDuration("API_SERVER_REDIS_HEALTHCHECK_TIMEOUT", 2*time.Second),
			},
			COS: COSConfig{
				Enabled:            getenvBool("API_SERVER_COS_ENABLED", false),
				SecretID:           getenv("API_SERVER_COS_SECRET_ID", ""),
				SecretKey:          getenv("API_SERVER_COS_SECRET_KEY", ""),
				Bucket:             getenv("API_SERVER_COS_BUCKET", ""),
				Region:             getenv("API_SERVER_COS_REGION", ""),
				PathPrefix:         getenv("API_SERVER_COS_PATH_PREFIX", "miniapp"),
				STSDuration:        getenvDuration("API_SERVER_COS_STS_DURATION", time.Hour),
				PresignedGETTTL:    getenvDuration("API_SERVER_COS_PRESIGNED_GET_TTL", 6*time.Hour),
				HealthCheckTimeout: getenvDuration("API_SERVER_COS_HEALTHCHECK_TIMEOUT", 2*time.Second),
			},
		},
		Persistence: PersistenceConfig{
			IAMBackend:  getenv("API_SERVER_IAM_BACKEND", "memory"),
			AutoMigrate: getenvBool("API_SERVER_AUTO_MIGRATE", false),
		},
	}

	return cfg, cfg.Validate()
}

func (c AppConfig) Validate() error {
	if c.Service.Name == "" {
		return errors.New("service name is required")
	}
	if c.Server.Host == "" {
		return errors.New("server host is required")
	}
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("server port %d is invalid", c.Server.Port)
	}
	if c.Server.ReadTimeout <= 0 || c.Server.WriteTimeout <= 0 || c.Server.ShutdownTimeout <= 0 {
		return errors.New("server timeouts must be positive")
	}
	if c.Auth.UserIDHeader == "" || c.Auth.RolesHeader == "" || c.Auth.PermissionsHeader == "" || c.Auth.AcademicBoundHeader == "" {
		return errors.New("auth headers must not be empty")
	}
	if c.Auth.AccessTokenTTL <= 0 {
		return errors.New("access token ttl must be positive")
	}
	if err := c.Dependencies.Postgres.Validate(); err != nil {
		return err
	}
	if err := c.Dependencies.Redis.Validate(); err != nil {
		return err
	}
	if err := c.Dependencies.COS.Validate(); err != nil {
		return err
	}
	if err := c.Persistence.Validate(c.Dependencies); err != nil {
		return err
	}

	return nil
}

func (c ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c PostgresConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c PostgresConfig) DSN() string {
	timeoutSeconds := int(math.Ceil(c.HealthCheckTimeout.Seconds()))
	if timeoutSeconds <= 0 {
		timeoutSeconds = 1
	}

	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s connect_timeout=%d",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.Database,
		c.SSLMode,
		timeoutSeconds,
	)
}

func (c PostgresConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Host == "" {
		return errors.New("postgres host is required when enabled")
	}
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("postgres port %d is invalid", c.Port)
	}
	if c.Database == "" {
		return errors.New("postgres database is required when enabled")
	}
	if c.User == "" {
		return errors.New("postgres user is required when enabled")
	}
	if c.SSLMode == "" {
		return errors.New("postgres ssl mode is required when enabled")
	}
	if c.HealthCheckTimeout <= 0 {
		return errors.New("postgres health check timeout must be positive")
	}

	return nil
}

func (c RedisConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c RedisConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Host == "" {
		return errors.New("redis host is required when enabled")
	}
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("redis port %d is invalid", c.Port)
	}
	if c.Database < 0 {
		return fmt.Errorf("redis db %d is invalid", c.Database)
	}
	if c.HealthCheckTimeout <= 0 {
		return errors.New("redis health check timeout must be positive")
	}

	return nil
}

func (c COSConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.SecretID == "" {
		return errors.New("cos secret id is required when enabled")
	}
	if c.SecretKey == "" {
		return errors.New("cos secret key is required when enabled")
	}
	if c.Bucket == "" {
		return errors.New("cos bucket is required when enabled")
	}
	if c.Region == "" {
		return errors.New("cos region is required when enabled")
	}
	if c.BucketAppID() == "" {
		return errors.New("cos bucket must use bucket-appid format when enabled")
	}
	if c.STSDuration <= 0 {
		return errors.New("cos sts duration must be positive")
	}
	if c.PresignedGETTTL <= 0 {
		return errors.New("cos presigned get ttl must be positive")
	}
	if c.HealthCheckTimeout <= 0 {
		return errors.New("cos health check timeout must be positive")
	}

	return nil
}

func (c COSConfig) BucketAppID() string {
	bucket := c.Bucket
	if bucket == "" {
		return ""
	}

	index := strings.LastIndex(bucket, "-")
	if index <= 0 || index >= len(bucket)-1 {
		return ""
	}

	appID := bucket[index+1:]
	if _, err := strconv.ParseInt(appID, 10, 64); err != nil {
		return ""
	}

	return appID
}

func (c PersistenceConfig) IAMBackendOrDefault() string {
	if c.IAMBackend == "" {
		return "memory"
	}

	return c.IAMBackend
}

func (c PersistenceConfig) Validate(dependencies DependenciesConfig) error {
	switch c.IAMBackendOrDefault() {
	case "memory":
		return nil
	case "postgres_redis":
		if !dependencies.Postgres.Enabled {
			return errors.New("postgres must be enabled when iam backend is postgres_redis")
		}
		if !dependencies.Redis.Enabled {
			return errors.New("redis must be enabled when iam backend is postgres_redis")
		}
		return nil
	default:
		return fmt.Errorf("iam backend %q is invalid", c.IAMBackend)
	}
}

func getenv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func getenvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		parsed, err := strconv.Atoi(value)
		if err == nil {
			return parsed
		}
	}

	return defaultValue
}

func getenvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		parsed, err := time.ParseDuration(value)
		if err == nil {
			return parsed
		}
	}

	return defaultValue
}

func getenvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		parsed, err := strconv.ParseBool(value)
		if err == nil {
			return parsed
		}
	}

	return defaultValue
}
