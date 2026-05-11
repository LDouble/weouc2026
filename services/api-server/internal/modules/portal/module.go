package portal

import (
	"github.com/gin-gonic/gin"
	portalrepo "github.com/liangluo/weouc2026/services/api-server/internal/modules/portal/repo"
	portalservice "github.com/liangluo/weouc2026/services/api-server/internal/modules/portal/service"
	"github.com/liangluo/weouc2026/services/api-server/internal/modules/portal/transport"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/audit"
)

type Dependencies struct {
	Repository    portalrepo.Repository
	AuditRecorder audit.Recorder
}

type Module struct {
	handler *transport.Handler
}

func NewModule(dependencies Dependencies) *Module {
	repository := dependencies.Repository
	if repository == nil {
		repository = portalrepo.NewInMemoryRepository()
	}
	service := portalservice.New(repository, dependencies.AuditRecorder)
	handler := transport.NewHandler(service)
	return &Module{handler: handler}
}

func (m *Module) RegisterRoutes(engine *gin.Engine) {
	transport.RegisterRoutes(engine, m.handler)
}
