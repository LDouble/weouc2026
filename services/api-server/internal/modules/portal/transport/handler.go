package transport

import (
	"strings"

	"github.com/gin-gonic/gin"
	portalservice "github.com/liangluo/weouc2026/services/api-server/internal/modules/portal/service"
	portaltypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/portal/types"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/auth"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/httpx"
)

type Handler struct {
	service *portalservice.Service
}

func NewHandler(service *portalservice.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetHome(c *gin.Context) {
	response, err := h.service.GetHome(c.Request.Context())
	if err != nil {
		httpx.AbortWithError(c, err)
		return
	}
	httpx.JSON(c, 200, response)
}

func (h *Handler) ListNotices(c *gin.Context) {
	query := portaltypes.NoticeQuery{
		Page:     parsePositiveInt(c.Query("page"), 1),
		PageSize: parsePositiveInt(c.Query("pageSize"), 20),
		Keyword:  strings.TrimSpace(c.Query("keyword")),
	}
	response, err := h.service.ListNotices(c.Request.Context(), query)
	if err != nil {
		httpx.AbortWithError(c, err)
		return
	}
	httpx.JSON(c, 200, response)
}

func (h *Handler) GetNotice(c *gin.Context) {
	response, err := h.service.GetNotice(c.Request.Context(), c.Param("id"))
	if err != nil {
		httpx.AbortWithError(c, err)
		return
	}
	httpx.JSON(c, 200, response)
}

func (h *Handler) PublishNotice(c *gin.Context) {
	var request portaltypes.NoticePublishRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpx.AbortWithError(c, httpx.BadRequest("公告发布参数格式错误", nil))
		return
	}
	response, err := h.service.PublishNotice(c.Request.Context(), auth.PrincipalFromContext(c), request)
	if err != nil {
		httpx.AbortWithError(c, err)
		return
	}
	httpx.JSON(c, 200, response)
}
