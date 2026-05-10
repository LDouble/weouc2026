package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	appconfig "github.com/liangluo/weouc2026/services/api-server/internal/platform/config"
	applogger "github.com/liangluo/weouc2026/services/api-server/internal/platform/logger"
)

func TestSystemRoutes(t *testing.T) {
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
		},
	}

	router := NewRouter(cfg, applogger.New(cfg))

	t.Run("healthz is public", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/healthz", nil)

		router.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", recorder.Code)
		}
	})

	t.Run("readyz stays ready when probes are disabled", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/readyz", nil)

		router.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", recorder.Code, recorder.Body.String())
		}
		if !strings.Contains(recorder.Body.String(), "\"status\":\"ready\"") {
			t.Fatalf("expected readiness response to be ready, got %s", recorder.Body.String())
		}
	})

	t.Run("profile requires auth", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/v1/system/profile", nil)

		router.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", recorder.Code)
		}
	})

	t.Run("profile returns principal after auth", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/v1/system/profile", nil)
		request.Header.Set("X-User-ID", "student-001")
		request.Header.Set("X-User-Roles", "student")
		request.Header.Set("X-User-Permissions", "contact:view")
		request.Header.Set("X-Academic-Bound", "true")

		router.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", recorder.Code)
		}

		var payload struct {
			Data struct {
				Auth struct {
					UserID        string `json:"user_id"`
					AcademicBound bool   `json:"academic_bound"`
				} `json:"auth"`
			} `json:"data"`
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
			t.Fatalf("unmarshal response failed: %v", err)
		}
		if payload.Data.Auth.UserID != "student-001" || !payload.Data.Auth.AcademicBound {
			t.Fatalf("unexpected payload: %s", recorder.Body.String())
		}
	})

	t.Run("readyz returns 503 when required dependency is unavailable", func(t *testing.T) {
		unavailableCfg := cfg
		unavailableCfg.Dependencies.Redis.Enabled = true
		unavailableCfg.Dependencies.Redis.Host = "127.0.0.1"
		unavailableCfg.Dependencies.Redis.Port = 1
		unavailableCfg.Dependencies.Redis.HealthCheckTimeout = 100 * time.Millisecond

		unavailableRouter := NewRouter(unavailableCfg, applogger.New(unavailableCfg))
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/readyz", nil)

		unavailableRouter.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusServiceUnavailable {
			t.Fatalf("expected 503, got %d: %s", recorder.Code, recorder.Body.String())
		}
		if !strings.Contains(recorder.Body.String(), "\"status\":\"not_ready\"") {
			t.Fatalf("expected readiness response to be not_ready, got %s", recorder.Body.String())
		}
	})
}
