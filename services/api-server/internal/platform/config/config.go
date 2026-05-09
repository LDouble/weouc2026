package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type AppConfig struct {
	Service ServiceConfig
	Server  ServerConfig
	Auth    AuthConfig
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

	return nil
}

func (c ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
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
