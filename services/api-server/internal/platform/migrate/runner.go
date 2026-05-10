package migrate

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"time"

	migrationfiles "github.com/liangluo/weouc2026/services/api-server/migrations"
)

func Run(ctx context.Context, db *sql.DB) error {
	return runFS(ctx, db, migrationfiles.Files)
}

func runFS(ctx context.Context, db *sql.DB, migrationFS fs.FS) error {
	if db == nil {
		return fmt.Errorf("migration db is nil")
	}

	if err := ensureMigrationsTable(ctx, db); err != nil {
		return err
	}

	entries, err := fs.ReadDir(migrationFS, ".")
	if err != nil {
		return fmt.Errorf("read migration dir failed: %w", err)
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".sql" {
			continue
		}
		names = append(names, entry.Name())
	}
	sort.Strings(names)

	for _, name := range names {
		applied, err := hasApplied(ctx, db, name)
		if err != nil {
			return err
		}
		if applied {
			continue
		}

		content, err := fs.ReadFile(migrationFS, name)
		if err != nil {
			return fmt.Errorf("read migration %s failed: %w", name, err)
		}

		if err := applyMigration(ctx, db, name, string(content)); err != nil {
			return err
		}
	}

	return nil
}

func ensureMigrationsTable(ctx context.Context, db *sql.DB) error {
	const query = `
CREATE TABLE IF NOT EXISTS schema_migrations (
    version TEXT PRIMARY KEY,
    applied_at TIMESTAMPTZ NOT NULL
)`

	if _, err := db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("ensure schema_migrations failed: %w", err)
	}

	return nil
}

func hasApplied(ctx context.Context, db *sql.DB, version string) (bool, error) {
	const query = `SELECT 1 FROM schema_migrations WHERE version = $1`

	var marker int
	if err := db.QueryRowContext(ctx, query, version).Scan(&marker); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("query migration %s failed: %w", version, err)
	}

	return true, nil
}

func applyMigration(ctx context.Context, db *sql.DB, version, content string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin migration %s failed: %w", version, err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, strings.TrimSpace(content)); err != nil {
		return fmt.Errorf("apply migration %s failed: %w", version, err)
	}
	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO schema_migrations (version, applied_at) VALUES ($1, $2)`,
		version,
		time.Now().UTC(),
	); err != nil {
		return fmt.Errorf("record migration %s failed: %w", version, err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit migration %s failed: %w", version, err)
	}

	return nil
}
