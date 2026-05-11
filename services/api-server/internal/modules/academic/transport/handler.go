package transport

import (
	"strings"

	"github.com/gin-gonic/gin"
	academicservice "github.com/liangluo/weouc2026/services/api-server/internal/modules/academic/service"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/auth"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/httpx"
)

type Handler struct {
	service *academicservice.Service
}

func NewHandler(service *academicservice.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListSemesters(c *gin.Context) {
	response, err := h.service.ListSemesters(c.Request.Context(), auth.PrincipalFromContext(c))
	if err != nil {
		httpx.AbortWithError(c, err)
		return
	}
	httpx.JSON(c, 200, response)
}

func (h *Handler) GetSchedule(c *gin.Context) {
	response, err := h.service.GetSchedule(
		c.Request.Context(),
		auth.PrincipalFromContext(c),
		strings.TrimSpace(c.Query("semester_id")),
	)
	if err != nil {
		httpx.AbortWithError(c, err)
		return
	}
	httpx.JSON(c, 200, response)
}

func (h *Handler) ListExams(c *gin.Context) {
	response, err := h.service.ListExams(
		c.Request.Context(),
		auth.PrincipalFromContext(c),
		strings.TrimSpace(c.Query("semester_id")),
	)
	if err != nil {
		httpx.AbortWithError(c, err)
		return
	}
	httpx.JSON(c, 200, response)
}

func (h *Handler) ListGrades(c *gin.Context) {
	response, err := h.service.ListGrades(
		c.Request.Context(),
		auth.PrincipalFromContext(c),
		strings.TrimSpace(c.Query("semester_id")),
	)
	if err != nil {
		httpx.AbortWithError(c, err)
		return
	}
	httpx.JSON(c, 200, response)
}
