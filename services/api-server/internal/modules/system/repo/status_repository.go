package repo

import (
	"context"
	"time"

	"github.com/liangluo/weouc2026/services/api-server/internal/modules/system/types"
	appconfig "github.com/liangluo/weouc2026/services/api-server/internal/platform/config"
	"github.com/liangluo/weouc2026/services/api-server/internal/providers/storage_provider"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type StatusRepository interface {
	ReadinessSnapshot(ctx context.Context) types.ReadinessStatus
}

type DependencyProbe interface {
	Check(ctx context.Context) types.DependencyStatus
}

type RuntimeStatusRepository struct {
	probes []DependencyProbe
}

type StaticProbe struct {
	name     string
	status   string
	required bool
	detail   string
}

type MySQLProbe struct {
	config appconfig.MySQLConfig
	db     interface {
		PingContext(ctx context.Context) error
	}
}

type MongoProbe struct {
	config appconfig.MongoConfig
	client *mongo.Client
}

type RedisProbe struct {
	config appconfig.RedisConfig
	client *redis.Client
}

type ObjectStorageProbe struct {
	config   appconfig.COSConfig
	provider storage_provider.Provider
}

func NewRuntimeStatusRepository(probes ...DependencyProbe) *RuntimeStatusRepository {
	cloned := make([]DependencyProbe, 0, len(probes))
	cloned = append(cloned, probes...)
	return &RuntimeStatusRepository{probes: cloned}
}

func NewStaticProbe(name, status string, required bool, detail string) *StaticProbe {
	return &StaticProbe{
		name:     name,
		status:   status,
		required: required,
		detail:   detail,
	}
}

func NewMySQLProbe(cfg appconfig.MySQLConfig, db interface {
	PingContext(ctx context.Context) error
}) *MySQLProbe {
	return &MySQLProbe{config: cfg, db: db}
}

func NewMongoProbe(cfg appconfig.MongoConfig, client *mongo.Client) *MongoProbe {
	return &MongoProbe{config: cfg, client: client}
}

func NewRedisProbe(cfg appconfig.RedisConfig, client *redis.Client) *RedisProbe {
	return &RedisProbe{config: cfg, client: client}
}

func NewObjectStorageProbe(cfg appconfig.COSConfig, provider storage_provider.Provider) *ObjectStorageProbe {
	return &ObjectStorageProbe{config: cfg, provider: provider}
}

func (r *RuntimeStatusRepository) ReadinessSnapshot(ctx context.Context) types.ReadinessStatus {
	dependencies := make([]types.DependencyStatus, 0, len(r.probes))
	ready := true
	for _, probe := range r.probes {
		dependency := probe.Check(ctx)
		if dependency.Required && dependency.Status != "ready" {
			ready = false
		}
		dependencies = append(dependencies, dependency)
	}

	status := "ready"
	if !ready {
		status = "not_ready"
	}

	return types.ReadinessStatus{
		Status:       status,
		Dependencies: dependencies,
		Timestamp:    time.Now().UTC(),
	}
}

func (p *StaticProbe) Check(context.Context) types.DependencyStatus {
	return types.DependencyStatus{
		Name:     p.name,
		Status:   p.status,
		Required: p.required,
		Detail:   p.detail,
	}
}

func (p *MySQLProbe) Check(ctx context.Context) types.DependencyStatus {
	if !p.config.Enabled {
		return types.DependencyStatus{
			Name:     "mysql",
			Status:   "skipped",
			Required: false,
			Detail:   "未启用 MySQL 健康探测",
		}
	}
	if p.db == nil {
		return types.DependencyStatus{
			Name:     "mysql",
			Status:   "not_ready",
			Required: true,
			Detail:   "MySQL 连接器初始化失败",
		}
	}

	checkCtx, cancel := context.WithTimeout(ctx, p.config.HealthCheckTimeout)
	defer cancel()

	if err := p.db.PingContext(checkCtx); err != nil {
		return types.DependencyStatus{
			Name:     "mysql",
			Status:   "not_ready",
			Required: true,
			Detail:   err.Error(),
		}
	}

	return types.DependencyStatus{
		Name:     "mysql",
		Status:   "ready",
		Required: true,
		Detail:   "连接正常",
	}
}

func (p *MongoProbe) Check(ctx context.Context) types.DependencyStatus {
	if !p.config.Enabled {
		return types.DependencyStatus{
			Name:     "mongo",
			Status:   "skipped",
			Required: false,
			Detail:   "未启用 MongoDB 健康探测",
		}
	}
	if p.client == nil {
		return types.DependencyStatus{
			Name:     "mongo",
			Status:   "not_ready",
			Required: true,
			Detail:   "MongoDB 连接器初始化失败",
		}
	}

	checkCtx, cancel := context.WithTimeout(ctx, p.config.HealthCheckTimeout)
	defer cancel()

	if err := p.client.Ping(checkCtx, readpref.Primary()); err != nil {
		return types.DependencyStatus{
			Name:     "mongo",
			Status:   "not_ready",
			Required: true,
			Detail:   err.Error(),
		}
	}

	return types.DependencyStatus{
		Name:     "mongo",
		Status:   "ready",
		Required: true,
		Detail:   "连接正常",
	}
}

func (p *RedisProbe) Check(ctx context.Context) types.DependencyStatus {
	if !p.config.Enabled {
		return types.DependencyStatus{
			Name:     "redis",
			Status:   "skipped",
			Required: false,
			Detail:   "未启用 Redis 健康探测",
		}
	}
	if p.client == nil {
		return types.DependencyStatus{
			Name:     "redis",
			Status:   "not_ready",
			Required: true,
			Detail:   "Redis 客户端初始化失败",
		}
	}

	checkCtx, cancel := context.WithTimeout(ctx, p.config.HealthCheckTimeout)
	defer cancel()

	if err := p.client.Ping(checkCtx).Err(); err != nil {
		return types.DependencyStatus{
			Name:     "redis",
			Status:   "not_ready",
			Required: true,
			Detail:   err.Error(),
		}
	}

	return types.DependencyStatus{
		Name:     "redis",
		Status:   "ready",
		Required: true,
		Detail:   "连接正常",
	}
}

func (p *ObjectStorageProbe) Check(ctx context.Context) types.DependencyStatus {
	if !p.config.Enabled {
		return types.DependencyStatus{
			Name:     "object_storage",
			Status:   "skipped",
			Required: false,
			Detail:   "未启用对象存储健康探测",
		}
	}
	if p.provider == nil {
		return types.DependencyStatus{
			Name:     "object_storage",
			Status:   "not_ready",
			Required: true,
			Detail:   "对象存储 Provider 初始化失败",
		}
	}

	checkCtx, cancel := context.WithTimeout(ctx, p.config.HealthCheckTimeout)
	defer cancel()

	if err := p.provider.Check(checkCtx); err != nil {
		return types.DependencyStatus{
			Name:     "object_storage",
			Status:   "not_ready",
			Required: true,
			Detail:   err.Error(),
		}
	}

	return types.DependencyStatus{
		Name:     "object_storage",
		Status:   "ready",
		Required: true,
		Detail:   "连接正常",
	}
}
