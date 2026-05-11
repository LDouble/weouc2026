package transport

import (
	"github.com/gin-gonic/gin"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/auth"
)

func RegisterRoutes(engine *gin.Engine, handler *Handler) {
	apiGroup := engine.Group("/api")
	apiGroup.POST("/auth/wechat/login", handler.LoginWithWeChat)
	apiGroup.POST("/auth/admin/login", handler.LoginWithPassword)

	protected := apiGroup.Group("")
	protected.Use(auth.RequireAuthenticated())
	protected.GET("/student", handler.GetStudentProfile)
	protected.POST("/student", handler.BindStudent)
	protected.PUT("/student", handler.UpdateStudent)
	protected.POST("/edu/send-captcha", handler.SendCaptcha)
}
