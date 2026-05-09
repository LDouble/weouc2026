package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	appconfig "github.com/liangluo/weouc2026/services/api-server/internal/platform/config"
)

func TestContextMiddlewareParsesHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := appconfig.AppConfig{
		Auth: appconfig.AuthConfig{
			UserIDHeader:        "X-User-ID",
			RolesHeader:         "X-User-Roles",
			PermissionsHeader:   "X-User-Permissions",
			AcademicBoundHeader: "X-Academic-Bound",
		},
	}

	router := gin.New()
	router.Use(ContextMiddleware(cfg))
	router.GET("/principal", func(c *gin.Context) {
		c.JSON(http.StatusOK, PrincipalFromContext(c))
	})

	request := httptest.NewRequest(http.MethodGet, "/principal", nil)
	request.Header.Set("X-User-ID", "u-1001")
	request.Header.Set("X-User-Roles", "student, editor")
	request.Header.Set("X-User-Permissions", "portal:publish, contact:view")
	request.Header.Set("X-Academic-Bound", "true")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	var principal Principal
	if err := json.Unmarshal(recorder.Body.Bytes(), &principal); err != nil {
		t.Fatalf("unmarshal principal failed: %v", err)
	}

	if !principal.Authenticated || principal.UserID != "u-1001" {
		t.Fatalf("unexpected principal: %+v", principal)
	}
	if len(principal.Roles) != 2 || !principal.AcademicBound {
		t.Fatalf("unexpected roles or academic bound flag: %+v", principal)
	}
}

func TestRequirePermissionRejectsMissingPermission(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := appconfig.AppConfig{
		Auth: appconfig.AuthConfig{
			UserIDHeader:        "X-User-ID",
			RolesHeader:         "X-User-Roles",
			PermissionsHeader:   "X-User-Permissions",
			AcademicBoundHeader: "X-Academic-Bound",
		},
	}

	router := gin.New()
	router.Use(ContextMiddleware(cfg))
	router.GET("/secure", RequirePermission("portal:publish"), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	request := httptest.NewRequest(http.MethodGet, "/secure", nil)
	request.Header.Set("X-User-ID", "u-1001")
	request.Header.Set("X-User-Permissions", "contact:view")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", recorder.Code)
	}
}
