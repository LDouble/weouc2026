package notification

import (
	"github.com/gin-gonic/gin"
	notificationrepo "github.com/liangluo/weouc2026/services/api-server/internal/modules/notification/repo"
	notificationservice "github.com/liangluo/weouc2026/services/api-server/internal/modules/notification/service"
	"github.com/liangluo/weouc2026/services/api-server/internal/modules/notification/transport"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/audit"
)

type Dependencies struct {
	Repository    notificationrepo.Repository
	AuditRecorder audit.Recorder
}

type Module struct {
	handler *transport.Handler
}

func NewModule(dependencies Dependencies) *Module {
	service := notificationservice.New(dependencies.Repository, dependencies.AuditRecorder)
	handler := transport.NewHandler(service)
	return &Module{handler: handler}
}

func (m *Module) RegisterRoutes(engine *gin.Engine) {
	transport.RegisterRoutes(engine, m.handler)
}
