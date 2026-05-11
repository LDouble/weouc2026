package audit

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

const auditColumns = `
id,
actor_id,
actor_name,
action,
resource_type,
resource_id,
result,
message,
details,
created_at`

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

func (s *PostgresStore) Record(ctx context.Context, entry Entry) error {
	if s.db == nil {
		return fmt.Errorf("postgres audit store db is nil")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin audit record tx failed: %w", err)
	}
	defer tx.Rollback()

	entry.ID, err = nextAuditID(ctx, tx, entry.ID)
	if err != nil {
		return err
	}
	entry.Result = firstNonEmpty(entry.Result, "success")
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = nowUTC()
	}
	if entry.Details != nil {
		entry.Details = cloneDetails(entry.Details)
	}

	var details any
	if entry.Details != nil {
		raw, marshalErr := json.Marshal(entry.Details)
		if marshalErr != nil {
			return fmt.Errorf("marshal audit details failed: %w", marshalErr)
		}
		details = string(raw)
	}

	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO audit_logs (
			id,
			actor_id,
			actor_name,
			action,
			resource_type,
			resource_id,
			result,
			message,
			details,
			created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9::jsonb, $10)`,
		entry.ID,
		entry.ActorID,
		entry.ActorName,
		entry.Action,
		entry.ResourceType,
		entry.ResourceID,
		entry.Result,
		entry.Message,
		details,
		entry.CreatedAt,
	); err != nil {
		return fmt.Errorf("insert audit log failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit audit record tx failed: %w", err)
	}

	return nil
}

func (s *PostgresStore) List(ctx context.Context) ([]Entry, error) {
	if s.db == nil {
		return nil, fmt.Errorf("postgres audit store db is nil")
	}

	rows, err := s.db.QueryContext(ctx, `SELECT `+auditColumns+` FROM audit_logs`)
	if err != nil {
		return nil, fmt.Errorf("list audit logs failed: %w", err)
	}
	defer rows.Close()

	entries := make([]Entry, 0)
	for rows.Next() {
		entry, scanErr := scanEntry(rows)
		if scanErr != nil {
			return nil, fmt.Errorf("scan audit log failed: %w", scanErr)
		}
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate audit logs failed: %w", err)
	}

	return entries, nil
}

type auditExecer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func nextAuditID(ctx context.Context, execer auditExecer, currentID string) (string, error) {
	if currentID != "" {
		return currentID, nil
	}

	if _, err := execer.ExecContext(
		ctx,
		`INSERT INTO audit_id_sequences (name, value) VALUES ($1, 0) ON CONFLICT (name) DO NOTHING`,
		"audit",
	); err != nil {
		return "", fmt.Errorf("ensure audit id sequence failed: %w", err)
	}

	var value int64
	if err := execer.QueryRowContext(
		ctx,
		`UPDATE audit_id_sequences SET value = value + 1 WHERE name = $1 RETURNING value`,
		"audit",
	).Scan(&value); err != nil {
		return "", fmt.Errorf("increment audit id sequence failed: %w", err)
	}

	return fmt.Sprintf("audit-%03d", value), nil
}

type auditScanner interface {
	Scan(dest ...any) error
}

func scanEntry(scanner auditScanner) (Entry, error) {
	var entry Entry
	var detailsRaw []byte

	if err := scanner.Scan(
		&entry.ID,
		&entry.ActorID,
		&entry.ActorName,
		&entry.Action,
		&entry.ResourceType,
		&entry.ResourceID,
		&entry.Result,
		&entry.Message,
		&detailsRaw,
		&entry.CreatedAt,
	); err != nil {
		return Entry{}, err
	}

	if len(detailsRaw) > 0 {
		if err := json.Unmarshal(detailsRaw, &entry.Details); err != nil {
			return Entry{}, fmt.Errorf("unmarshal audit details failed: %w", err)
		}
	}

	return cloneEntry(entry), nil
}

var nowUTC = func() time.Time {
	return time.Now().UTC()
}
