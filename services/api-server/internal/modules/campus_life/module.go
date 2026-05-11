package campus_life

import (
	"github.com/gin-gonic/gin"
	clrepo "github.com/liangluo/weouc2026/services/api-server/internal/modules/campus_life/repo"
	clservice "github.com/liangluo/weouc2026/services/api-server/internal/modules/campus_life/service"
	"github.com/liangluo/weouc2026/services/api-server/internal/modules/campus_life/transport"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/audit"
	"github.com/liangluo/weouc2026/services/api-server/internal/providers/storage_provider"
)

type Dependencies struct {
	Repository      clrepo.Repository
	StorageProvider storage_provider.Provider
	AuditRecorder   audit.Recorder
}

type Module struct {
	handler *transport.Handler
}

func NewModule(dependencies Dependencies) *Module {
	service := clservice.New(dependencies.Repository, dependencies.StorageProvider, dependencies.AuditRecorder)
	handler := transport.NewHandler(service)
	return &Module{handler: handler}
}

func (m *Module) RegisterRoutes(engine *gin.Engine) {
	transport.RegisterRoutes(engine, m.handler)
}
