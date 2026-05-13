package migrate

import (
	"context"
	"fmt"

	iamrepo "github.com/liangluo/weouc2026/services/api-server/internal/modules/iam/repo"
	appconfig "github.com/liangluo/weouc2026/services/api-server/internal/platform/config"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/persistence"
)

func Bootstrap(ctx context.Context, cfg appconfig.AppConfig, clients *persistence.Clients) error {
	if cfg.Persistence.IAMBackendOrDefault() != "mysql_redis" {
		return nil
	}
	if clients == nil || clients.MySQL == nil {
		return fmt.Errorf("mysql client is required for iam bootstrap")
	}

	if err := iamrepo.AutoMigrateMySQL(ctx, clients.MySQL); err != nil {
		return fmt.Errorf("auto migrate iam mysql schema failed: %w", err)
	}

	return nil
}
