package config

import (
	"errors"
	"fmt"
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
	MySQL MySQLConfig
	Mongo MongoConfig
	Redis RedisConfig
	COS   COSConfig
}

type PersistenceConfig struct {
	IAMBackend          string
	CampusLifeBackend   string
	PortalBackend       string
	NotificationBackend string
	AnalyticsBackend    string
	AutoMigrate         bool
}

type MySQLConfig struct {
	Enabled            bool
	Host               string
	Port               int
	Database           string
	User               string
	Password           string
	Params             string
	HealthCheckTimeout time.Duration
}

type MongoConfig struct {
	Enabled            bool
	URI                string
	Database           string
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
			MySQL: MySQLConfig{
				Enabled:            getenvBool("API_SERVER_MYSQL_ENABLED", true),
				Host:               getenv("API_SERVER_MYSQL_HOST", "127.0.0.1"),
				Port:               getenvInt("API_SERVER_MYSQL_PORT", 3306),
				Database:           getenv("API_SERVER_MYSQL_DATABASE", "weouc"),
				User:               getenv("API_SERVER_MYSQL_USER", "weouc"),
				Password:           getenv("API_SERVER_MYSQL_PASSWORD", "weouc"),
				Params:             getenv("API_SERVER_MYSQL_PARAMS", "charset=utf8mb4&parseTime=True&loc=Local"),
				HealthCheckTimeout: getenvDuration("API_SERVER_MYSQL_HEALTHCHECK_TIMEOUT", 2*time.Second),
			},
			Mongo: MongoConfig{
				Enabled:            getenvBool("API_SERVER_MONGO_ENABLED", true),
				URI:                getenv("API_SERVER_MONGO_URI", "mongodb://127.0.0.1:27017/?directConnection=true"),
				Database:           getenv("API_SERVER_MONGO_DATABASE", "weouc"),
				HealthCheckTimeout: getenvDuration("API_SERVER_MONGO_HEALTHCHECK_TIMEOUT", 2*time.Second),
			},
			Redis: RedisConfig{
				Enabled:            getenvBool("API_SERVER_REDIS_ENABLED", true),
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
			IAMBackend:          getenv("API_SERVER_IAM_BACKEND", "mysql_redis"),
			CampusLifeBackend:   getenv("API_SERVER_CAMPUS_LIFE_BACKEND", "mongo"),
			PortalBackend:       getenv("API_SERVER_PORTAL_BACKEND", "mongo"),
			NotificationBackend: getenv("API_SERVER_NOTIFICATION_BACKEND", "mongo"),
			AnalyticsBackend:    getenv("API_SERVER_ANALYTICS_BACKEND", "mongo"),
			AutoMigrate:         getenvBool("API_SERVER_AUTO_MIGRATE", false),
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
	if err := c.Dependencies.MySQL.Validate(); err != nil {
		return err
	}
	if err := c.Dependencies.Mongo.Validate(); err != nil {
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

func (c MySQLConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c MySQLConfig) DSN() string {
	params := strings.TrimSpace(c.Params)
	if params == "" {
		params = "charset=utf8mb4&parseTime=True&loc=Local"
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", c.User, c.Password, c.Host, c.Port, c.Database, params)
}

func (c MySQLConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Host == "" {
		return errors.New("mysql host is required when enabled")
	}
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("mysql port %d is invalid", c.Port)
	}
	if c.Database == "" {
		return errors.New("mysql database is required when enabled")
	}
	if c.User == "" {
		return errors.New("mysql user is required when enabled")
	}
	if c.HealthCheckTimeout <= 0 {
		return errors.New("mysql health check timeout must be positive")
	}

	return nil
}

func (c MongoConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if strings.TrimSpace(c.URI) == "" {
		return errors.New("mongo uri is required when enabled")
	}
	if c.Database == "" {
		return errors.New("mongo database is required when enabled")
	}
	if c.HealthCheckTimeout <= 0 {
		return errors.New("mongo health check timeout must be positive")
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
		return "mysql_redis"
	}

	return c.IAMBackend
}

func (c PersistenceConfig) CampusLifeBackendOrDefault() string {
	if c.CampusLifeBackend == "" {
		return "mongo"
	}

	return c.CampusLifeBackend
}

func (c PersistenceConfig) PortalBackendOrDefault() string {
	if c.PortalBackend == "" {
		return "mongo"
	}

	return c.PortalBackend
}

func (c PersistenceConfig) NotificationBackendOrDefault() string {
	if c.NotificationBackend == "" {
		return "mongo"
	}

	return c.NotificationBackend
}

func (c PersistenceConfig) AnalyticsBackendOrDefault() string {
	if c.AnalyticsBackend == "" {
		return "mongo"
	}

	return c.AnalyticsBackend
}

func (c PersistenceConfig) Validate(dependencies DependenciesConfig) error {
	switch c.IAMBackendOrDefault() {
	case "mysql_redis":
		if !dependencies.MySQL.Enabled {
			return errors.New("mysql must be enabled when iam backend is mysql_redis")
		}
		if !dependencies.Redis.Enabled {
			return errors.New("redis must be enabled when iam backend is mysql_redis")
		}
	default:
		return fmt.Errorf("iam backend %q is invalid", c.IAMBackend)
	}

	switch c.CampusLifeBackendOrDefault() {
	case "mongo":
		if !dependencies.Mongo.Enabled {
			return errors.New("mongo must be enabled when campus_life backend is mongo")
		}
	default:
		return fmt.Errorf("campus_life backend %q is invalid", c.CampusLifeBackend)
	}

	switch c.PortalBackendOrDefault() {
	case "mongo":
		if !dependencies.Mongo.Enabled {
			return errors.New("mongo must be enabled when portal backend is mongo")
		}
	default:
		return fmt.Errorf("portal backend %q is invalid", c.PortalBackend)
	}

	switch c.NotificationBackendOrDefault() {
	case "mongo":
		if !dependencies.Mongo.Enabled {
			return errors.New("mongo must be enabled when notification backend is mongo")
		}
	default:
		return fmt.Errorf("notification backend %q is invalid", c.NotificationBackend)
	}

	switch c.AnalyticsBackendOrDefault() {
	case "mongo":
		if !dependencies.Mongo.Enabled {
			return errors.New("mongo must be enabled when analytics backend is mongo")
		}
	default:
		return fmt.Errorf("analytics backend %q is invalid", c.AnalyticsBackend)
	}

	return nil
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
