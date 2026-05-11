package transport

import (
	"strings"

	"github.com/gin-gonic/gin"
	notificationservice "github.com/liangluo/weouc2026/services/api-server/internal/modules/notification/service"
	notificationtypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/notification/types"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/auth"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/httpx"
)

type Handler struct {
	service *notificationservice.Service
}

func NewHandler(service *notificationservice.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListMessages(c *gin.Context) {
	query := notificationtypes.MessageQuery{
		Page:       parsePositiveInt(c.Query("page"), 1),
		PageSize:   parsePositiveInt(c.Query("pageSize"), 20),
		Category:   strings.TrimSpace(c.Query("category")),
		UnreadOnly: parseBool(c.Query("unread_only")),
	}
	response, err := h.service.ListMessages(c.Request.Context(), auth.PrincipalFromContext(c), query)
	if err != nil {
		httpx.AbortWithError(c, err)
		return
	}
	httpx.JSON(c, 200, response)
}

func (h *Handler) GetUnreadCount(c *gin.Context) {
	response, err := h.service.GetUnreadCount(c.Request.Context(), auth.PrincipalFromContext(c))
	if err != nil {
		httpx.AbortWithError(c, err)
		return
	}
	httpx.JSON(c, 200, response)
}

func (h *Handler) PublishMessage(c *gin.Context) {
	var request notificationtypes.PublishRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpx.AbortWithError(c, httpx.BadRequest("通知发布参数格式错误", nil))
		return
	}
	response, err := h.service.PublishMessage(c.Request.Context(), auth.PrincipalFromContext(c), request)
	if err != nil {
		httpx.AbortWithError(c, err)
		return
	}
	httpx.JSON(c, 200, response)
}

func (h *Handler) MarkRead(c *gin.Context) {
	var request notificationtypes.MarkReadRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpx.AbortWithError(c, httpx.BadRequest("通知已读参数格式错误", nil))
		return
	}
	if err := h.service.MarkRead(c.Request.Context(), auth.PrincipalFromContext(c), request.MessageID); err != nil {
		httpx.AbortWithError(c, err)
		return
	}
	httpx.JSON(c, 200, gin.H{"success": true})
}
