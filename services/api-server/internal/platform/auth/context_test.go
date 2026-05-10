package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	appconfig "github.com/liangluo/weouc2026/services/api-server/internal/platform/config"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/httpx"
)

type mockResolver struct {
	principal Principal
	err       error
}

func (m mockResolver) ResolveToken(context.Context, string) (Principal, error) {
	if m.err != nil {
		return Principal{}, m.err
	}

	return m.principal, nil
}

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
	router.Use(ContextMiddleware(cfg, nil))
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
	router.Use(ContextMiddleware(cfg, nil))
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

func TestContextMiddlewareResolvesBearerToken(t *testing.T) {
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
	router.Use(ContextMiddleware(cfg, mockResolver{
		principal: Principal{
			Authenticated: true,
			UserID:        "u-2002",
			DisplayName:   "海大同学",
			Roles:         []string{"student"},
			Permissions:   []string{"contact:view"},
			AcademicBound: true,
		},
	}))
	router.GET("/principal", func(c *gin.Context) {
		c.JSON(http.StatusOK, PrincipalFromContext(c))
	})

	request := httptest.NewRequest(http.MethodGet, "/principal", nil)
	request.Header.Set("Authorization", "Bearer token-001")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	var principal Principal
	if err := json.Unmarshal(recorder.Body.Bytes(), &principal); err != nil {
		t.Fatalf("unmarshal principal failed: %v", err)
	}

	if principal.UserID != "u-2002" || !principal.AcademicBound || principal.DisplayName != "海大同学" {
		t.Fatalf("unexpected principal: %+v", principal)
	}
}

func TestContextMiddlewareRejectsInvalidBearerToken(t *testing.T) {
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
	router.Use(ContextMiddleware(cfg, mockResolver{err: httpx.Unauthorized("登录状态已失效，请重新登录")}))
	router.GET("/principal", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	request := httptest.NewRequest(http.MethodGet, "/principal", nil)
	request.Header.Set("Authorization", "Bearer expired-token")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", recorder.Code)
	}
}
