package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	notificationtypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/notification/types"
)

const messageColumns = `
id,
title,
content,
category,
target_scope,
target_user_ids,
action_url,
publisher_user_id,
publisher,
created_at,
read_by_user_ids`

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) ListMessages(ctx context.Context) ([]notificationtypes.MessageItem, error) {
	if r.db == nil {
		return nil, fmt.Errorf("postgres notification repository db is nil")
	}

	rows, err := r.db.QueryContext(ctx, `SELECT `+messageColumns+` FROM notification_messages`)
	if err != nil {
		return nil, fmt.Errorf("list notification messages failed: %w", err)
	}
	defer rows.Close()

	items := make([]notificationtypes.MessageItem, 0)
	for rows.Next() {
		item, scanErr := scanMessage(rows)
		if scanErr != nil {
			return nil, fmt.Errorf("scan notification message failed: %w", scanErr)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate notification messages failed: %w", err)
	}

	return items, nil
}

func (r *PostgresRepository) GetMessage(ctx context.Context, id string) (notificationtypes.MessageItem, error) {
	if r.db == nil {
		return notificationtypes.MessageItem{}, fmt.Errorf("postgres notification repository db is nil")
	}

	row := r.db.QueryRowContext(ctx, `SELECT `+messageColumns+` FROM notification_messages WHERE id = $1 LIMIT 1`, id)
	item, err := scanMessage(row)
	if errors.Is(err, sql.ErrNoRows) {
		return notificationtypes.MessageItem{}, ErrNotFound
	}
	if err != nil {
		return notificationtypes.MessageItem{}, fmt.Errorf("get notification message failed: %w", err)
	}

	return item, nil
}

func (r *PostgresRepository) SaveMessage(
	ctx context.Context,
	item notificationtypes.MessageItem,
) (notificationtypes.MessageItem, error) {
	if r.db == nil {
		return notificationtypes.MessageItem{}, fmt.Errorf("postgres notification repository db is nil")
	}

	saved, err := saveMessage(ctx, r.db, item)
	if err != nil {
		return notificationtypes.MessageItem{}, fmt.Errorf("save notification message failed: %w", err)
	}

	return saved, nil
}

func (r *PostgresRepository) UpdateMessage(
	ctx context.Context,
	id string,
	mutate func(*notificationtypes.MessageItem) error,
) (notificationtypes.MessageItem, error) {
	if r.db == nil {
		return notificationtypes.MessageItem{}, fmt.Errorf("postgres notification repository db is nil")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return notificationtypes.MessageItem{}, fmt.Errorf("begin notification update tx failed: %w", err)
	}
	defer tx.Rollback()

	current, err := scanMessage(
		tx.QueryRowContext(ctx, `SELECT `+messageColumns+` FROM notification_messages WHERE id = $1 FOR UPDATE`, id),
	)
	if errors.Is(err, sql.ErrNoRows) {
		return notificationtypes.MessageItem{}, ErrNotFound
	}
	if err != nil {
		return notificationtypes.MessageItem{}, fmt.Errorf("load notification message for update failed: %w", err)
	}
	if err := mutate(&current); err != nil {
		return notificationtypes.MessageItem{}, err
	}

	updated, err := saveMessage(ctx, tx, current)
	if err != nil {
		return notificationtypes.MessageItem{}, fmt.Errorf("update notification message failed: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return notificationtypes.MessageItem{}, fmt.Errorf("commit notification update failed: %w", err)
	}

	return updated, nil
}

func (r *PostgresRepository) NextID(ctx context.Context, prefix string) (string, error) {
	if r.db == nil {
		return "", fmt.Errorf("postgres notification repository db is nil")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("begin notification id tx failed: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO notification_id_sequences (name, value) VALUES ($1, 0) ON CONFLICT (name) DO NOTHING`,
		prefix,
	); err != nil {
		return "", fmt.Errorf("ensure notification id sequence failed: %w", err)
	}

	var value int64
	if err := tx.QueryRowContext(
		ctx,
		`UPDATE notification_id_sequences SET value = value + 1 WHERE name = $1 RETURNING value`,
		prefix,
	).Scan(&value); err != nil {
		return "", fmt.Errorf("increment notification id sequence failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("commit notification id tx failed: %w", err)
	}

	return fmt.Sprintf("%s-%03d", prefix, value), nil
}

type notificationQueryRower interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func saveMessage(
	ctx context.Context,
	querier notificationQueryRower,
	item notificationtypes.MessageItem,
) (notificationtypes.MessageItem, error) {
	targetUserIDs, err := json.Marshal(item.TargetUserIDs)
	if err != nil {
		return notificationtypes.MessageItem{}, fmt.Errorf("marshal notification target users failed: %w", err)
	}
	readByUserIDs, err := json.Marshal(item.ReadByUserIDs)
	if err != nil {
		return notificationtypes.MessageItem{}, fmt.Errorf("marshal notification read state failed: %w", err)
	}

	row := querier.QueryRowContext(
		ctx,
		`INSERT INTO notification_messages (
			id,
			title,
			content,
			category,
			target_scope,
			target_user_ids,
			action_url,
			publisher_user_id,
			publisher,
			created_at,
			read_by_user_ids
		) VALUES ($1, $2, $3, $4, $5, $6::jsonb, $7, $8, $9, $10, $11::jsonb)
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			content = EXCLUDED.content,
			category = EXCLUDED.category,
			target_scope = EXCLUDED.target_scope,
			target_user_ids = EXCLUDED.target_user_ids,
			action_url = EXCLUDED.action_url,
			publisher_user_id = EXCLUDED.publisher_user_id,
			publisher = EXCLUDED.publisher,
			read_by_user_ids = EXCLUDED.read_by_user_ids
		RETURNING `+messageColumns,
		item.ID,
		item.Title,
		item.Content,
		item.Category,
		item.TargetScope,
		string(targetUserIDs),
		item.ActionURL,
		item.PublisherUserID,
		item.Publisher,
		item.CreatedAt,
		string(readByUserIDs),
	)

	saved, err := scanMessage(row)
	if err != nil {
		return notificationtypes.MessageItem{}, err
	}

	return saved, nil
}

type notificationScanner interface {
	Scan(dest ...any) error
}

func scanMessage(scanner notificationScanner) (notificationtypes.MessageItem, error) {
	var item notificationtypes.MessageItem
	var targetUserIDsRaw []byte
	var readByUserIDsRaw []byte

	if err := scanner.Scan(
		&item.ID,
		&item.Title,
		&item.Content,
		&item.Category,
		&item.TargetScope,
		&targetUserIDsRaw,
		&item.ActionURL,
		&item.PublisherUserID,
		&item.Publisher,
		&item.CreatedAt,
		&readByUserIDsRaw,
	); err != nil {
		return notificationtypes.MessageItem{}, err
	}

	if len(targetUserIDsRaw) > 0 {
		if err := json.Unmarshal(targetUserIDsRaw, &item.TargetUserIDs); err != nil {
			return notificationtypes.MessageItem{}, fmt.Errorf("unmarshal notification target users failed: %w", err)
		}
	}
	if len(readByUserIDsRaw) > 0 {
		if err := json.Unmarshal(readByUserIDsRaw, &item.ReadByUserIDs); err != nil {
			return notificationtypes.MessageItem{}, fmt.Errorf("unmarshal notification read state failed: %w", err)
		}
	}

	return cloneMessage(item), nil
}
