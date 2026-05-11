package transport

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/auth"
)

func RegisterRoutes(engine *gin.Engine, handler *Handler) {
	admin := engine.Group("/api/admin/analytics")
	admin.Use(auth.RequireAuthenticated(), auth.RequirePermission("analytics:view"))
	admin.GET("/dashboard", handler.GetDashboard)
	admin.GET("/audit-logs", handler.ListAuditLogs)
}

func parsePositiveInt(raw string, defaultValue int) int {
	value, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || value <= 0 {
		return defaultValue
	}
	return value
}
