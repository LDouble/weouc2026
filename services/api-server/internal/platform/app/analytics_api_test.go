package app

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	appconfig "github.com/liangluo/weouc2026/services/api-server/internal/platform/config"
	applogger "github.com/liangluo/weouc2026/services/api-server/internal/platform/logger"
)

func TestAnalyticsAuditAPIs(t *testing.T) {
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

	recorder := performJSONRequest(t, router, http.MethodGet, "/api/admin/analytics/dashboard", token, nil)
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected analytics dashboard 403 without permission, got %d: %s", recorder.Code, recorder.Body.String())
	}

	recorder = performJSONRequest(t, router, http.MethodPost, "/api/edu/send-captcha", token, map[string]any{
		"sid": "20260001",
	})
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected captcha send 200, got %d: %s", recorder.Code, recorder.Body.String())
	}

	recorder = performJSONRequest(t, router, http.MethodPost, "/api/student", token, map[string]any{
		"student_id": "20260001",
		"password":   "password-001",
		"captcha":    "123456",
	})
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected bind student 200, got %d: %s", recorder.Code, recorder.Body.String())
	}

	recorder = performJSONRequest(t, router, http.MethodPost, "/api/market/publish", token, map[string]any{
		"title":   "九成新台灯",
		"desc":    "宿舍搬空出一个白色台灯。",
		"price":   "35",
		"contact": "wx-market-a1",
	})
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected market publish 200, got %d: %s", recorder.Code, recorder.Body.String())
	}
	var marketPublish struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &marketPublish); err != nil {
		t.Fatalf("unmarshal market publish payload failed: %v", err)
	}

	recorder = performJSONRequestWithHeaders(t, router, http.MethodPost, "/api/admin/portal/notices/publish", map[string]string{
		"Content-Type":       "application/json",
		"X-User-ID":          "admin-001",
		"X-User-Permissions": "portal:publish",
	}, map[string]any{
		"title":   "维护窗口通知",
		"content": "今晚 23:00 至明日 01:00 进行例行维护。",
		"pinned":  true,
	})
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected portal publish 200, got %d: %s", recorder.Code, recorder.Body.String())
	}

	recorder = performJSONRequestWithHeaders(t, router, http.MethodPost, "/api/admin/notification/publish", map[string]string{
		"Content-Type":       "application/json",
		"X-User-ID":          "admin-001",
		"X-User-Permissions": "notification:publish",
	}, map[string]any{
		"title":        "系统维护提醒",
		"content":      "发布与审核链路将短暂进入只读窗口。",
		"category":     "system",
		"target_scope": "all",
	})
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected notification publish 200, got %d: %s", recorder.Code, recorder.Body.String())
	}

	recorder = performJSONRequestWithHeaders(t, router, http.MethodPost, "/api/admin/campus-life/review/update", map[string]string{
		"Content-Type":       "application/json",
		"X-User-ID":          "admin-001",
		"X-User-Permissions": "campus_life:moderate",
	}, map[string]any{
		"content_type":  "market",
		"content_id":    marketPublish.Data.ID,
		"review_status": "published",
	})
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected review update 200, got %d: %s", recorder.Code, recorder.Body.String())
	}

	recorder = performJSONRequestWithHeaders(t, router, http.MethodGet, "/api/admin/analytics/dashboard", map[string]string{
		"X-User-ID":          "admin-001",
		"X-User-Permissions": "analytics:view",
	}, nil)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected analytics dashboard 200, got %d: %s", recorder.Code, recorder.Body.String())
	}

	var dashboard struct {
		Data struct {
			Summary struct {
				TotalAuditLogs int `json:"total_audit_logs"`
				LoginCount     int `json:"login_count"`
				BindCount      int `json:"bind_count"`
				PublishCount   int `json:"publish_count"`
				ReviewCount    int `json:"review_count"`
			} `json:"summary"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &dashboard); err != nil {
		t.Fatalf("unmarshal analytics dashboard failed: %v", err)
	}
	if dashboard.Data.Summary.TotalAuditLogs < 5 {
		t.Fatalf("expected at least 5 audit logs, got %d: %s", dashboard.Data.Summary.TotalAuditLogs, recorder.Body.String())
	}
	if dashboard.Data.Summary.LoginCount < 1 || dashboard.Data.Summary.BindCount < 1 {
		t.Fatalf("expected login/bind counts recorded, got %s", recorder.Body.String())
	}
	if dashboard.Data.Summary.PublishCount < 3 || dashboard.Data.Summary.ReviewCount < 1 {
		t.Fatalf("expected publish/review counts recorded, got %s", recorder.Body.String())
	}

	recorder = performJSONRequestWithHeaders(t, router, http.MethodGet, "/api/admin/analytics/audit-logs?action=campus_life.review.update", map[string]string{
		"X-User-ID":          "admin-001",
		"X-User-Permissions": "analytics:view",
	}, nil)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected analytics audit logs 200, got %d: %s", recorder.Code, recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), marketPublish.Data.ID) || !strings.Contains(recorder.Body.String(), "\"review_status\":\"published\"") {
		t.Fatalf("expected review audit log payload, got %s", recorder.Body.String())
	}
}
