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

func TestAcademicAPIs(t *testing.T) {
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

	recorder := performJSONRequest(t, router, http.MethodGet, "/api/academic/semesters", "", nil)
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected academic semesters 401 without auth, got %d: %s", recorder.Code, recorder.Body.String())
	}

	token := loginAndGetToken(t, router)

	recorder = performJSONRequest(t, router, http.MethodGet, "/api/academic/semesters", token, nil)
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected academic semesters 403 before bind, got %d: %s", recorder.Code, recorder.Body.String())
	}

	recorder = performJSONRequest(t, router, http.MethodPost, "/api/edu/send-captcha", token, map[string]any{
		"sid": "20260002",
	})
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected captcha send 200, got %d: %s", recorder.Code, recorder.Body.String())
	}

	recorder = performJSONRequest(t, router, http.MethodPost, "/api/student", token, map[string]any{
		"student_id": "20260002",
		"password":   "password-002",
		"captcha":    "123456",
	})
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected bind student 200, got %d: %s", recorder.Code, recorder.Body.String())
	}

	recorder = performJSONRequest(t, router, http.MethodGet, "/api/academic/semesters", token, nil)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected academic semesters 200, got %d: %s", recorder.Code, recorder.Body.String())
	}

	var semesters struct {
		Data struct {
			CurrentSemesterID string `json:"current_semester_id"`
			List              []struct {
				ID        string `json:"id"`
				IsCurrent bool   `json:"is_current"`
			} `json:"list"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &semesters); err != nil {
		t.Fatalf("unmarshal semesters payload failed: %v", err)
	}
	if semesters.Data.CurrentSemesterID == "" || len(semesters.Data.List) == 0 {
		t.Fatalf("expected non-empty semesters payload, got %s", recorder.Body.String())
	}

	recorder = performJSONRequest(t, router, http.MethodGet, "/api/academic/schedule", token, nil)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected academic schedule 200, got %d: %s", recorder.Code, recorder.Body.String())
	}

	var schedule struct {
		Data struct {
			Semester struct {
				ID string `json:"id"`
			} `json:"semester"`
			List []struct {
				ID string `json:"id"`
			} `json:"list"`
			Summary struct {
				CourseCount int `json:"course_count"`
			} `json:"summary"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &schedule); err != nil {
		t.Fatalf("unmarshal schedule payload failed: %v", err)
	}
	if schedule.Data.Semester.ID != semesters.Data.CurrentSemesterID {
		t.Fatalf("expected current semester schedule, got %s", recorder.Body.String())
	}
	if schedule.Data.Summary.CourseCount == 0 || len(schedule.Data.List) == 0 {
		t.Fatalf("expected non-empty schedule, got %s", recorder.Body.String())
	}

	recorder = performJSONRequest(
		t,
		router,
		http.MethodGet,
		"/api/academic/exams?semester_id="+semesters.Data.CurrentSemesterID,
		token,
		nil,
	)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected academic exams 200, got %d: %s", recorder.Code, recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), "\"semester_id\":\""+semesters.Data.CurrentSemesterID+"\"") {
		t.Fatalf("expected exams payload to contain selected semester, got %s", recorder.Body.String())
	}

	recorder = performJSONRequest(t, router, http.MethodGet, "/api/academic/grades", token, nil)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected academic grades 200, got %d: %s", recorder.Code, recorder.Body.String())
	}

	var grades struct {
		Data struct {
			Summary struct {
				CourseCount       int     `json:"course_count"`
				PassedCount       int     `json:"passed_count"`
				AverageScore      float64 `json:"average_score"`
				AverageGradePoint float64 `json:"average_grade_point"`
			} `json:"summary"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &grades); err != nil {
		t.Fatalf("unmarshal grades payload failed: %v", err)
	}
	if grades.Data.Summary.CourseCount == 0 || grades.Data.Summary.PassedCount == 0 {
		t.Fatalf("expected non-empty grades summary, got %s", recorder.Body.String())
	}
	if grades.Data.Summary.AverageScore <= 0 || grades.Data.Summary.AverageGradePoint <= 0 {
		t.Fatalf("expected positive grade averages, got %s", recorder.Body.String())
	}

	recorder = performJSONRequestWithHeaders(t, router, http.MethodGet, "/api/admin/analytics/audit-logs?action=academic.grades.view", map[string]string{
		"X-User-ID":          "admin-001",
		"X-User-Permissions": "analytics:view",
	}, nil)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected analytics audit logs 200, got %d: %s", recorder.Code, recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), "\"student_id\":\"20260002\"") || !strings.Contains(recorder.Body.String(), "\"action\":\"academic.grades.view\"") {
		t.Fatalf("expected academic grades audit log payload, got %s", recorder.Body.String())
	}
}
