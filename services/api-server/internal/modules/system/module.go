package system

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/liangluo/weouc2026/services/api-server/internal/modules/system/repo"
	"github.com/liangluo/weouc2026/services/api-server/internal/modules/system/service"
	"github.com/liangluo/weouc2026/services/api-server/internal/modules/system/transport"
	appconfig "github.com/liangluo/weouc2026/services/api-server/internal/platform/config"
)

type Module struct {
	handler *transport.Handler
	closers []io.Closer
}

func NewModule(appConfig appconfig.AppConfig) *Module {
	statusRepo, closers := repo.NewRuntimeStatusRepository(appConfig.Dependencies)
	systemService := service.New(statusRepo, appConfig)
	handler := transport.NewHandler(systemService)

	return &Module{
		handler: handler,
		closers: closers,
	}
}

func (m *Module) RegisterRoutes(engine *gin.Engine) {
	transport.RegisterRoutes(engine, m.handler)
}

func (m *Module) Closers() []io.Closer {
	if len(m.closers) == 0 {
		return nil
	}

	closers := make([]io.Closer, 0, len(m.closers))
	closers = append(closers, m.closers...)
	return closers
}
