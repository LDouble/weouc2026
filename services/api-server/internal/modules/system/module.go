package system

import (
	"github.com/gin-gonic/gin"
	moduleconfig "github.com/liangluo/weouc2026/services/api-server/internal/modules/system/config"
	"github.com/liangluo/weouc2026/services/api-server/internal/modules/system/repo"
	"github.com/liangluo/weouc2026/services/api-server/internal/modules/system/service"
	"github.com/liangluo/weouc2026/services/api-server/internal/modules/system/transport"
	appconfig "github.com/liangluo/weouc2026/services/api-server/internal/platform/config"
)

type Module struct {
	handler *transport.Handler
}

func NewModule(appConfig appconfig.AppConfig) *Module {
	moduleCfg := moduleconfig.DefaultModuleConfig()
	statusRepo := repo.NewStaticStatusRepository(moduleCfg)
	systemService := service.New(statusRepo, appConfig)
	handler := transport.NewHandler(systemService)

	return &Module{handler: handler}
}

func (m *Module) RegisterRoutes(engine *gin.Engine) {
	transport.RegisterRoutes(engine, m.handler)
}
