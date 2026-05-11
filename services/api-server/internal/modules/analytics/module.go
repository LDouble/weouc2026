package analytics

import (
	"github.com/gin-gonic/gin"
	analyticsrepo "github.com/liangluo/weouc2026/services/api-server/internal/modules/analytics/repo"
	analyticsservice "github.com/liangluo/weouc2026/services/api-server/internal/modules/analytics/service"
	"github.com/liangluo/weouc2026/services/api-server/internal/modules/analytics/transport"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/audit"
)

type Dependencies struct {
	AuditStore audit.Repository
}

type Module struct {
	handler *transport.Handler
}

func NewModule(dependencies Dependencies) *Module {
	repository := analyticsrepo.NewAuditRepository(dependencies.AuditStore)
	service := analyticsservice.New(repository)
	handler := transport.NewHandler(service)
	return &Module{handler: handler}
}

func (m *Module) RegisterRoutes(engine *gin.Engine) {
	transport.RegisterRoutes(engine, m.handler)
}
