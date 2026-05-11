package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	cltypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/campus_life/types"
)

func TestPostgresRepositoryGetMarket(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("new sqlmock failed: %v", err)
	}
	defer db.Close()

	extraRaw, err := json.Marshal(cltypes.MarketExtra{
		Category:      "digital",
		Price:         "4299",
		OriginalPrice: "6499",
		Condition:     "9成新",
		TradeMode:     "校内当面交易",
		Contact:       "wx-hd-ipad",
		Images:        []string{"miniapp/market/u-1/20260510/ipad.png"},
	})
	if err != nil {
		t.Fatalf("marshal market extra failed: %v", err)
	}

	rows := sqlmock.NewRows([]string{
		"id", "title", "description", "publisher_user_id", "publisher", "publisher_initial",
		"image", "created_at", "likes", "liked_by_user_ids", "extra",
	}).AddRow(
		"market-201",
		"九成新 iPad Pro 11 寸",
		"M2 芯片，日常记笔记和刷题使用。",
		"user-001",
		"海大同学",
		"海",
		"miniapp/market/u-1/20260510/ipad.png",
		time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC),
		18,
		`{"user-002":true}`,
		extraRaw,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT ` + marketColumns + ` FROM campus_markets WHERE id = $1 LIMIT 1`)).
		WithArgs("market-201").
		WillReturnRows(rows)

	repository := NewPostgresRepository(db)
	item, err := repository.GetMarket(context.Background(), "market-201")
	if err != nil {
		t.Fatalf("GetMarket returned error: %v", err)
	}
	if item.ID != "market-201" || item.Extra.Category != "digital" || !item.LikedByUserIDs["user-002"] {
		t.Fatalf("unexpected market payload: %+v", item)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestPostgresRepositoryUpdateErrandReturnsNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("new sqlmock failed: %v", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT ` + errandColumns + ` FROM campus_errands WHERE id = $1 FOR UPDATE`)).
		WithArgs("missing-errand").
		WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()

	repository := NewPostgresRepository(db)
	_, err = repository.UpdateErrand(context.Background(), "missing-errand", func(item *cltypes.ErrandItem) error {
		return nil
	})
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestPostgresRepositoryNextID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("new sqlmock failed: %v", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO campus_life_id_sequences (name, value) VALUES ($1, 0) ON CONFLICT (name) DO NOTHING`,
	)).
		WithArgs("market").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(regexp.QuoteMeta(
		`UPDATE campus_life_id_sequences SET value = value + 1 WHERE name = $1 RETURNING value`,
	)).
		WithArgs("market").
		WillReturnRows(sqlmock.NewRows([]string{"value"}).AddRow(int64(103)))
	mock.ExpectCommit()

	repository := NewPostgresRepository(db)
	id, err := repository.NextID(context.Background(), "market")
	if err != nil {
		t.Fatalf("NextID returned error: %v", err)
	}
	if id != "market-103" {
		t.Fatalf("expected market-103, got %q", id)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestPostgresRepositoryGetCarpool(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("new sqlmock failed: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id", "category", "route_from", "route_to", "travel_at", "type_label", "seats_text",
		"price_text", "note", "tags", "contact", "review_status", "publisher_user_id",
		"publisher", "publisher_initial", "created_at",
	}).AddRow(
		"carpool-201",
		"tomorrow",
		"海大北门",
		"福州南站",
		time.Date(2026, 5, 12, 18, 30, 0, 0, time.FixedZone("CST", 8*3600)),
		"明日顺路",
		"余座 2",
		"人均 20 元",
		"可带 1 个行李箱",
		`["明天出发","顺路可拼"]`,
		"wx-carpool-201",
		"published",
		"user-201",
		"返校同学",
		"返",
		time.Date(2026, 5, 11, 9, 0, 0, 0, time.UTC),
	)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT ` + carpoolColumns + ` FROM campus_carpools WHERE id = $1 LIMIT 1`)).
		WithArgs("carpool-201").
		WillReturnRows(rows)

	repository := NewPostgresRepository(db)
	item, err := repository.GetCarpool(context.Background(), "carpool-201")
	if err != nil {
		t.Fatalf("GetCarpool returned error: %v", err)
	}
	if item.Category != "tomorrow" || item.From != "海大北门" || len(item.Tags) != 2 {
		t.Fatalf("unexpected carpool payload: %+v", item)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}
