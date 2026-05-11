package notification

import (
	"github.com/gin-gonic/gin"
	notificationrepo "github.com/liangluo/weouc2026/services/api-server/internal/modules/notification/repo"
	notificationservice "github.com/liangluo/weouc2026/services/api-server/internal/modules/notification/service"
	"github.com/liangluo/weouc2026/services/api-server/internal/modules/notification/transport"
)

type Dependencies struct {
	Repository notificationrepo.Repository
}

type Module struct {
	handler *transport.Handler
}

func NewModule(dependencies Dependencies) *Module {
	repository := dependencies.Repository
	if repository == nil {
		repository = notificationrepo.NewInMemoryRepository()
	}
	service := notificationservice.New(repository)
	handler := transport.NewHandler(service)
	return &Module{handler: handler}
}

func (m *Module) RegisterRoutes(engine *gin.Engine) {
	transport.RegisterRoutes(engine, m.handler)
}
