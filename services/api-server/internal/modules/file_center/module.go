package file_center

import (
	"github.com/gin-gonic/gin"
	fcconfig "github.com/liangluo/weouc2026/services/api-server/internal/modules/file_center/config"
	fcservice "github.com/liangluo/weouc2026/services/api-server/internal/modules/file_center/service"
	"github.com/liangluo/weouc2026/services/api-server/internal/modules/file_center/transport"
	"github.com/liangluo/weouc2026/services/api-server/internal/providers/storage_provider"
)

type Dependencies struct {
	StorageProvider storage_provider.Provider
}

type Module struct {
	handler *transport.Handler
}

func NewModule(dependencies Dependencies) *Module {
	service := fcservice.New(fcconfig.DefaultModuleConfig(), dependencies.StorageProvider)
	handler := transport.NewHandler(service)
	return &Module{handler: handler}
}

func (m *Module) RegisterRoutes(engine *gin.Engine) {
	transport.RegisterRoutes(engine, m.handler)
}
