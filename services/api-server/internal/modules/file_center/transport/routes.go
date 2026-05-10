package transport

import (
	"github.com/gin-gonic/gin"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/auth"
)

func RegisterRoutes(engine *gin.Engine, handler *Handler) {
	api := engine.Group("/api")
	protected := api.Group("")
	protected.Use(auth.RequireAuthenticated())
	protected.GET("/upload/cos-sts", handler.GetCOSSTS)
	protected.POST("/upload/presigned-get", handler.GetPresignedGET)
}
