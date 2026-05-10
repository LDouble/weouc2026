package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"io"

	_ "github.com/jackc/pgx/v5/stdlib"
	appconfig "github.com/liangluo/weouc2026/services/api-server/internal/platform/config"
	"github.com/redis/go-redis/v9"
)

type Clients struct {
	Postgres *sql.DB
	Redis    *redis.Client
	closers  []io.Closer
}

func Open(cfg appconfig.AppConfig) (*Clients, error) {
	clients := &Clients{}

	if cfg.Dependencies.Postgres.Enabled {
		db, err := sql.Open("pgx", cfg.Dependencies.Postgres.DSN())
		if err != nil {
			return nil, fmt.Errorf("open postgres failed: %w", err)
		}
		clients.Postgres = db
		clients.closers = append(clients.closers, db)
	}

	if cfg.Dependencies.Redis.Enabled {
		client := redis.NewClient(&redis.Options{
			Addr:     cfg.Dependencies.Redis.Address(),
			Username: cfg.Dependencies.Redis.Username,
			Password: cfg.Dependencies.Redis.Password,
			DB:       cfg.Dependencies.Redis.Database,
		})
		clients.Redis = client
		clients.closers = append(clients.closers, client)
	}

	return clients, nil
}

func (c *Clients) EnsureIAMBackendReady(ctx context.Context, cfg appconfig.AppConfig) error {
	if cfg.Persistence.IAMBackendOrDefault() != "postgres_redis" {
		return nil
	}
	if c.Postgres == nil {
		return fmt.Errorf("postgres client is not configured")
	}
	postgresCtx, cancel := context.WithTimeout(ctx, cfg.Dependencies.Postgres.HealthCheckTimeout)
	defer cancel()
	if err := c.Postgres.PingContext(postgresCtx); err != nil {
		return fmt.Errorf("postgres ping failed: %w", err)
	}
	if c.Redis == nil {
		return fmt.Errorf("redis client is not configured")
	}
	redisCtx, cancel := context.WithTimeout(ctx, cfg.Dependencies.Redis.HealthCheckTimeout)
	defer cancel()
	if err := c.Redis.Ping(redisCtx).Err(); err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}

	return nil
}

func (c *Clients) Closers() []io.Closer {
	if len(c.closers) == 0 {
		return nil
	}

	closers := make([]io.Closer, 0, len(c.closers))
	closers = append(closers, c.closers...)
	return closers
}
