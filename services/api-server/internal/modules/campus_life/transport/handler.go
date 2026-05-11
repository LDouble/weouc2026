package transport

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	clservice "github.com/liangluo/weouc2026/services/api-server/internal/modules/campus_life/service"
	cltypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/campus_life/types"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/auth"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/httpx"
)

type Handler struct {
	service *clservice.Service
}

func NewHandler(service *clservice.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListFeed(c *gin.Context) {
	query := cltypes.FeedQuery{
		Pagination: paginationFromContext(c),
		FeedTypes:  queryArray(c, "feed_types"),
		Keyword:    strings.TrimSpace(c.Query("keyword")),
		UserRole:   strings.TrimSpace(c.Query("user_role")),
	}
	response, err := h.service.ListFeed(c.Request.Context(), auth.PrincipalFromContext(c), query)
	if err != nil {
		httpx.AbortWithError(c, err)
		return
	}
	httpx.JSON(c, http.StatusOK, response)
}

func (h *Handler) ListMarket(c *gin.Context) {
	query := cltypes.MarketQuery{
		Pagination: paginationFromContext(c),
		Category:   strings.TrimSpace(c.Query("category")),
		Keyword:    strings.TrimSpace(c.Query("keyword")),
	}
	response, err := h.service.ListMarket(c.Request.Context(), auth.PrincipalFromContext(c), query)
	if err != nil {
		httpx.AbortWithError(c, err)
		return
	}
	httpx.JSON(c, http.StatusOK, response)
}

func (h *Handler) GetMarketDetail(c *gin.Context) {
	response, err := h.service.GetMarketDetail(c.Request.Context(), auth.PrincipalFromContext(c), c.Param("id"))
	if err != nil {
		httpx.AbortWithError(c, err)
		return
	}
	httpx.JSON(c, http.StatusOK, response)
}

func (h *Handler) PublishMarket(c *gin.Context) {
	var request cltypes.MarketPublishRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpx.AbortWithError(c, httpx.BadRequest("二手发布参数格式错误", nil))
		return
	}
	response, err := h.service.PublishMarket(c.Request.Context(), auth.PrincipalFromContext(c), request)
	if err != nil {
		httpx.AbortWithError(c, err)
		return
	}
	httpx.JSON(c, http.StatusOK, response)
}

func (h *Handler) FavoriteMarket(c *gin.Context) {
	var request cltypes.FavoriteMarketRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpx.AbortWithError(c, httpx.BadRequest("收藏参数格式错误", nil))
		return
	}
	if err := h.service.FavoriteMarket(c.Request.Context(), auth.PrincipalFromContext(c), request); err != nil {
		httpx.AbortWithError(c, err)
		return
	}
	httpx.JSON(c, http.StatusOK, gin.H{"success": true})
}

func (h *Handler) ListErrands(c *gin.Context) {
	query := cltypes.ErrandQuery{
		Pagination: paginationFromContext(c),
		Category:   strings.TrimSpace(c.Query("category")),
		Keyword:    strings.TrimSpace(c.Query("keyword")),
		UserRole:   strings.TrimSpace(c.Query("user_role")),
	}
	response, err := h.service.ListErrands(c.Request.Context(), auth.PrincipalFromContext(c), query)
	if err != nil {
		httpx.AbortWithError(c, err)
		return
	}
	httpx.JSON(c, http.StatusOK, response)
}

func (h *Handler) GetErrandDetail(c *gin.Context) {
	response, err := h.service.GetErrandDetail(c.Request.Context(), auth.PrincipalFromContext(c), c.Param("id"))
	if err != nil {
		httpx.AbortWithError(c, err)
		return
	}
	httpx.JSON(c, http.StatusOK, response)
}

func (h *Handler) PublishErrand(c *gin.Context) {
	var request cltypes.ErrandPublishRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpx.AbortWithError(c, httpx.BadRequest("跑腿发布参数格式错误", nil))
		return
	}
	response, err := h.service.PublishErrand(c.Request.Context(), auth.PrincipalFromContext(c), request)
	if err != nil {
		httpx.AbortWithError(c, err)
		return
	}
	httpx.JSON(c, http.StatusOK, response)
}

func (h *Handler) AcceptErrand(c *gin.Context) {
	var request cltypes.ErrandActionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpx.AbortWithError(c, httpx.BadRequest("接单参数格式错误", nil))
		return
	}
	if err := h.service.AcceptErrand(c.Request.Context(), auth.PrincipalFromContext(c), request.TaskID); err != nil {
		httpx.AbortWithError(c, err)
		return
	}
	httpx.JSON(c, http.StatusOK, gin.H{"success": true})
}

