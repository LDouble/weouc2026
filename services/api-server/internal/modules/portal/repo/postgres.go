package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	portaltypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/portal/types"
)

const bannerColumns = `
id,
title,
description,
image_url,
action_url,
sort_order,
created_at`

const noticeColumns = `
id,
title,
summary,
content,
audience,
tags,
pinned,
publisher_user_id,
publisher,
published_at,
created_at`

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) ListBanners(ctx context.Context) ([]portaltypes.BannerItem, error) {
	if r.db == nil {
		return nil, fmt.Errorf("postgres portal repository db is nil")
	}

	rows, err := r.db.QueryContext(ctx, `SELECT `+bannerColumns+` FROM portal_banners`)
	if err != nil {
		return nil, fmt.Errorf("list portal banners failed: %w", err)
	}
	defer rows.Close()

	items := make([]portaltypes.BannerItem, 0)
	for rows.Next() {
		item, scanErr := scanBanner(rows)
		if scanErr != nil {
			return nil, fmt.Errorf("scan portal banner failed: %w", scanErr)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate portal banners failed: %w", err)
	}

	return items, nil
}

func (r *PostgresRepository) GetBanner(ctx context.Context, id string) (portaltypes.BannerItem, error) {
	if r.db == nil {
		return portaltypes.BannerItem{}, fmt.Errorf("postgres portal repository db is nil")
	}

	row := r.db.QueryRowContext(ctx, `SELECT `+bannerColumns+` FROM portal_banners WHERE id = $1 LIMIT 1`, id)
	item, err := scanBanner(row)
	if errors.Is(err, sql.ErrNoRows) {
		return portaltypes.BannerItem{}, ErrNotFound
	}
	if err != nil {
		return portaltypes.BannerItem{}, fmt.Errorf("get portal banner failed: %w", err)
	}

	return item, nil
}

func (r *PostgresRepository) SaveBanner(ctx context.Context, item portaltypes.BannerItem) (portaltypes.BannerItem, error) {
	if r.db == nil {
		return portaltypes.BannerItem{}, fmt.Errorf("postgres portal repository db is nil")
	}

	saved, err := saveBanner(ctx, r.db, item)
	if err != nil {
		return portaltypes.BannerItem{}, fmt.Errorf("save portal banner failed: %w", err)
	}

	return saved, nil
}

func (r *PostgresRepository) DeleteBanner(ctx context.Context, id string) error {
	if r.db == nil {
		return fmt.Errorf("postgres portal repository db is nil")
	}

	result, err := r.db.ExecContext(ctx, `DELETE FROM portal_banners WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete portal banner failed: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read portal banner rows affected failed: %w", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *PostgresRepository) ListNotices(ctx context.Context) ([]portaltypes.NoticeItem, error) {
	if r.db == nil {
		return nil, fmt.Errorf("postgres portal repository db is nil")
	}

	rows, err := r.db.QueryContext(ctx, `SELECT `+noticeColumns+` FROM portal_notices`)
	if err != nil {
		return nil, fmt.Errorf("list portal notices failed: %w", err)
	}
	defer rows.Close()

	items := make([]portaltypes.NoticeItem, 0)
	for rows.Next() {
		item, scanErr := scanNotice(rows)
		if scanErr != nil {
			return nil, fmt.Errorf("scan portal notice failed: %w", scanErr)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate portal notices failed: %w", err)
	}

	return items, nil
}

func (r *PostgresRepository) GetNotice(ctx context.Context, id string) (portaltypes.NoticeItem, error) {
	if r.db == nil {
		return portaltypes.NoticeItem{}, fmt.Errorf("postgres portal repository db is nil")
	}

	row := r.db.QueryRowContext(ctx, `SELECT `+noticeColumns+` FROM portal_notices WHERE id = $1 LIMIT 1`, id)
	item, err := scanNotice(row)
	if errors.Is(err, sql.ErrNoRows) {
		return portaltypes.NoticeItem{}, ErrNotFound
	}
	if err != nil {
		return portaltypes.NoticeItem{}, fmt.Errorf("get portal notice failed: %w", err)
	}

	return item, nil
}

func (r *PostgresRepository) SaveNotice(ctx context.Context, item portaltypes.NoticeItem) (portaltypes.NoticeItem, error) {
	if r.db == nil {
		return portaltypes.NoticeItem{}, fmt.Errorf("postgres portal repository db is nil")
	}

	saved, err := saveNotice(ctx, r.db, item)
	if err != nil {
		return portaltypes.NoticeItem{}, fmt.Errorf("save portal notice failed: %w", err)
	}

	return saved, nil
}

func (r *PostgresRepository) DeleteNotice(ctx context.Context, id string) error {
	if r.db == nil {
		return fmt.Errorf("postgres portal repository db is nil")
	}

	result, err := r.db.ExecContext(ctx, `DELETE FROM portal_notices WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete portal notice failed: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read portal notice rows affected failed: %w", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *PostgresRepository) NextID(ctx context.Context, prefix string) (string, error) {
	if r.db == nil {
		return "", fmt.Errorf("postgres portal repository db is nil")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("begin portal id tx failed: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO portal_id_sequences (name, value) VALUES ($1, 0) ON CONFLICT (name) DO NOTHING`,
		prefix,
	); err != nil {
		return "", fmt.Errorf("ensure portal id sequence failed: %w", err)
	}

	var value int64
	if err := tx.QueryRowContext(
		ctx,
		`UPDATE portal_id_sequences SET value = value + 1 WHERE name = $1 RETURNING value`,
		prefix,
	).Scan(&value); err != nil {
		return "", fmt.Errorf("increment portal id sequence failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("commit portal id tx failed: %w", err)
	}

	return fmt.Sprintf("%s-%03d", prefix, value), nil
}

type portalQueryRower interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func saveBanner(ctx context.Context, querier portalQueryRower, item portaltypes.BannerItem) (portaltypes.BannerItem, error) {
	row := querier.QueryRowContext(
		ctx,
		`INSERT INTO portal_banners (
			id,
			title,
			description,
			image_url,
			action_url,
			sort_order,
			created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			description = EXCLUDED.description,
			image_url = EXCLUDED.image_url,
			action_url = EXCLUDED.action_url,
			sort_order = EXCLUDED.sort_order
		RETURNING `+bannerColumns,
		item.ID,
		item.Title,
		item.Description,
		item.ImageURL,
		item.ActionURL,
		item.Sort,
		item.CreatedAt,
	)

	saved, err := scanBanner(row)
	if err != nil {
		return portaltypes.BannerItem{}, err
	}
	return saved, nil
}

func saveNotice(ctx context.Context, querier portalQueryRower, item portaltypes.NoticeItem) (portaltypes.NoticeItem, error) {
	tags, err := json.Marshal(item.Tags)
	if err != nil {
		return portaltypes.NoticeItem{}, fmt.Errorf("marshal portal notice tags failed: %w", err)
	}

	row := querier.QueryRowContext(
		ctx,
		`INSERT INTO portal_notices (
			id,
			title,
			summary,
			content,
			audience,
			tags,
			pinned,
			publisher_user_id,
			publisher,
			published_at,
			created_at
		) VALUES ($1, $2, $3, $4, $5, $6::jsonb, $7, $8, $9, $10, $11)
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			summary = EXCLUDED.summary,
			content = EXCLUDED.content,
			audience = EXCLUDED.audience,
			tags = EXCLUDED.tags,
			pinned = EXCLUDED.pinned,
			publisher_user_id = EXCLUDED.publisher_user_id,
			publisher = EXCLUDED.publisher,
			published_at = EXCLUDED.published_at
		RETURNING `+noticeColumns,
		item.ID,
		item.Title,
		item.Summary,
		item.Content,
		item.Audience,
		string(tags),
		item.Pinned,
		item.PublisherUserID,
		item.Publisher,
		item.PublishedAt,
		item.CreatedAt,
	)

	saved, err := scanNotice(row)
	if err != nil {
		return portaltypes.NoticeItem{}, err
	}

	return saved, nil
}

type portalScanner interface {
	Scan(dest ...any) error
}

func scanBanner(scanner portalScanner) (portaltypes.BannerItem, error) {
	var item portaltypes.BannerItem
	if err := scanner.Scan(
		&item.ID,
		&item.Title,
		&item.Description,
		&item.ImageURL,
		&item.ActionURL,
		&item.Sort,
		&item.CreatedAt,
	); err != nil {
		return portaltypes.BannerItem{}, err
	}
	return item, nil
}

func scanNotice(scanner portalScanner) (portaltypes.NoticeItem, error) {
	var item portaltypes.NoticeItem
	var tagsRaw []byte

	if err := scanner.Scan(
		&item.ID,
		&item.Title,
		&item.Summary,
		&item.Content,
		&item.Audience,
		&tagsRaw,
		&item.Pinned,
		&item.PublisherUserID,
		&item.Publisher,
		&item.PublishedAt,
		&item.CreatedAt,
	); err != nil {
		return portaltypes.NoticeItem{}, err
	}

	if len(tagsRaw) > 0 {
		if err := json.Unmarshal(tagsRaw, &item.Tags); err != nil {
			return portaltypes.NoticeItem{}, fmt.Errorf("unmarshal portal notice tags failed: %w", err)
		}
	}

	return cloneNotice(item), nil
}
