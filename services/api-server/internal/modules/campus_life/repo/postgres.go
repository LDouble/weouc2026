package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	cltypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/campus_life/types"
)

const marketColumns = `
id,
title,
description,
publisher_user_id,
publisher,
publisher_initial,
image,
created_at,
likes,
liked_by_user_ids,
extra`

const errandColumns = `
id,
title,
description,
category,
route_start,
route_end,
deadline,
reward,
contact,
urgent,
images,
status,
publisher_user_id,
publisher,
publisher_initial,
acceptor_user_id,
created_at`

const resourceColumns = `
id,
title,
description,
publisher_user_id,
publisher,
publisher_initial,
created_at,
extra`

const lostFoundColumns = `
id,
title,
description,
publisher_user_id,
publisher,
publisher_initial,
created_at,
extra`

const carpoolColumns = `
id,
category,
route_from,
route_to,
travel_at,
type_label,
seats_text,
price_text,
note,
tags,
contact,
review_status,
publisher_user_id,
publisher,
publisher_initial,
created_at`

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) ListMarkets(ctx context.Context) ([]cltypes.MarketItem, error) {
	if r.db == nil {
		return nil, fmt.Errorf("postgres campus_life repository db is nil")
	}

	rows, err := r.db.QueryContext(ctx, `SELECT `+marketColumns+` FROM campus_markets`)
	if err != nil {
		return nil, fmt.Errorf("list markets failed: %w", err)
	}
	defer rows.Close()

	items := make([]cltypes.MarketItem, 0)
	for rows.Next() {
		item, scanErr := scanMarket(rows)
		if scanErr != nil {
			return nil, fmt.Errorf("scan market failed: %w", scanErr)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate markets failed: %w", err)
	}

	return items, nil
}

func (r *PostgresRepository) GetMarket(ctx context.Context, id string) (cltypes.MarketItem, error) {
	if r.db == nil {
		return cltypes.MarketItem{}, fmt.Errorf("postgres campus_life repository db is nil")
	}

	row := r.db.QueryRowContext(ctx, `SELECT `+marketColumns+` FROM campus_markets WHERE id = $1 LIMIT 1`, id)
	item, err := scanMarket(row)
	if errors.Is(err, sql.ErrNoRows) {
		return cltypes.MarketItem{}, ErrNotFound
	}
	if err != nil {
		return cltypes.MarketItem{}, fmt.Errorf("get market failed: %w", err)
	}

	return item, nil
}

func (r *PostgresRepository) SaveMarket(ctx context.Context, item cltypes.MarketItem) (cltypes.MarketItem, error) {
	if r.db == nil {
		return cltypes.MarketItem{}, fmt.Errorf("postgres campus_life repository db is nil")
	}

	saved, err := saveMarket(ctx, r.db, item)
	if err != nil {
		return cltypes.MarketItem{}, fmt.Errorf("save market failed: %w", err)
	}

	return saved, nil
}

func (r *PostgresRepository) UpdateMarket(
	ctx context.Context,
	id string,
	mutate func(*cltypes.MarketItem) error,
) (cltypes.MarketItem, error) {
	if r.db == nil {
		return cltypes.MarketItem{}, fmt.Errorf("postgres campus_life repository db is nil")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return cltypes.MarketItem{}, fmt.Errorf("begin market update tx failed: %w", err)
	}
	defer tx.Rollback()

	current, err := scanMarket(tx.QueryRowContext(ctx, `SELECT `+marketColumns+` FROM campus_markets WHERE id = $1 FOR UPDATE`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return cltypes.MarketItem{}, ErrNotFound
	}
	if err != nil {
		return cltypes.MarketItem{}, fmt.Errorf("load market for update failed: %w", err)
	}
	if err := mutate(&current); err != nil {
		return cltypes.MarketItem{}, err
	}

	updated, err := saveMarket(ctx, tx, current)
	if err != nil {
		return cltypes.MarketItem{}, fmt.Errorf("update market failed: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return cltypes.MarketItem{}, fmt.Errorf("commit market update failed: %w", err)
	}

	return updated, nil
}

func (r *PostgresRepository) ListErrands(ctx context.Context) ([]cltypes.ErrandItem, error) {
	if r.db == nil {
		return nil, fmt.Errorf("postgres campus_life repository db is nil")
	}

	rows, err := r.db.QueryContext(ctx, `SELECT `+errandColumns+` FROM campus_errands`)
	if err != nil {
		return nil, fmt.Errorf("list errands failed: %w", err)
	}
	defer rows.Close()

	items := make([]cltypes.ErrandItem, 0)
	for rows.Next() {
		item, scanErr := scanErrand(rows)
		if scanErr != nil {
			return nil, fmt.Errorf("scan errand failed: %w", scanErr)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate errands failed: %w", err)
	}

	return items, nil
}

func (r *PostgresRepository) GetErrand(ctx context.Context, id string) (cltypes.ErrandItem, error) {
	if r.db == nil {
		return cltypes.ErrandItem{}, fmt.Errorf("postgres campus_life repository db is nil")
	}

	row := r.db.QueryRowContext(ctx, `SELECT `+errandColumns+` FROM campus_errands WHERE id = $1 LIMIT 1`, id)
	item, err := scanErrand(row)
	if errors.Is(err, sql.ErrNoRows) {
		return cltypes.ErrandItem{}, ErrNotFound
	}
	if err != nil {
		return cltypes.ErrandItem{}, fmt.Errorf("get errand failed: %w", err)
	}

	return item, nil
}

func (r *PostgresRepository) SaveErrand(ctx context.Context, item cltypes.ErrandItem) (cltypes.ErrandItem, error) {
	if r.db == nil {
		return cltypes.ErrandItem{}, fmt.Errorf("postgres campus_life repository db is nil")
	}

	saved, err := saveErrand(ctx, r.db, item)
	if err != nil {
		return cltypes.ErrandItem{}, fmt.Errorf("save errand failed: %w", err)
	}

	return saved, nil
}

func (r *PostgresRepository) UpdateErrand(
	ctx context.Context,
	id string,
	mutate func(*cltypes.ErrandItem) error,
) (cltypes.ErrandItem, error) {
	if r.db == nil {
		return cltypes.ErrandItem{}, fmt.Errorf("postgres campus_life repository db is nil")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return cltypes.ErrandItem{}, fmt.Errorf("begin errand update tx failed: %w", err)
	}
	defer tx.Rollback()

	current, err := scanErrand(tx.QueryRowContext(ctx, `SELECT `+errandColumns+` FROM campus_errands WHERE id = $1 FOR UPDATE`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return cltypes.ErrandItem{}, ErrNotFound
	}
	if err != nil {
		return cltypes.ErrandItem{}, fmt.Errorf("load errand for update failed: %w", err)
	}
	if err := mutate(&current); err != nil {
		return cltypes.ErrandItem{}, err
	}

	updated, err := saveErrand(ctx, tx, current)
	if err != nil {
		return cltypes.ErrandItem{}, fmt.Errorf("update errand failed: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return cltypes.ErrandItem{}, fmt.Errorf("commit errand update failed: %w", err)
	}

	return updated, nil
}

func (r *PostgresRepository) ListResources(ctx context.Context) ([]cltypes.ResourceItem, error) {
	if r.db == nil {
		return nil, fmt.Errorf("postgres campus_life repository db is nil")
	}

	rows, err := r.db.QueryContext(ctx, `SELECT `+resourceColumns+` FROM campus_resources`)
	if err != nil {
		return nil, fmt.Errorf("list resources failed: %w", err)
	}
	defer rows.Close()

	items := make([]cltypes.ResourceItem, 0)
	for rows.Next() {
		item, scanErr := scanResource(rows)
		if scanErr != nil {
			return nil, fmt.Errorf("scan resource failed: %w", scanErr)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate resources failed: %w", err)
	}

	return items, nil
}

func (r *PostgresRepository) GetResource(ctx context.Context, id string) (cltypes.ResourceItem, error) {
	if r.db == nil {
		return cltypes.ResourceItem{}, fmt.Errorf("postgres campus_life repository db is nil")
	}

	row := r.db.QueryRowContext(ctx, `SELECT `+resourceColumns+` FROM campus_resources WHERE id = $1 LIMIT 1`, id)
	item, err := scanResource(row)
	if errors.Is(err, sql.ErrNoRows) {
		return cltypes.ResourceItem{}, ErrNotFound
	}
	if err != nil {
		return cltypes.ResourceItem{}, fmt.Errorf("get resource failed: %w", err)
	}

	return item, nil
}

func (r *PostgresRepository) SaveResource(ctx context.Context, item cltypes.ResourceItem) (cltypes.ResourceItem, error) {
	if r.db == nil {
		return cltypes.ResourceItem{}, fmt.Errorf("postgres campus_life repository db is nil")
	}

	saved, err := saveResource(ctx, r.db, item)
	if err != nil {
		return cltypes.ResourceItem{}, fmt.Errorf("save resource failed: %w", err)
	}

	return saved, nil
}

func (r *PostgresRepository) ListLostFound(ctx context.Context) ([]cltypes.LostFoundItem, error) {
	if r.db == nil {
		return nil, fmt.Errorf("postgres campus_life repository db is nil")
	}

	rows, err := r.db.QueryContext(ctx, `SELECT `+lostFoundColumns+` FROM campus_lost_founds`)
	if err != nil {
		return nil, fmt.Errorf("list lost_found failed: %w", err)
	}
	defer rows.Close()

	items := make([]cltypes.LostFoundItem, 0)
	for rows.Next() {
		item, scanErr := scanLostFound(rows)
		if scanErr != nil {
			return nil, fmt.Errorf("scan lost_found failed: %w", scanErr)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate lost_found failed: %w", err)
	}

	return items, nil
}

func (r *PostgresRepository) GetLostFound(ctx context.Context, id string) (cltypes.LostFoundItem, error) {
	if r.db == nil {
		return cltypes.LostFoundItem{}, fmt.Errorf("postgres campus_life repository db is nil")
	}

	row := r.db.QueryRowContext(ctx, `SELECT `+lostFoundColumns+` FROM campus_lost_founds WHERE id = $1 LIMIT 1`, id)
	item, err := scanLostFound(row)
	if errors.Is(err, sql.ErrNoRows) {
		return cltypes.LostFoundItem{}, ErrNotFound
	}
	if err != nil {
		return cltypes.LostFoundItem{}, fmt.Errorf("get lost_found failed: %w", err)
	}

	return item, nil
}

func (r *PostgresRepository) SaveLostFound(ctx context.Context, item cltypes.LostFoundItem) (cltypes.LostFoundItem, error) {
	if r.db == nil {
		return cltypes.LostFoundItem{}, fmt.Errorf("postgres campus_life repository db is nil")
	}

	saved, err := saveLostFound(ctx, r.db, item)
	if err != nil {
		return cltypes.LostFoundItem{}, fmt.Errorf("save lost_found failed: %w", err)
	}

	return saved, nil
}

func (r *PostgresRepository) ListCarpools(ctx context.Context) ([]cltypes.CarpoolItem, error) {
	if r.db == nil {
		return nil, fmt.Errorf("postgres campus_life repository db is nil")
	}

	rows, err := r.db.QueryContext(ctx, `SELECT `+carpoolColumns+` FROM campus_carpools`)
	if err != nil {
		return nil, fmt.Errorf("list carpools failed: %w", err)
	}
	defer rows.Close()

	items := make([]cltypes.CarpoolItem, 0)
	for rows.Next() {
		item, scanErr := scanCarpool(rows)
		if scanErr != nil {
			return nil, fmt.Errorf("scan carpool failed: %w", scanErr)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate carpools failed: %w", err)
	}

	return items, nil
}

func (r *PostgresRepository) GetCarpool(ctx context.Context, id string) (cltypes.CarpoolItem, error) {
	if r.db == nil {
		return cltypes.CarpoolItem{}, fmt.Errorf("postgres campus_life repository db is nil")
	}

	row := r.db.QueryRowContext(ctx, `SELECT `+carpoolColumns+` FROM campus_carpools WHERE id = $1 LIMIT 1`, id)
	item, err := scanCarpool(row)
	if errors.Is(err, sql.ErrNoRows) {
		return cltypes.CarpoolItem{}, ErrNotFound
	}
	if err != nil {
		return cltypes.CarpoolItem{}, fmt.Errorf("get carpool failed: %w", err)
	}

	return item, nil
}

func (r *PostgresRepository) SaveCarpool(ctx context.Context, item cltypes.CarpoolItem) (cltypes.CarpoolItem, error) {
	if r.db == nil {
		return cltypes.CarpoolItem{}, fmt.Errorf("postgres campus_life repository db is nil")
	}

	saved, err := saveCarpool(ctx, r.db, item)
	if err != nil {
		return cltypes.CarpoolItem{}, fmt.Errorf("save carpool failed: %w", err)
	}

	return saved, nil
}

func (r *PostgresRepository) NextID(ctx context.Context, prefix string) (string, error) {
	if r.db == nil {
		return "", fmt.Errorf("postgres campus_life repository db is nil")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("begin next id tx failed: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO campus_life_id_sequences (name, value) VALUES ($1, 0) ON CONFLICT (name) DO NOTHING`,
		prefix,
	); err != nil {
		return "", fmt.Errorf("ensure next id sequence failed: %w", err)
	}

	var nextValue int64
	if err := tx.QueryRowContext(
		ctx,
		`UPDATE campus_life_id_sequences SET value = value + 1 WHERE name = $1 RETURNING value`,
		prefix,
	).Scan(&nextValue); err != nil {
		return "", fmt.Errorf("increment next id sequence failed: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("commit next id tx failed: %w", err)
	}

	return fmt.Sprintf("%s-%d", prefix, nextValue), nil
}

type queryRower interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type rowScanner interface {
	Scan(dest ...any) error
}

func saveMarket(ctx context.Context, querier queryRower, item cltypes.MarketItem) (cltypes.MarketItem, error) {
	likedByRaw, err := json.Marshal(item.LikedByUserIDs)
	if err != nil {
		return cltypes.MarketItem{}, fmt.Errorf("marshal market likes failed: %w", err)
	}
	extraRaw, err := json.Marshal(item.Extra)
	if err != nil {
		return cltypes.MarketItem{}, fmt.Errorf("marshal market extra failed: %w", err)
	}

	row := querier.QueryRowContext(
		ctx,
		`INSERT INTO campus_markets (
			id,
			title,
			description,
			publisher_user_id,
			publisher,
			publisher_initial,
			image,
			created_at,
			likes,
			liked_by_user_ids,
			extra
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10::jsonb, $11::jsonb)
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			description = EXCLUDED.description,
			publisher_user_id = EXCLUDED.publisher_user_id,
			publisher = EXCLUDED.publisher,
			publisher_initial = EXCLUDED.publisher_initial,
			image = EXCLUDED.image,
			created_at = EXCLUDED.created_at,
			likes = EXCLUDED.likes,
			liked_by_user_ids = EXCLUDED.liked_by_user_ids,
			extra = EXCLUDED.extra
		RETURNING `+marketColumns,
		item.ID,
		item.Title,
		item.Desc,
		item.PublisherUserID,
		item.Publisher,
		item.PublisherInitial,
		item.Image,
		item.CreatedAt,
		item.Likes,
		string(likedByRaw),
		string(extraRaw),
	)

	return scanMarket(row)
}

func scanMarket(scanner rowScanner) (cltypes.MarketItem, error) {
	var item cltypes.MarketItem
	var likedByRaw []byte
	var extraRaw []byte
	if err := scanner.Scan(
		&item.ID,
		&item.Title,
		&item.Desc,
		&item.PublisherUserID,
		&item.Publisher,
		&item.PublisherInitial,
		&item.Image,
		&item.CreatedAt,
		&item.Likes,
		&likedByRaw,
		&extraRaw,
	); err != nil {
		return cltypes.MarketItem{}, err
	}
	if len(likedByRaw) > 0 {
		if err := json.Unmarshal(likedByRaw, &item.LikedByUserIDs); err != nil {
			return cltypes.MarketItem{}, fmt.Errorf("unmarshal market likes failed: %w", err)
		}
	}
	if item.LikedByUserIDs == nil {
		item.LikedByUserIDs = map[string]bool{}
	}
	if len(extraRaw) > 0 {
		if err := json.Unmarshal(extraRaw, &item.Extra); err != nil {
			return cltypes.MarketItem{}, fmt.Errorf("unmarshal market extra failed: %w", err)
		}
	}

	return item, nil
}

func saveErrand(ctx context.Context, querier queryRower, item cltypes.ErrandItem) (cltypes.ErrandItem, error) {
	imagesRaw, err := json.Marshal(item.Images)
	if err != nil {
		return cltypes.ErrandItem{}, fmt.Errorf("marshal errand images failed: %w", err)
	}

	row := querier.QueryRowContext(
		ctx,
		`INSERT INTO campus_errands (
			id,
			title,
			description,
			category,
			route_start,
			route_end,
			deadline,
			reward,
			contact,
			urgent,
			images,
			status,
			publisher_user_id,
			publisher,
			publisher_initial,
			acceptor_user_id,
			created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11::jsonb, $12, $13, $14, $15, $16, $17)
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			description = EXCLUDED.description,
			category = EXCLUDED.category,
			route_start = EXCLUDED.route_start,
			route_end = EXCLUDED.route_end,
			deadline = EXCLUDED.deadline,
			reward = EXCLUDED.reward,
			contact = EXCLUDED.contact,
			urgent = EXCLUDED.urgent,
			images = EXCLUDED.images,
			status = EXCLUDED.status,
			publisher_user_id = EXCLUDED.publisher_user_id,
			publisher = EXCLUDED.publisher,
			publisher_initial = EXCLUDED.publisher_initial,
			acceptor_user_id = EXCLUDED.acceptor_user_id,
			created_at = EXCLUDED.created_at
		RETURNING `+errandColumns,
		item.ID,
		item.Title,
		item.Desc,
		item.Category,
		item.RouteStart,
		item.RouteEnd,
		item.Deadline,
		item.Reward,
		item.Contact,
		item.Urgent,
		string(imagesRaw),
		item.Status,
		item.PublisherUserID,
		item.Publisher,
		item.PublisherInitial,
		item.AcceptorUserID,
		item.CreatedAt,
	)

	return scanErrand(row)
}

func scanErrand(scanner rowScanner) (cltypes.ErrandItem, error) {
	var item cltypes.ErrandItem
	var imagesRaw []byte
	if err := scanner.Scan(
		&item.ID,
		&item.Title,
		&item.Desc,
		&item.Category,
		&item.RouteStart,
		&item.RouteEnd,
		&item.Deadline,
		&item.Reward,
		&item.Contact,
		&item.Urgent,
		&imagesRaw,
		&item.Status,
		&item.PublisherUserID,
		&item.Publisher,
		&item.PublisherInitial,
		&item.AcceptorUserID,
		&item.CreatedAt,
	); err != nil {
		return cltypes.ErrandItem{}, err
	}
	if len(imagesRaw) > 0 {
		if err := json.Unmarshal(imagesRaw, &item.Images); err != nil {
			return cltypes.ErrandItem{}, fmt.Errorf("unmarshal errand images failed: %w", err)
		}
	}

	return item, nil
}

func saveResource(ctx context.Context, querier queryRower, item cltypes.ResourceItem) (cltypes.ResourceItem, error) {
	extraRaw, err := json.Marshal(item.Extra)
	if err != nil {
		return cltypes.ResourceItem{}, fmt.Errorf("marshal resource extra failed: %w", err)
	}

	row := querier.QueryRowContext(
		ctx,
		`INSERT INTO campus_resources (
			id,
			title,
			description,
			publisher_user_id,
			publisher,
			publisher_initial,
			created_at,
			extra
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8::jsonb)
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			description = EXCLUDED.description,
			publisher_user_id = EXCLUDED.publisher_user_id,
			publisher = EXCLUDED.publisher,
			publisher_initial = EXCLUDED.publisher_initial,
			created_at = EXCLUDED.created_at,
			extra = EXCLUDED.extra
		RETURNING `+resourceColumns,
		item.ID,
		item.Title,
		item.Desc,
		item.PublisherUserID,
		item.Publisher,
		item.PublisherInitial,
		item.CreatedAt,
		string(extraRaw),
	)

	return scanResource(row)
}

func scanResource(scanner rowScanner) (cltypes.ResourceItem, error) {
	var item cltypes.ResourceItem
	var extraRaw []byte
	if err := scanner.Scan(
		&item.ID,
		&item.Title,
		&item.Desc,
		&item.PublisherUserID,
		&item.Publisher,
		&item.PublisherInitial,
		&item.CreatedAt,
		&extraRaw,
	); err != nil {
		return cltypes.ResourceItem{}, err
	}
	if len(extraRaw) > 0 {
		if err := json.Unmarshal(extraRaw, &item.Extra); err != nil {
			return cltypes.ResourceItem{}, fmt.Errorf("unmarshal resource extra failed: %w", err)
		}
	}

	return item, nil
}

func saveLostFound(ctx context.Context, querier queryRower, item cltypes.LostFoundItem) (cltypes.LostFoundItem, error) {
	extraRaw, err := json.Marshal(item.Extra)
	if err != nil {
		return cltypes.LostFoundItem{}, fmt.Errorf("marshal lost_found extra failed: %w", err)
	}

	row := querier.QueryRowContext(
		ctx,
		`INSERT INTO campus_lost_founds (
			id,
			title,
			description,
			publisher_user_id,
			publisher,
			publisher_initial,
			created_at,
			extra
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8::jsonb)
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			description = EXCLUDED.description,
			publisher_user_id = EXCLUDED.publisher_user_id,
			publisher = EXCLUDED.publisher,
			publisher_initial = EXCLUDED.publisher_initial,
			created_at = EXCLUDED.created_at,
			extra = EXCLUDED.extra
		RETURNING `+lostFoundColumns,
		item.ID,
		item.Title,
		item.Desc,
		item.PublisherUserID,
		item.Publisher,
		item.PublisherInitial,
		item.CreatedAt,
		string(extraRaw),
	)

	return scanLostFound(row)
}

func saveCarpool(ctx context.Context, querier queryRower, item cltypes.CarpoolItem) (cltypes.CarpoolItem, error) {
	tagsRaw, err := json.Marshal(item.Tags)
	if err != nil {
		return cltypes.CarpoolItem{}, fmt.Errorf("marshal carpool tags failed: %w", err)
	}

	row := querier.QueryRowContext(
		ctx,
		`INSERT INTO campus_carpools (
			id,
			category,
			route_from,
			route_to,
			travel_at,
			type_label,
			seats_text,
			price_text,
			note,
			tags,
			contact,
			review_status,
			publisher_user_id,
			publisher,
			publisher_initial,
			created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10::jsonb, $11, $12, $13, $14, $15, $16)
		ON CONFLICT (id) DO UPDATE SET
			category = EXCLUDED.category,
			route_from = EXCLUDED.route_from,
			route_to = EXCLUDED.route_to,
			travel_at = EXCLUDED.travel_at,
			type_label = EXCLUDED.type_label,
			seats_text = EXCLUDED.seats_text,
			price_text = EXCLUDED.price_text,
			note = EXCLUDED.note,
			tags = EXCLUDED.tags,
			contact = EXCLUDED.contact,
			review_status = EXCLUDED.review_status,
			publisher_user_id = EXCLUDED.publisher_user_id,
			publisher = EXCLUDED.publisher,
			publisher_initial = EXCLUDED.publisher_initial,
			created_at = EXCLUDED.created_at
		RETURNING `+carpoolColumns,
		item.ID,
		item.Category,
		item.From,
		item.To,
		item.TravelAt,
		item.Type,
		item.SeatsText,
		item.Price,
		item.Note,
		string(tagsRaw),
		item.Contact,
		item.ReviewStatus,
		item.PublisherUserID,
		item.Publisher,
		item.PublisherInitial,
		item.CreatedAt,
	)

	return scanCarpool(row)
}

func scanLostFound(scanner rowScanner) (cltypes.LostFoundItem, error) {
	var item cltypes.LostFoundItem
	var extraRaw []byte
	if err := scanner.Scan(
		&item.ID,
		&item.Title,
		&item.Desc,
		&item.PublisherUserID,
		&item.Publisher,
		&item.PublisherInitial,
		&item.CreatedAt,
		&extraRaw,
	); err != nil {
		return cltypes.LostFoundItem{}, err
	}
	if len(extraRaw) > 0 {
		if err := json.Unmarshal(extraRaw, &item.Extra); err != nil {
			return cltypes.LostFoundItem{}, fmt.Errorf("unmarshal lost_found extra failed: %w", err)
		}
	}

	return item, nil
}

func scanCarpool(scanner rowScanner) (cltypes.CarpoolItem, error) {
	var item cltypes.CarpoolItem
	var tagsRaw []byte
	if err := scanner.Scan(
		&item.ID,
		&item.Category,
		&item.From,
		&item.To,
		&item.TravelAt,
		&item.Type,
		&item.SeatsText,
		&item.Price,
		&item.Note,
		&tagsRaw,
		&item.Contact,
		&item.ReviewStatus,
		&item.PublisherUserID,
		&item.Publisher,
		&item.PublisherInitial,
		&item.CreatedAt,
	); err != nil {
		return cltypes.CarpoolItem{}, err
	}
	if len(tagsRaw) > 0 {
		if err := json.Unmarshal(tagsRaw, &item.Tags); err != nil {
			return cltypes.CarpoolItem{}, fmt.Errorf("unmarshal carpool tags failed: %w", err)
		}
	}

	return item, nil
}
