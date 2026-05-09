package transport

import (
	"github.com/gin-gonic/gin"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/auth"
)

func RegisterRoutes(engine *gin.Engine, handler *Handler) {
	engine.GET("/healthz", handler.Healthz)
	engine.GET("/readyz", handler.Readyz)

	systemGroup := engine.Group("/api/v1/system")
	systemGroup.GET("/profile", auth.RequireAuthenticated(), handler.Profile)
}
