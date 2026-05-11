package repo

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	notificationtypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/notification/types"
)

func TestPostgresRepositoryGetMessage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("new sqlmock failed: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id", "title", "content", "category", "target_scope", "target_user_ids", "action_url", "publisher_user_id", "publisher", "created_at", "read_by_user_ids",
	}).AddRow(
		"notification-201",
		"系统提醒",
		"今晚 23 点开始维护。",
		"system",
		"users",
		`["u-1","u-2"]`,
		"/pages/home/index",
		"admin-001",
		"校园运营中心",
		time.Date(2026, 5, 11, 14, 0, 0, 0, time.UTC),
		`{"u-1":"2026-05-11T14:30:00Z"}`,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT ` + messageColumns + ` FROM notification_messages WHERE id = $1 LIMIT 1`)).
		WithArgs("notification-201").
		WillReturnRows(rows)

	repository := NewPostgresRepository(db)
	item, err := repository.GetMessage(context.Background(), "notification-201")
	if err != nil {
		t.Fatalf("GetMessage returned error: %v", err)
	}
	if item.ID != "notification-201" || len(item.TargetUserIDs) != 2 || item.ReadByUserIDs["u-1"].IsZero() {
		t.Fatalf("unexpected notification payload: %+v", item)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestPostgresRepositoryUpdateMessageReturnsNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("new sqlmock failed: %v", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT ` + messageColumns + ` FROM notification_messages WHERE id = $1 FOR UPDATE`)).
		WithArgs("missing-message").
		WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()

	repository := NewPostgresRepository(db)
	_, err = repository.UpdateMessage(context.Background(), "missing-message", func(item *notificationtypes.MessageItem) error {
		return nil
	})
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestPostgresRepositorySaveMessage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("new sqlmock failed: %v", err)
	}
	defer db.Close()

	item := notificationtypes.MessageItem{
		ID:              "notification-301",
		Title:           "发布成功",
		Content:         "内容已发布。",
		Category:        "system",
		TargetScope:     "all",
		TargetUserIDs:   []string{},
		ActionURL:       "/pages/home/index",
		PublisherUserID: "admin-001",
		Publisher:       "管理员",
		CreatedAt:       time.Date(2026, 5, 11, 15, 0, 0, 0, time.UTC),
		ReadByUserIDs:   map[string]time.Time{},
	}

	rows := sqlmock.NewRows([]string{
		"id", "title", "content", "category", "target_scope", "target_user_ids", "action_url", "publisher_user_id", "publisher", "created_at", "read_by_user_ids",
	}).AddRow(
		item.ID,
		item.Title,
		item.Content,
		item.Category,
		item.TargetScope,
		`[]`,
		item.ActionURL,
		item.PublisherUserID,
		item.Publisher,
		item.CreatedAt,
		`{}`,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO notification_messages (`)).
		WithArgs(
			item.ID,
			item.Title,
			item.Content,
			item.Category,
			item.TargetScope,
			`[]`,
			item.ActionURL,
			item.PublisherUserID,
			item.Publisher,
			item.CreatedAt,
			`{}`,
		).
		WillReturnRows(rows)

	repository := NewPostgresRepository(db)
	saved, err := repository.SaveMessage(context.Background(), item)
	if err != nil {
		t.Fatalf("SaveMessage returned error: %v", err)
	}
	if saved.ID != item.ID || len(saved.TargetUserIDs) != 0 {
		t.Fatalf("unexpected saved notification: %+v", saved)
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
		`INSERT INTO notification_id_sequences (name, value) VALUES ($1, 0) ON CONFLICT (name) DO NOTHING`,
	)).
		WithArgs("notification").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(regexp.QuoteMeta(
		`UPDATE notification_id_sequences SET value = value + 1 WHERE name = $1 RETURNING value`,
	)).
		WithArgs("notification").
		WillReturnRows(sqlmock.NewRows([]string{"value"}).AddRow(int64(501)))
	mock.ExpectCommit()

	repository := NewPostgresRepository(db)
	id, err := repository.NextID(context.Background(), "notification")
	if err != nil {
		t.Fatalf("NextID returned error: %v", err)
	}
	if id != "notification-501" {
		t.Fatalf("expected notification-501, got %q", id)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}
