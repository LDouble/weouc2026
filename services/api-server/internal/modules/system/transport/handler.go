package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"
	systemservice "github.com/liangluo/weouc2026/services/api-server/internal/modules/system/service"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/auth"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/httpx"
)

type Handler struct {
	service *systemservice.Service
}

func NewHandler(service *systemservice.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Healthz(c *gin.Context) {
	httpx.JSON(c, http.StatusOK, h.service.Health(c.Request.Context()))
}

func (h *Handler) Readyz(c *gin.Context) {
	status := h.service.Ready(c.Request.Context())
	statusCode := http.StatusOK
	if !status.IsReady() {
		statusCode = http.StatusServiceUnavailable
	}

	httpx.JSON(c, statusCode, status)
}

func (h *Handler) Profile(c *gin.Context) {
	principal := auth.PrincipalFromContext(c)
	httpx.JSON(c, http.StatusOK, h.service.Profile(c.Request.Context(), principal))
}
