package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"
	fcservice "github.com/liangluo/weouc2026/services/api-server/internal/modules/file_center/service"
	fctypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/file_center/types"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/auth"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/httpx"
)

type Handler struct {
	service *fcservice.Service
}

func NewHandler(service *fcservice.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetCOSSTS(c *gin.Context) {
	response, err := h.service.IssueUploadCredentials(
		c.Request.Context(),
		auth.PrincipalFromContext(c),
		c.Query("scene"),
	)
	if err != nil {
		httpx.AbortWithError(c, err)
		return
	}

	httpx.JSON(c, http.StatusOK, response)
}

func (h *Handler) GetPresignedGET(c *gin.Context) {
	var request fctypes.PresignedGetRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpx.AbortWithError(c, httpx.BadRequest("预签名参数格式错误", nil))
		return
	}

	response, err := h.service.PresignGet(c.Request.Context(), request)
	if err != nil {
		httpx.AbortWithError(c, err)
		return
	}

	httpx.JSON(c, http.StatusOK, response)
}
