package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	appconfig "github.com/liangluo/weouc2026/services/api-server/internal/platform/config"
	applogger "github.com/liangluo/weouc2026/services/api-server/internal/platform/logger"
)

func TestMiniappCoreAPIs(t *testing.T) {
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
			AccessTokenTTL:      24 * time.Hour,
		},
	}

	router, err := NewRouter(cfg, applogger.New(cfg))
	if err != nil {
		t.Fatalf("NewRouter() returned error: %v", err)
	}
	token := loginAndGetToken(t, router)

	t.Run("student profile returns 404 before binding", func(t *testing.T) {
		recorder := performJSONRequest(t, router, http.MethodGet, "/api/student", token, nil)
		if recorder.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d: %s", recorder.Code, recorder.Body.String())
		}
	})

	t.Run("bind student profile and expose contact after binding", func(t *testing.T) {
		recorder := performJSONRequest(t, router, http.MethodPost, "/api/edu/send-captcha", token, map[string]any{
			"sid": "20260001",
		})
		if recorder.Code != http.StatusOK {
			t.Fatalf("expected captcha status 200, got %d: %s", recorder.Code, recorder.Body.String())
		}

		recorder = performJSONRequest(t, router, http.MethodPost, "/api/student", token, map[string]any{
			"student_id": "20260001",
			"password":   "password-001",
			"captcha":    "123456",
		})
		if recorder.Code != http.StatusOK {
			t.Fatalf("expected bind status 200, got %d: %s", recorder.Code, recorder.Body.String())
		}

		var bindPayload struct {
			Data struct {
				StudentID string `json:"student_id"`
				IsBound   bool   `json:"is_bound"`
			} `json:"data"`
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &bindPayload); err != nil {
			t.Fatalf("unmarshal bind payload failed: %v", err)
		}
		if bindPayload.Data.StudentID != "20260001" || !bindPayload.Data.IsBound {
			t.Fatalf("unexpected bind payload: %s", recorder.Body.String())
		}

		recorder = performJSONRequest(t, router, http.MethodGet, "/api/student", token, nil)
		if recorder.Code != http.StatusOK {
			t.Fatalf("expected student status 200, got %d: %s", recorder.Code, recorder.Body.String())
		}

		recorder = performJSONRequest(t, router, http.MethodGet, "/api/market/detail/market-101", "", nil)
		if recorder.Code != http.StatusOK {
			t.Fatalf("expected market detail status 200, got %d: %s", recorder.Code, recorder.Body.String())
		}
		var publicDetail struct {
			Data struct {
				CanViewContact bool `json:"can_view_contact"`
				Extra          struct {
					Contact string `json:"contact"`
				} `json:"extra"`
			} `json:"data"`
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &publicDetail); err != nil {
			t.Fatalf("unmarshal public market detail failed: %v", err)
		}
		if publicDetail.Data.CanViewContact || publicDetail.Data.Extra.Contact != "" {
			t.Fatalf("expected public market detail to hide contact: %s", recorder.Body.String())
		}

		recorder = performJSONRequest(t, router, http.MethodGet, "/api/market/detail/market-101", token, nil)
		if recorder.Code != http.StatusOK {
			t.Fatalf("expected bound market detail status 200, got %d: %s", recorder.Code, recorder.Body.String())
		}
		var authedDetail struct {
			Data struct {
				CanViewContact bool `json:"can_view_contact"`
				Extra          struct {
					Contact string `json:"contact"`
				} `json:"extra"`
			} `json:"data"`
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &authedDetail); err != nil {
			t.Fatalf("unmarshal authed market detail failed: %v", err)
		}
		if !authedDetail.Data.CanViewContact || authedDetail.Data.Extra.Contact == "" {
			t.Fatalf("expected bound market detail to expose contact: %s", recorder.Body.String())
		}
	})

	t.Run("accept errand updates detail role and status", func(t *testing.T) {
		recorder := performJSONRequest(t, router, http.MethodPost, "/api/errand/accept", token, map[string]any{
			"task_id": "errand-101",
		})
		if recorder.Code != http.StatusOK {
			t.Fatalf("expected accept status 200, got %d: %s", recorder.Code, recorder.Body.String())
		}

		recorder = performJSONRequest(t, router, http.MethodGet, "/api/errand/detail/errand-101", token, nil)
		if recorder.Code != http.StatusOK {
			t.Fatalf("expected errand detail status 200, got %d: %s", recorder.Code, recorder.Body.String())
		}

		var payload struct {
			Data struct {
				UserRole string `json:"user_role"`
				Item     struct {
					Status     string `json:"status"`
					IsAccepted bool   `json:"is_accepted"`
				} `json:"item"`
			} `json:"data"`
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
			t.Fatalf("unmarshal errand detail failed: %v", err)
		}
		if payload.Data.UserRole != "acceptor" || payload.Data.Item.Status != "accepted" || !payload.Data.Item.IsAccepted {
			t.Fatalf("unexpected errand detail payload: %s", recorder.Body.String())
		}
	})
}

func loginAndGetToken(t *testing.T, router http.Handler) string {
	t.Helper()

	recorder := performJSONRequest(t, router, http.MethodPost, "/api/auth/wechat/login", "", map[string]any{
		"code":   "wx-code-001",
		"app_id": "wx-test-app",
	})
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected login status 200, got %d: %s", recorder.Code, recorder.Body.String())
	}

	var payload struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal login payload failed: %v", err)
	}
	if payload.Data.Token == "" {
		t.Fatalf("expected non-empty token, got payload %s", recorder.Body.String())
	}
	return payload.Data.Token
}

func performJSONRequest(t *testing.T, router http.Handler, method, url, token string, body any) *httptest.ResponseRecorder {
	t.Helper()

	var payload []byte
	var err error
	if body != nil {
		payload, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body failed: %v", err)
		}
	}

	request := httptest.NewRequest(method, url, bytes.NewReader(payload))
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	return recorder
}
