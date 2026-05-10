package transport

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	cltypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/campus_life/types"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/auth"
)

func RegisterRoutes(engine *gin.Engine, handler *Handler) {
	api := engine.Group("/api")

	api.GET("/feed/list", handler.ListFeed)
	api.GET("/market/list", handler.ListMarket)
	api.GET("/market/detail/:id", handler.GetMarketDetail)
	api.GET("/errand/list", handler.ListErrands)
	api.GET("/errand/detail/:id", handler.GetErrandDetail)
	api.GET("/resource/list", handler.ListResources)
	api.GET("/resource/detail/:id", handler.GetResourceDetail)
	api.GET("/lostFound/list", handler.ListLostFound)
	api.GET("/lostFound/detail/:id", handler.GetLostFoundDetail)

	protected := api.Group("")
	protected.Use(auth.RequireAuthenticated())
	protected.POST("/market/publish", handler.PublishMarket)
	protected.POST("/market/favorite", handler.FavoriteMarket)
	protected.POST("/errand/publish", handler.PublishErrand)
	protected.POST("/errand/accept", handler.AcceptErrand)
	protected.POST("/errand/cancel-publish", handler.CancelErrandPublish)
	protected.POST("/errand/cancel-accept", handler.CancelErrandAccept)
	protected.POST("/resource/publish", handler.PublishResource)
	protected.POST("/lostFound/publish", handler.PublishLostFound)
}

func paginationFromContext(c *gin.Context) cltypes.Pagination {
	return cltypes.Pagination{
		Page:     parsePositiveInt(c.Query("page"), 1),
		PageSize: parsePositiveInt(c.Query("pageSize"), 20),
	}
}

func queryArray(c *gin.Context, key string) []string {
	values := c.QueryArray(key)
	if len(values) > 0 {
		return values
	}

	raw := strings.TrimSpace(c.Query(key))
	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		value := strings.Trim(strings.TrimSpace(part), "[]\"")
		if value != "" {
			result = append(result, value)
		}
	}
	return result
}

func parsePositiveInt(raw string, defaultValue int) int {
	value, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || value <= 0 {
		return defaultValue
	}
	return value
}
