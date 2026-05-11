package transport

import (
	"strings"

	"github.com/gin-gonic/gin"
	analyticsservice "github.com/liangluo/weouc2026/services/api-server/internal/modules/analytics/service"
	analyticstypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/analytics/types"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/httpx"
)

type Handler struct {
	service *analyticsservice.Service
}

func NewHandler(service *analyticsservice.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetDashboard(c *gin.Context) {
	response, err := h.service.GetDashboard(c.Request.Context())
	if err != nil {
		httpx.AbortWithError(c, err)
		return
	}
	httpx.JSON(c, 200, response)
}

func (h *Handler) ListAuditLogs(c *gin.Context) {
	query := analyticstypes.AuditLogQuery{
		Page:         parsePositiveInt(c.Query("page"), 1),
		PageSize:     parsePositiveInt(c.Query("pageSize"), 20),
		ActorID:      strings.TrimSpace(c.Query("actor_id")),
		Action:       strings.TrimSpace(c.Query("action")),
		ResourceType: strings.TrimSpace(c.Query("resource_type")),
		ResourceID:   strings.TrimSpace(c.Query("resource_id")),
	}
	response, err := h.service.ListAuditLogs(c.Request.Context(), query)
	if err != nil {
		httpx.AbortWithError(c, err)
		return
	}
	httpx.JSON(c, 200, response)
}
