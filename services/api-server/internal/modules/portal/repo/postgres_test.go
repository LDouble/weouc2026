package repo

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	portaltypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/portal/types"
)

func TestPostgresRepositoryGetNotice(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("new sqlmock failed: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id", "title", "summary", "content", "audience", "tags", "pinned", "publisher_user_id", "publisher", "published_at", "created_at",
	}).AddRow(
		"notice-201",
		"停机维护通知",
		"今晚 23 点开始维护。",
		"发布和审核链路将短暂只读。",
		"all",
		`["运维","公告"]`,
		true,
		"admin-001",
		"校园运营中心",
		time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC),
		time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC),
	)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT ` + noticeColumns + ` FROM portal_notices WHERE id = $1 LIMIT 1`)).
		WithArgs("notice-201").
		WillReturnRows(rows)

	repository := NewPostgresRepository(db)
	item, err := repository.GetNotice(context.Background(), "notice-201")
	if err != nil {
		t.Fatalf("GetNotice returned error: %v", err)
	}
	if item.ID != "notice-201" || item.Publisher != "校园运营中心" || len(item.Tags) != 2 {
		t.Fatalf("unexpected notice payload: %+v", item)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestPostgresRepositorySaveNotice(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("new sqlmock failed: %v", err)
	}
	defer db.Close()

	item := portaltypes.NoticeItem{
		ID:              "notice-301",
		Title:           "新公告",
		Summary:         "摘要",
		Content:         "正文",
		Audience:        "student",
		Tags:            []string{"测试"},
		Pinned:          false,
		PublisherUserID: "admin-001",
		Publisher:       "管理员",
		PublishedAt:     time.Date(2026, 5, 11, 13, 0, 0, 0, time.UTC),
		CreatedAt:       time.Date(2026, 5, 11, 13, 0, 0, 0, time.UTC),
	}

	rows := sqlmock.NewRows([]string{
		"id", "title", "summary", "content", "audience", "tags", "pinned", "publisher_user_id", "publisher", "published_at", "created_at",
	}).AddRow(
		item.ID,
		item.Title,
		item.Summary,
		item.Content,
		item.Audience,
		`["测试"]`,
		item.Pinned,
		item.PublisherUserID,
		item.Publisher,
		item.PublishedAt,
		item.CreatedAt,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO portal_notices (`)).
		WithArgs(
			item.ID,
			item.Title,
			item.Summary,
			item.Content,
			item.Audience,
			`["测试"]`,
			item.Pinned,
			item.PublisherUserID,
			item.Publisher,
			item.PublishedAt,
			item.CreatedAt,
		).
		WillReturnRows(rows)

	repository := NewPostgresRepository(db)
	saved, err := repository.SaveNotice(context.Background(), item)
	if err != nil {
		t.Fatalf("SaveNotice returned error: %v", err)
	}
	if saved.ID != item.ID || len(saved.Tags) != 1 || saved.Tags[0] != "测试" {
		t.Fatalf("unexpected saved notice: %+v", saved)
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
		`INSERT INTO portal_id_sequences (name, value) VALUES ($1, 0) ON CONFLICT (name) DO NOTHING`,
	)).
		WithArgs("notice").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(regexp.QuoteMeta(
		`UPDATE portal_id_sequences SET value = value + 1 WHERE name = $1 RETURNING value`,
	)).
		WithArgs("notice").
		WillReturnRows(sqlmock.NewRows([]string{"value"}).AddRow(int64(301)))
	mock.ExpectCommit()

	repository := NewPostgresRepository(db)
	id, err := repository.NextID(context.Background(), "notice")
	if err != nil {
		t.Fatalf("NextID returned error: %v", err)
	}
	if id != "notice-301" {
		t.Fatalf("expected notice-301, got %q", id)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestPostgresRepositoryGetNoticeReturnsNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("new sqlmock failed: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT ` + noticeColumns + ` FROM portal_notices WHERE id = $1 LIMIT 1`)).
		WithArgs("missing-notice").
		WillReturnError(sql.ErrNoRows)

	repository := NewPostgresRepository(db)
	_, err = repository.GetNotice(context.Background(), "missing-notice")
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestPostgresRepositoryListNoticesWithKeyword(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("new sqlmock failed: %v", err)
	}
	defer db.Close()

	query := NoticeListQuery{Keyword: "维护"}
	sqlText, args := buildNoticeListQuery(query)
	rows := sqlmock.NewRows([]string{
		"id", "title", "summary", "content", "audience", "tags", "pinned", "publisher_user_id", "publisher", "published_at", "created_at",
	}).AddRow(
		"notice-301",
		"停机维护通知",
		"今晚维护",
		"发布链路将短暂只读。",
		"all",
		`["运维"]`,
		true,
		"admin-001",
		"校园运营中心",
		time.Date(2026, 5, 12, 12, 0, 0, 0, time.UTC),
		time.Date(2026, 5, 12, 12, 0, 0, 0, time.UTC),
	)

	mock.ExpectQuery(regexp.QuoteMeta(sqlText)).
		WithArgs(toDriverValues(args)...).
		WillReturnRows(rows)

	repository := NewPostgresRepository(db)
	items, err := repository.ListNotices(context.Background(), query)
	if err != nil {
		t.Fatalf("ListNotices returned error: %v", err)
	}
	if len(items) != 1 || items[0].ID != "notice-301" {
		t.Fatalf("unexpected notice payloads: %+v", items)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func toDriverValues(args []any) []driver.Value {
	values := make([]driver.Value, 0, len(args))
	for _, arg := range args {
		values = append(values, arg)
	}
	return values
}
