package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	appconfig "github.com/liangluo/weouc2026/services/api-server/internal/platform/config"
	applogger "github.com/liangluo/weouc2026/services/api-server/internal/platform/logger"
)

func TestMiniappCoreAPIs(t *testing.T) {
	t.Skip("requires mysql + mongo + redis integration environment")

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

	t.Run("carpool review workflow gates public visibility", func(t *testing.T) {
		recorder := performJSONRequest(t, router, http.MethodPost, "/api/carpool/publish", token, map[string]any{
			"category":    "tomorrow",
			"travel_date": "2026-05-12",
			"travel_time": "18:30",
			"from":        "海大南门",
			"to":          "福州南站",
			"type":        "明日顺路",
			"seats_text":  "余座 2",
			"price":       "人均 20 元",
			"note":        "可带一个行李箱",
			"tags":        []string{"明天出发"},
			"contact":     "wx-carpool-301",
		})
		if recorder.Code != http.StatusOK {
			t.Fatalf("expected carpool publish status 200, got %d: %s", recorder.Code, recorder.Body.String())
		}

		var publishPayload struct {
			Data struct {
				ID string `json:"id"`
			} `json:"data"`
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &publishPayload); err != nil {
			t.Fatalf("unmarshal carpool publish payload failed: %v", err)
		}
		if publishPayload.Data.ID == "" {
			t.Fatalf("expected carpool id, got %s", recorder.Body.String())
		}

		detailURL := "/api/carpool/detail/" + publishPayload.Data.ID
		recorder = performJSONRequest(t, router, http.MethodGet, detailURL, "", nil)
		if recorder.Code != http.StatusNotFound {
			t.Fatalf("expected public pending carpool detail 404, got %d: %s", recorder.Code, recorder.Body.String())
		}

		recorder = performJSONRequest(t, router, http.MethodGet, detailURL, token, nil)
		if recorder.Code != http.StatusOK {
			t.Fatalf("expected owner pending carpool detail 200, got %d: %s", recorder.Code, recorder.Body.String())
		}

		var ownerDetail struct {
			Data struct {
				Status string `json:"status"`
			} `json:"data"`
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &ownerDetail); err != nil {
			t.Fatalf("unmarshal owner carpool detail failed: %v", err)
		}
		if ownerDetail.Data.Status != "reviewing" {
			t.Fatalf("expected reviewing status, got %s", recorder.Body.String())
		}

		recorder = performJSONRequest(t, router, http.MethodGet, "/api/carpool/list", "", nil)
		if recorder.Code != http.StatusOK {
			t.Fatalf("expected public carpool list 200, got %d: %s", recorder.Code, recorder.Body.String())
		}
		if strings.Contains(recorder.Body.String(), publishPayload.Data.ID) {
			t.Fatalf("expected pending carpool hidden from public list: %s", recorder.Body.String())
		}

		recorder = performJSONRequestWithHeaders(t, router, http.MethodPost, "/api/admin/campus-life/review/update", map[string]string{
			"Content-Type":       "application/json",
			"X-User-ID":          "admin-001",
			"X-User-Permissions": "campus_life:moderate",
		}, map[string]any{
			"content_type":  "carpool",
			"content_id":    publishPayload.Data.ID,
			"review_status": "published",
		})
		if recorder.Code != http.StatusOK {
			t.Fatalf("expected review update 200, got %d: %s", recorder.Code, recorder.Body.String())
		}

		recorder = performJSONRequest(t, router, http.MethodGet, "/api/carpool/list", "", nil)
		if recorder.Code != http.StatusOK {
			t.Fatalf("expected public carpool list after review 200, got %d: %s", recorder.Code, recorder.Body.String())
		}
		if !strings.Contains(recorder.Body.String(), publishPayload.Data.ID) {
			t.Fatalf("expected approved carpool visible in public list: %s", recorder.Body.String())
		}
	})

	t.Run("meetup review and join workflow", func(t *testing.T) {
		tomorrow := time.Now().In(time.FixedZone("Asia/Shanghai", 8*3600)).AddDate(0, 0, 1)
		startAt := tomorrow.Format("2006-01-02") + "T19:00:00+08:00"
		deadlineAt := tomorrow.Format("2006-01-02") + "T17:30:00+08:00"
		recorder := performJSONRequest(t, router, http.MethodPost, "/api/meetup/publish", token, map[string]any{
			"category":         "study",
			"title":            "高数晚自习组队",
			"desc":             "想找 2 位同学一起在图书馆刷题。",
			"location":         "图书馆五楼北区",
			"start_at":         startAt,
			"deadline_at":      deadlineAt,
			"max_participants": 3,
			"fee_text":         "免费",
			"tags":             []string{"期末复习", "刷题"},
			"contact":          "wx-meetup-301",
		})
		if recorder.Code != http.StatusOK {
			t.Fatalf("expected meetup publish status 200, got %d: %s", recorder.Code, recorder.Body.String())
		}

		var publishPayload struct {
			Data struct {
				ID string `json:"id"`
			} `json:"data"`
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &publishPayload); err != nil {
			t.Fatalf("unmarshal meetup publish payload failed: %v", err)
		}
		if publishPayload.Data.ID == "" {
			t.Fatalf("expected meetup id, got %s", recorder.Body.String())
		}

		detailURL := "/api/meetup/detail/" + publishPayload.Data.ID
		recorder = performJSONRequest(t, router, http.MethodGet, detailURL, "", nil)
		if recorder.Code != http.StatusNotFound {
			t.Fatalf("expected public pending meetup detail 404, got %d: %s", recorder.Code, recorder.Body.String())
		}

		recorder = performJSONRequestWithHeaders(t, router, http.MethodPost, "/api/admin/campus-life/review/update", map[string]string{
			"Content-Type":       "application/json",
			"X-User-ID":          "admin-001",
			"X-User-Permissions": "campus_life:moderate",
		}, map[string]any{
			"content_type":  "meetup",
			"content_id":    publishPayload.Data.ID,
			"review_status": "published",
		})
		if recorder.Code != http.StatusOK {
			t.Fatalf("expected meetup review update 200, got %d: %s", recorder.Code, recorder.Body.String())
		}

		recorder = performJSONRequestWithHeaders(t, router, http.MethodPost, "/api/meetup/join", map[string]string{
			"Content-Type": "application/json",
			"X-User-ID":    "participant-001",
		}, map[string]any{
			"meetup_id": publishPayload.Data.ID,
		})
		if recorder.Code != http.StatusOK {
			t.Fatalf("expected meetup join 200, got %d: %s", recorder.Code, recorder.Body.String())
		}

		recorder = performJSONRequestWithHeaders(t, router, http.MethodGet, detailURL, map[string]string{
			"X-User-ID":          "participant-001",
			"X-Academic-Bound":   "true",
			"X-User-Permissions": "",
		}, nil)
		if recorder.Code != http.StatusOK {
			t.Fatalf("expected participant meetup detail 200, got %d: %s", recorder.Code, recorder.Body.String())
		}

		var participantDetail struct {
			Data struct {
				Joined         bool   `json:"joined"`
				UserRole       string `json:"user_role"`
				CanCancelJoin  bool   `json:"can_cancel_join"`
				JoinedCount    int    `json:"joined_count"`
				RemainingSeats int    `json:"remaining_seats"`
			} `json:"data"`
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &participantDetail); err != nil {
			t.Fatalf("unmarshal participant meetup detail failed: %v", err)
		}
		if !participantDetail.Data.Joined || participantDetail.Data.UserRole != "participant" || !participantDetail.Data.CanCancelJoin {
			t.Fatalf("unexpected participant meetup detail payload: %s", recorder.Body.String())
		}
		if participantDetail.Data.JoinedCount != 2 || participantDetail.Data.RemainingSeats != 1 {
			t.Fatalf("unexpected meetup participant counters: %s", recorder.Body.String())
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

func performJSONRequestWithHeaders(
	t *testing.T,
	router http.Handler,
	method, url string,
	headers map[string]string,
	body any,
) *httptest.ResponseRecorder {
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
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	return recorder
}
