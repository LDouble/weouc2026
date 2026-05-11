package transport

import (
	"github.com/gin-gonic/gin"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/auth"
)

func RegisterRoutes(engine *gin.Engine, handler *Handler) {
	group := engine.Group("/api/academic")
	group.Use(auth.RequireAuthenticated())
	group.GET("/semesters", handler.ListSemesters)
	group.GET("/schedule", handler.GetSchedule)
	group.GET("/exams", handler.ListExams)
	group.GET("/grades", handler.ListGrades)
}
