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

func TestPortalAndNotificationAPIs(t *testing.T) {
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

	t.Run("portal home is public and notice publish is permission guarded", func(t *testing.T) {
		recorder := performJSONRequest(t, router, http.MethodGet, "/api/portal/home", "", nil)
		if recorder.Code != http.StatusOK {
			t.Fatalf("expected portal home 200, got %d: %s", recorder.Code, recorder.Body.String())
		}
		if !strings.Contains(recorder.Body.String(), "\"banners\"") || !strings.Contains(recorder.Body.String(), "\"notices\"") {
			t.Fatalf("expected portal home payload, got %s", recorder.Body.String())
		}

		recorder = performJSONRequest(t, router, http.MethodPost, "/api/admin/portal/notices/publish", token, map[string]any{
			"title":   "普通学生不能发布公告",
			"content": "没有 portal:publish 权限时应被拒绝。",
		})
		if recorder.Code != http.StatusForbidden {
			t.Fatalf("expected portal publish 403 without permission, got %d: %s", recorder.Code, recorder.Body.String())
		}

		recorder = performJSONRequestWithHeaders(t, router, http.MethodPost, "/api/admin/portal/notices/publish", map[string]string{
			"Content-Type":       "application/json",
			"X-User-ID":          "admin-001",
			"X-User-Permissions": "portal:publish",
		}, map[string]any{
			"title":    "明日停机维护通知",
			"summary":  "今晚 23:00 至明日 01:00 将进行数据库迁移。",
			"content":  "为保障后端迁移执行，今晚 23:00 至明日 01:00 暂停部分内容发布能力。",
			"audience": "all",
			"tags":     []string{"运维", "通知"},
			"pinned":   true,
		})
		if recorder.Code != http.StatusOK {
			t.Fatalf("expected admin portal publish 200, got %d: %s", recorder.Code, recorder.Body.String())
		}

		var publishPayload struct {
			Data struct {
				ID string `json:"id"`
			} `json:"data"`
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &publishPayload); err != nil {
			t.Fatalf("unmarshal portal publish payload failed: %v", err)
		}
		if publishPayload.Data.ID == "" {
			t.Fatalf("expected portal notice id, got %s", recorder.Body.String())
		}

		recorder = performJSONRequest(t, router, http.MethodGet, "/api/portal/notices", "", nil)
		if recorder.Code != http.StatusOK {
			t.Fatalf("expected portal list 200, got %d: %s", recorder.Code, recorder.Body.String())
		}
		if !strings.Contains(recorder.Body.String(), publishPayload.Data.ID) {
			t.Fatalf("expected portal list to contain new notice, got %s", recorder.Body.String())
		}

		recorder = performJSONRequest(t, router, http.MethodGet, "/api/portal/notices/"+publishPayload.Data.ID, "", nil)
		if recorder.Code != http.StatusOK {
			t.Fatalf("expected portal detail 200, got %d: %s", recorder.Code, recorder.Body.String())
		}
		if !strings.Contains(recorder.Body.String(), "数据库迁移") {
			t.Fatalf("expected portal detail content, got %s", recorder.Body.String())
		}
	})

	t.Run("notification publish list and read workflow", func(t *testing.T) {
		recorder := performJSONRequest(t, router, http.MethodGet, "/api/notification/list", "", nil)
		if recorder.Code != http.StatusUnauthorized {
			t.Fatalf("expected notification list 401, got %d: %s", recorder.Code, recorder.Body.String())
		}

		recorder = performJSONRequest(t, router, http.MethodGet, "/api/notification/unread-count", token, nil)
		if recorder.Code != http.StatusOK {
			t.Fatalf("expected unread count 200, got %d: %s", recorder.Code, recorder.Body.String())
		}
		var unreadBefore struct {
			Data struct {
				Count int `json:"count"`
			} `json:"data"`
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &unreadBefore); err != nil {
			t.Fatalf("unmarshal unread count before failed: %v", err)
		}

		recorder = performJSONRequest(t, router, http.MethodPost, "/api/admin/notification/publish", token, map[string]any{
			"title":   "普通学生不能发通知",
			"content": "没有 notification:publish 权限时应被拒绝。",
		})
		if recorder.Code != http.StatusForbidden {
			t.Fatalf("expected notification publish 403 without permission, got %d: %s", recorder.Code, recorder.Body.String())
		}

		recorder = performJSONRequestWithHeaders(t, router, http.MethodPost, "/api/admin/notification/publish", map[string]string{
			"Content-Type":       "application/json",
			"X-User-ID":          "admin-001",
			"X-User-Permissions": "notification:publish",
		}, map[string]any{
			"title":        "今晚 23 点起暂停内容发布",
			"content":      "发布、审核和消息写入链路将短暂进入只读维护窗口。",
			"category":     "system",
			"target_scope": "all",
			"action_url":   "/pages/home/index",
		})
		if recorder.Code != http.StatusOK {
			t.Fatalf("expected admin notification publish 200, got %d: %s", recorder.Code, recorder.Body.String())
		}

		var publishPayload struct {
			Data struct {
				ID string `json:"id"`
			} `json:"data"`
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &publishPayload); err != nil {
			t.Fatalf("unmarshal notification publish payload failed: %v", err)
		}
		if publishPayload.Data.ID == "" {
			t.Fatalf("expected notification id, got %s", recorder.Body.String())
		}

		recorder = performJSONRequest(t, router, http.MethodGet, "/api/notification/unread-count", token, nil)
		if recorder.Code != http.StatusOK {
			t.Fatalf("expected unread count after publish 200, got %d: %s", recorder.Code, recorder.Body.String())
		}
		var unreadAfterPublish struct {
			Data struct {
				Count int `json:"count"`
			} `json:"data"`
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &unreadAfterPublish); err != nil {
			t.Fatalf("unmarshal unread count after publish failed: %v", err)
		}
		if unreadAfterPublish.Data.Count != unreadBefore.Data.Count+1 {
			t.Fatalf("expected unread count +1 after publish, before=%d after=%d", unreadBefore.Data.Count, unreadAfterPublish.Data.Count)
		}

		recorder = performJSONRequest(t, router, http.MethodGet, "/api/notification/list", token, nil)
		if recorder.Code != http.StatusOK {
			t.Fatalf("expected notification list 200, got %d: %s", recorder.Code, recorder.Body.String())
		}
		if !strings.Contains(recorder.Body.String(), publishPayload.Data.ID) {
			t.Fatalf("expected notification list to contain new message, got %s", recorder.Body.String())
		}
		if !strings.Contains(recorder.Body.String(), "\"read\":false") {
			t.Fatalf("expected unread notification state, got %s", recorder.Body.String())
		}

		recorder = performJSONRequest(t, router, http.MethodPost, "/api/notification/read", token, map[string]any{
			"message_id": publishPayload.Data.ID,
		})
		if recorder.Code != http.StatusOK {
			t.Fatalf("expected notification read 200, got %d: %s", recorder.Code, recorder.Body.String())
		}

		recorder = performJSONRequest(t, router, http.MethodGet, "/api/notification/unread-count", token, nil)
		if recorder.Code != http.StatusOK {
			t.Fatalf("expected unread count after read 200, got %d: %s", recorder.Code, recorder.Body.String())
		}
		var unreadAfterRead struct {
			Data struct {
				Count int `json:"count"`
			} `json:"data"`
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &unreadAfterRead); err != nil {
			t.Fatalf("unmarshal unread count after read failed: %v", err)
		}
		if unreadAfterRead.Data.Count != unreadBefore.Data.Count {
			t.Fatalf("expected unread count restored after read, before=%d after=%d", unreadBefore.Data.Count, unreadAfterRead.Data.Count)
		}

		recorder = performJSONRequest(t, router, http.MethodGet, "/api/notification/list", token, nil)
		if recorder.Code != http.StatusOK {
			t.Fatalf("expected notification list after read 200, got %d: %s", recorder.Code, recorder.Body.String())
		}
		if !strings.Contains(recorder.Body.String(), publishPayload.Data.ID) || !strings.Contains(recorder.Body.String(), "\"read\":true") {
			t.Fatalf("expected notification to become read, got %s", recorder.Body.String())
		}
	})
}
