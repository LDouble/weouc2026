package transport

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/auth"
)

func RegisterRoutes(engine *gin.Engine, handler *Handler) {
	api := engine.Group("/api")

	user := api.Group("/notification")
	user.Use(auth.RequireAuthenticated())
	user.GET("/list", handler.ListMessages)
	user.GET("/unread-count", handler.GetUnreadCount)
	user.POST("/read", handler.MarkRead)

	admin := api.Group("/admin/notification")
	admin.Use(auth.RequireAuthenticated(), auth.RequirePermission("notification:publish"))
	admin.POST("/publish", handler.PublishMessage)
}

func parsePositiveInt(raw string, defaultValue int) int {
	value, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || value <= 0 {
		return defaultValue
	}
	return value
}

func parseBool(raw string) bool {
	value, err := strconv.ParseBool(strings.TrimSpace(raw))
	return err == nil && value
}
