package service

import (
	"context"
	"time"

	"github.com/liangluo/weouc2026/services/api-server/internal/modules/system/repo"
	"github.com/liangluo/weouc2026/services/api-server/internal/modules/system/types"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/auth"
	appconfig "github.com/liangluo/weouc2026/services/api-server/internal/platform/config"
)

type Service struct {
	statusRepo repo.StatusRepository
	appConfig  appconfig.AppConfig
}

func New(statusRepo repo.StatusRepository, appConfig appconfig.AppConfig) *Service {
	return &Service{
		statusRepo: statusRepo,
		appConfig:  appConfig,
	}
}

func (s *Service) Health(context.Context) types.HealthStatus {
	return types.HealthStatus{
		Status:    "ok",
		Service:   s.appConfig.Service.Name,
		Version:   s.appConfig.Service.Version,
		Timestamp: time.Now().UTC(),
	}
}

func (s *Service) Ready(ctx context.Context) types.ReadinessStatus {
	return s.statusRepo.ReadinessSnapshot(ctx)
}

func (s *Service) Profile(_ context.Context, principal auth.Principal) types.Profile {
	return types.Profile{
		Service: types.ProfileService{
			Name:        s.appConfig.Service.Name,
			Environment: s.appConfig.Service.Environment,
			Version:     s.appConfig.Service.Version,
		},
		Auth: types.ProfileAuth{
			Authenticated: principal.Authenticated,
			UserID:        principal.UserID,
			Roles:         principal.Roles,
			Permissions:   principal.Permissions,
			AcademicBound: principal.AcademicBound,
		},
	}
}
