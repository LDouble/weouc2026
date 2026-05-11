package transport

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/auth"
)

func RegisterRoutes(engine *gin.Engine, handler *Handler) {
	api := engine.Group("/api")
	api.GET("/portal/home", handler.GetHome)
	api.GET("/portal/notices", handler.ListNotices)
	api.GET("/portal/notices/:id", handler.GetNotice)

	admin := api.Group("/admin/portal")
	admin.Use(auth.RequireAuthenticated(), auth.RequirePermission("portal:publish"))
	admin.GET("/banners", handler.ListBanners)
	admin.GET("/banners/:id", handler.GetBanner)
	admin.POST("/banners", handler.CreateBanner)
	admin.PUT("/banners/:id", handler.UpdateBanner)
	admin.DELETE("/banners/:id", handler.DeleteBanner)
	admin.POST("/notices/publish", handler.PublishNotice)
	admin.PUT("/notices/:id", handler.UpdateNotice)
	admin.DELETE("/notices/:id", handler.DeleteNotice)
}

func parsePositiveInt(raw string, defaultValue int) int {
	value, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || value <= 0 {
		return defaultValue
	}
	return value
}
