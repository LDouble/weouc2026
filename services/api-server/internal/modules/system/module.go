package system

import (
	"github.com/gin-gonic/gin"
	"github.com/liangluo/weouc2026/services/api-server/internal/modules/system/repo"
	"github.com/liangluo/weouc2026/services/api-server/internal/modules/system/service"
	"github.com/liangluo/weouc2026/services/api-server/internal/modules/system/transport"
	appconfig "github.com/liangluo/weouc2026/services/api-server/internal/platform/config"
)

type Dependencies struct {
	StatusRepository repo.StatusRepository
}

type Module struct {
	handler *transport.Handler
}

func NewModule(appConfig appconfig.AppConfig, dependencies Dependencies) *Module {
	statusRepo := dependencies.StatusRepository
	if statusRepo == nil {
		statusRepo = repo.NewRuntimeStatusRepository(
			repo.NewStaticProbe("postgres", "skipped", false, "未启用 PostgreSQL 健康探测"),
			repo.NewStaticProbe("redis", "skipped", false, "未启用 Redis 健康探测"),
			repo.NewStaticProbe("object_storage", "skipped", false, "当前阶段未接入对象存储健康探测"),
		)
	}
	systemService := service.New(statusRepo, appConfig)
	handler := transport.NewHandler(systemService)

	return &Module{
		handler: handler,
	}
}

func (m *Module) RegisterRoutes(engine *gin.Engine) {
	transport.RegisterRoutes(engine, m.handler)
}
