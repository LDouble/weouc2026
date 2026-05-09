package repo

import (
	"context"
	"time"

	moduleconfig "github.com/liangluo/weouc2026/services/api-server/internal/modules/system/config"
	"github.com/liangluo/weouc2026/services/api-server/internal/modules/system/types"
)

type StatusRepository interface {
	ReadinessSnapshot(ctx context.Context) types.ReadinessStatus
}

type StaticStatusRepository struct {
	dependencyNames []string
}

func NewStaticStatusRepository(cfg moduleconfig.ModuleConfig) *StaticStatusRepository {
	names := append([]string(nil), cfg.DependencyNames...)
	return &StaticStatusRepository{dependencyNames: names}
}

func (r *StaticStatusRepository) ReadinessSnapshot(context.Context) types.ReadinessStatus {
	dependencies := make([]types.DependencyStatus, 0, len(r.dependencyNames))
	for _, name := range r.dependencyNames {
		dependencies = append(dependencies, types.DependencyStatus{
			Name:     name,
			Status:   "skipped",
			Required: false,
			Detail:   "基础骨架阶段未接入真实依赖",
		})
	}

	return types.ReadinessStatus{
		Status:       "ready",
		Dependencies: dependencies,
		Timestamp:    time.Now().UTC(),
	}
}
