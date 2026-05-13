package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"io"

	mysqlDriver "github.com/go-sql-driver/mysql"
	appconfig "github.com/liangluo/weouc2026/services/api-server/internal/platform/config"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Clients struct {
	MySQL       *gorm.DB
	MySQLSQLDB  *sql.DB
	Mongo       *mongo.Database
	MongoClient *mongo.Client
	Redis       *redis.Client
	closers     []io.Closer
}

func Open(cfg appconfig.AppConfig) (*Clients, error) {
	clients := &Clients{}

	if cfg.Dependencies.MySQL.Enabled {
		gormDB, sqlDB, err := openMySQL(cfg.Dependencies.MySQL)
		if err != nil {
			return nil, fmt.Errorf("open mysql failed: %w", err)
		}
		clients.MySQL = gormDB
		clients.MySQLSQLDB = sqlDB
		clients.closers = append(clients.closers, sqlDB)
	}

	if cfg.Dependencies.Mongo.Enabled {
		client, database, err := openMongo(cfg.Dependencies.Mongo)
		if err != nil {
			return nil, fmt.Errorf("open mongo failed: %w", err)
		}
		clients.Mongo = database
		clients.MongoClient = client
		clients.closers = append(clients.closers, mongoCloser{client: client})
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

type mongoCloser struct {
	client *mongo.Client
}

func (c mongoCloser) Close() error {
	if c.client == nil {
		return nil
	}

	return c.client.Disconnect(context.Background())
}

func openMySQL(cfg appconfig.MySQLConfig) (*gorm.DB, *sql.DB, error) {
	driverConfig, err := mysqlDriver.ParseDSN(cfg.DSN())
	if err != nil {
		return nil, nil, fmt.Errorf("parse mysql dsn failed: %w", err)
	}

	timeout := cfg.HealthCheckTimeout
	driverConfig.Timeout = timeout
	driverConfig.ReadTimeout = timeout
	driverConfig.WriteTimeout = timeout

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       driverConfig.FormatDSN(),
		DefaultStringSize:         256,
		DisableDatetimePrecision:  false,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		return nil, nil, err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("resolve mysql sql db failed: %w", err)
	}

	return gormDB, sqlDB, nil
}

func openMongo(cfg appconfig.MongoConfig) (*mongo.Client, *mongo.Database, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, nil, err
	}

	return client, client.Database(cfg.Database), nil
}

func (c *Clients) EnsureRuntimeBackendsReady(ctx context.Context, cfg appconfig.AppConfig) error {
	requiresMySQL := cfg.Persistence.IAMBackendOrDefault() == "mysql_redis"
	requiresMongo := cfg.Persistence.CampusLifeBackendOrDefault() == "mongo" ||
		cfg.Persistence.PortalBackendOrDefault() == "mongo" ||
		cfg.Persistence.NotificationBackendOrDefault() == "mongo" ||
		cfg.Persistence.AnalyticsBackendOrDefault() == "mongo"
	requiresRedis := cfg.Persistence.IAMBackendOrDefault() == "mysql_redis"
	if requiresMySQL {
		if c.MySQLSQLDB == nil {
			return fmt.Errorf("mysql client is not configured")
		}
		mysqlCtx, cancel := context.WithTimeout(ctx, cfg.Dependencies.MySQL.HealthCheckTimeout)
		defer cancel()
		if err := c.MySQLSQLDB.PingContext(mysqlCtx); err != nil {
			return fmt.Errorf("mysql ping failed: %w", err)
		}
	}

	if requiresMongo {
		if c.MongoClient == nil {
			return fmt.Errorf("mongo client is not configured")
		}
		mongoCtx, cancel := context.WithTimeout(ctx, cfg.Dependencies.Mongo.HealthCheckTimeout)
		defer cancel()
		if err := c.MongoClient.Ping(mongoCtx, readpref.Primary()); err != nil {
			return fmt.Errorf("mongo ping failed: %w", err)
		}
	}

	if requiresRedis {
		if c.Redis == nil {
			return fmt.Errorf("redis client is not configured")
		}
		redisCtx, cancel := context.WithTimeout(ctx, cfg.Dependencies.Redis.HealthCheckTimeout)
		defer cancel()
		if err := c.Redis.Ping(redisCtx).Err(); err != nil {
			return fmt.Errorf("redis ping failed: %w", err)
		}
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
