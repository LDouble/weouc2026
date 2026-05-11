package iam

import (
	"github.com/gin-gonic/gin"
	iamconfig "github.com/liangluo/weouc2026/services/api-server/internal/modules/iam/config"
	"github.com/liangluo/weouc2026/services/api-server/internal/modules/iam/repo"
	iamservice "github.com/liangluo/weouc2026/services/api-server/internal/modules/iam/service"
	"github.com/liangluo/weouc2026/services/api-server/internal/modules/iam/transport"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/audit"
	appconfig "github.com/liangluo/weouc2026/services/api-server/internal/platform/config"
	"github.com/liangluo/weouc2026/services/api-server/internal/providers/academic_provider"
	"github.com/liangluo/weouc2026/services/api-server/internal/providers/wechat_provider"
)

type Dependencies struct {
	UserRepository    repo.UserRepository
	SessionRepository repo.SessionRepository
	CaptchaRepository repo.CaptchaRepository
	WeChatProvider    wechat_provider.Provider
	AcademicProvider  academic_provider.Provider
	AuditRecorder     audit.Recorder
}

type Module struct {
	service *iamservice.Service
	handler *transport.Handler
}

func NewModule(appCfg appconfig.AppConfig, dependencies Dependencies) *Module {
	moduleConfig := iamconfig.New(appCfg.Auth.AccessTokenTTL)
	service := iamservice.New(
		moduleConfig,
		dependencies.UserRepository,
		dependencies.SessionRepository,
		dependencies.CaptchaRepository,
		dependencies.WeChatProvider,
		dependencies.AcademicProvider,
		dependencies.AuditRecorder,
	)
	handler := transport.NewHandler(service)

	return &Module{
		service: service,
		handler: handler,
	}
}

func (m *Module) RegisterRoutes(engine *gin.Engine) {
	transport.RegisterRoutes(engine, m.handler)
}

func (m *Module) AuthResolver() *iamservice.Service {
	return m.service
}
