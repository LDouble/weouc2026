package academic

import (
	"github.com/gin-gonic/gin"
	academicrepo "github.com/liangluo/weouc2026/services/api-server/internal/modules/academic/repo"
	academicservice "github.com/liangluo/weouc2026/services/api-server/internal/modules/academic/service"
	"github.com/liangluo/weouc2026/services/api-server/internal/modules/academic/transport"
	iamrepo "github.com/liangluo/weouc2026/services/api-server/internal/modules/iam/repo"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/audit"
	"github.com/liangluo/weouc2026/services/api-server/internal/providers/academic_provider"
)

type Dependencies struct {
	UserRepository   iamrepo.UserRepository
	AcademicProvider academic_provider.Provider
	AuditRecorder    audit.Recorder
}

type Module struct {
	handler *transport.Handler
}

func NewModule(dependencies Dependencies) *Module {
	repository := academicrepo.NewProviderRepository(dependencies.AcademicProvider)
	service := academicservice.New(repository, dependencies.UserRepository, dependencies.AuditRecorder)
	handler := transport.NewHandler(service)

	return &Module{handler: handler}
}

func (m *Module) RegisterRoutes(engine *gin.Engine) {
	transport.RegisterRoutes(engine, m.handler)
}
