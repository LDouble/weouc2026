package repo

import (
	"context"
	"database/sql"
	"io"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/liangluo/weouc2026/services/api-server/internal/modules/system/types"
	appconfig "github.com/liangluo/weouc2026/services/api-server/internal/platform/config"
	"github.com/redis/go-redis/v9"
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

type PostgresProbe struct {
	config appconfig.PostgresConfig
	db     *sql.DB
}

type RedisProbe struct {
	config appconfig.RedisConfig
	client *redis.Client
}

func NewRuntimeStatusRepository(cfg appconfig.DependenciesConfig) (*RuntimeStatusRepository, []io.Closer) {
	probes := []DependencyProbe{
		NewPostgresProbe(cfg.Postgres),
		NewRedisProbe(cfg.Redis),
		NewStaticProbe("object_storage", "skipped", false, "当前阶段未接入对象存储健康探测"),
	}
	closers := make([]io.Closer, 0, len(probes))
	for _, probe := range probes {
		closer, ok := probe.(io.Closer)
		if !ok {
			continue
		}
		closers = append(closers, closer)
	}

	return &RuntimeStatusRepository{probes: probes}, closers
}

func NewStaticProbe(name, status string, required bool, detail string) *StaticProbe {
	return &StaticProbe{
		name:     name,
		status:   status,
		required: required,
		detail:   detail,
	}
}

func NewPostgresProbe(cfg appconfig.PostgresConfig) *PostgresProbe {
	if !cfg.Enabled {
		return &PostgresProbe{config: cfg}
	}

	db, err := sql.Open("pgx", cfg.DSN())
	if err != nil {
		return &PostgresProbe{config: cfg}
	}

	return &PostgresProbe{
		config: cfg,
		db:     db,
	}
}

func NewRedisProbe(cfg appconfig.RedisConfig) *RedisProbe {
	if !cfg.Enabled {
		return &RedisProbe{config: cfg}
	}

	return &RedisProbe{
		config: cfg,
		client: redis.NewClient(&redis.Options{
			Addr:     cfg.Address(),
			Username: cfg.Username,
			Password: cfg.Password,
			DB:       cfg.Database,
		}),
	}
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

func (p *PostgresProbe) Check(ctx context.Context) types.DependencyStatus {
	if !p.config.Enabled {
		return types.DependencyStatus{
			Name:     "postgres",
			Status:   "skipped",
			Required: false,
			Detail:   "未启用 PostgreSQL 健康探测",
		}
	}
	if p.db == nil {
		return types.DependencyStatus{
			Name:     "postgres",
			Status:   "not_ready",
			Required: true,
			Detail:   "PostgreSQL 连接器初始化失败",
		}
	}

	checkCtx, cancel := context.WithTimeout(ctx, p.config.HealthCheckTimeout)
	defer cancel()

	if err := p.db.PingContext(checkCtx); err != nil {
		return types.DependencyStatus{
			Name:     "postgres",
			Status:   "not_ready",
			Required: true,
			Detail:   err.Error(),
		}
	}

	return types.DependencyStatus{
		Name:     "postgres",
		Status:   "ready",
		Required: true,
		Detail:   "连接正常",
	}
}

func (p *PostgresProbe) Close() error {
	if p.db == nil {
		return nil
	}

	return p.db.Close()
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

func (p *RedisProbe) Close() error {
	if p.client == nil {
		return nil
	}

	return p.client.Close()
}
