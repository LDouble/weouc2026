package auth

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	appconfig "github.com/liangluo/weouc2026/services/api-server/internal/platform/config"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/httpx"
)

const principalContextKey = "auth_principal"

type Principal struct {
	Authenticated bool     `json:"authenticated"`
	UserID        string   `json:"user_id,omitempty"`
	Roles         []string `json:"roles,omitempty"`
	Permissions   []string `json:"permissions,omitempty"`
	AcademicBound bool     `json:"academic_bound"`
}

func ContextMiddleware(cfg appconfig.AppConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		principal := Principal{
			UserID:        strings.TrimSpace(c.GetHeader(cfg.Auth.UserIDHeader)),
			Roles:         splitHeaderValues(c.GetHeader(cfg.Auth.RolesHeader)),
			Permissions:   splitHeaderValues(c.GetHeader(cfg.Auth.PermissionsHeader)),
			AcademicBound: parseBoolHeader(c.GetHeader(cfg.Auth.AcademicBoundHeader)),
		}
		principal.Authenticated = principal.UserID != ""

		c.Set(principalContextKey, principal)
		c.Next()
	}
}

func RequireAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !PrincipalFromContext(c).Authenticated {
			httpx.AbortWithError(c, httpx.Unauthorized("需要登录后访问"))
			return
		}

		c.Next()
	}
}

func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		principal := PrincipalFromContext(c)
		if !principal.Authenticated {
			httpx.AbortWithError(c, httpx.Unauthorized("需要登录后访问"))
			return
		}
		if !principal.HasPermission(permission) {
			httpx.AbortWithError(c, httpx.Forbidden("当前账号缺少所需权限", map[string]any{
				"required_permission": permission,
			}))
			return
		}

		c.Next()
	}
}

func PrincipalFromContext(c *gin.Context) Principal {
	value, exists := c.Get(principalContextKey)
	if !exists {
		return Principal{}
	}

	principal, ok := value.(Principal)
	if !ok {
		return Principal{}
	}

	return principal
}

func (p Principal) HasPermission(permission string) bool {
	if permission == "" {
		return true
	}

	for _, candidate := range p.Permissions {
		if candidate == permission || candidate == "*" {
			return true
		}
	}

	return false
}

func splitHeaderValues(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value != "" {
			result = append(result, value)
		}
	}

	return result
}

func parseBoolHeader(raw string) bool {
	value, err := strconv.ParseBool(strings.TrimSpace(raw))
	return err == nil && value
}