func (h *Handler) CancelErrandPublish(c *gin.Context) {
	var request cltypes.ErrandActionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpx.AbortWithError(c, httpx.BadRequest("取消发布参数格式错误", nil))
		return
	}
	if err := h.service.CancelErrandPublish(c.Request.Context(), auth.PrincipalFromContext(c), request.TaskID); err != nil {
		httpx.AbortWithError(c, err)
		return
	}
	httpx.JSON(c, http.StatusOK, gin.H{"success": true})
}

func (h *Handler) CancelErrandAccept(c *gin.Context) {
	var request cltypes.ErrandActionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpx.AbortWithError(c, httpx.BadRequest("取消接单参数格式错误", nil))
		return
	}
	if err := h.service.CancelErrandAccept(c.Request.Context(), auth.PrincipalFromContext(c), request.TaskID); err != nil {
		httpx.AbortWithError(c, err)
		return
	}
	httpx.JSON(c, http.StatusOK, gin.H{"success": true})
}

func (h *Handler) ListResources(c *gin.Context) {
	query := cltypes.ResourceQuery{
		Pagination: paginationFromContext(c),
		Category:   strings.TrimSpace(c.Query("category")),
		Keyword:    strings.TrimSpace(c.Query("keyword")),
	}
	response, err := h.service.ListResources(c.Request.Context(), query)
	if err != nil {
		httpx.AbortWithError(c, err)
		return
	}
	httpx.JSON(c, http.StatusOK, response)
}

func (h *Handler) GetResourceDetail(c *gin.Context) {
	response, err := h.service.GetResourceDetail(c.Request.Context(), c.Param("id"))
	if err != nil {
		httpx.AbortWithError(c, err)
		return
	}
	httpx.JSON(c, http.StatusOK, response)
}

func (h *Handler) PublishResource(c *gin.Context) {
	var request cltypes.ResourcePublishRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpx.AbortWithError(c, httpx.BadRequest("资料发布参数格式错误", nil))
		return
	}
	response, err := h.service.PublishResource(c.Request.Context(), auth.PrincipalFromContext(c), request)
	if err != nil {
		httpx.AbortWithError(c, err)
		return
	}
	httpx.JSON(c, http.StatusOK, response)
}

func (h *Handler) ListLostFound(c *gin.Context) {
	query := cltypes.LostFoundQuery{
		Pagination: paginationFromContext(c),
		Category:   strings.TrimSpace(c.Query("category")),
		Keyword:    strings.TrimSpace(c.Query("keyword")),
		Type:       strings.TrimSpace(c.Query("type")),
	}
	response, err := h.service.ListLostFound(c.Request.Context(), auth.PrincipalFromContext(c), query)
	if err != nil {
		httpx.AbortWithError(c, err)
		return
	}
	httpx.JSON(c, http.StatusOK, response)
}

func (h *Handler) GetLostFoundDetail(c *gin.Context) {
	response, err := h.service.GetLostFoundDetail(c.Request.Context(), auth.PrincipalFromContext(c), c.Param("id"))
	if err != nil {
		httpx.AbortWithError(c, err)
		return
	}
	httpx.JSON(c, http.StatusOK, response)
}

func (h *Handler) ListCarpools(c *gin.Context) {
	query := cltypes.CarpoolQuery{
		Pagination: paginationFromContext(c),
		Category:   strings.TrimSpace(c.Query("category")),
		Keyword:    strings.TrimSpace(c.Query("keyword")),
	}
	response, err := h.service.ListCarpools(c.Request.Context(), auth.PrincipalFromContext(c), query)
	if err != nil {
		httpx.AbortWithError(c, err)
		return
	}
	httpx.JSON(c, http.StatusOK, response)
}

func (h *Handler) GetCarpoolDetail(c *gin.Context) {
	response, err := h.service.GetCarpoolDetail(c.Request.Context(), auth.PrincipalFromContext(c), c.Param("id"))
	if err != nil {
		httpx.AbortWithError(c, err)
		return
	}
	httpx.JSON(c, http.StatusOK, response)
}

func (h *Handler) PublishCarpool(c *gin.Context) {
	var request cltypes.CarpoolPublishRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpx.AbortWithError(c, httpx.BadRequest("拼车发布参数格式错误", nil))
		return
	}
	response, err := h.service.PublishCarpool(c.Request.Context(), auth.PrincipalFromContext(c), request)
	if err != nil {
		httpx.AbortWithError(c, err)
		return
	}
	httpx.JSON(c, http.StatusOK, response)
}

func (h *Handler) PublishLostFound(c *gin.Context) {
	var request cltypes.LostFoundPublishRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpx.AbortWithError(c, httpx.BadRequest("失物招领发布参数格式错误", nil))
		return
	}
	response, err := h.service.PublishLostFound(c.Request.Context(), auth.PrincipalFromContext(c), request)
	if err != nil {
		httpx.AbortWithError(c, err)
		return
	}
	httpx.JSON(c, http.StatusOK, response)
}
