package audit

import (
	"context"
	"database/sql/driver"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPostgresStoreRecord(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("new sqlmock failed: %v", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO audit_id_sequences (name, value) VALUES ($1, 0) ON CONFLICT (name) DO NOTHING`,
	)).
		WithArgs("audit").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(regexp.QuoteMeta(
		`UPDATE audit_id_sequences SET value = value + 1 WHERE name = $1 RETURNING value`,
	)).
		WithArgs("audit").
		WillReturnRows(sqlmock.NewRows([]string{"value"}).AddRow(int64(901)))
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO audit_logs (`)).
		WithArgs(
			"audit-901",
			"user-001",
			"海大同学",
			"portal.notice.publish",
			"portal_notice",
			"notice-301",
			"success",
			"门户公告发布成功",
			`{"audience":"all"}`,
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	store := NewPostgresStore(db)
	err = store.Record(context.Background(), Entry{
		ActorID:      "user-001",
		ActorName:    "海大同学",
		Action:       "portal.notice.publish",
		ResourceType: "portal_notice",
		ResourceID:   "notice-301",
		Message:      "门户公告发布成功",
		Details: map[string]any{
			"audience": "all",
		},
	})
	if err != nil {
		t.Fatalf("Record returned error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestPostgresStoreList(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("new sqlmock failed: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id", "actor_id", "actor_name", "action", "resource_type", "resource_id", "result", "message", "details", "created_at",
	}).AddRow(
		"audit-901",
		"user-001",
		"海大同学",
		"notification.read",
		"notification_message",
		"notification-101",
		"success",
		"通知标记为已读",
		`{"category":"system"}`,
		time.Date(2026, 5, 11, 16, 0, 0, 0, time.UTC),
	)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT ` + auditColumns + ` FROM audit_logs`)).
		WillReturnRows(rows)

	store := NewPostgresStore(db)
	entries, err := store.List(context.Background(), ListQuery{})
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(entries) != 1 || entries[0].Action != "notification.read" || entries[0].Details["category"] != "system" {
		t.Fatalf("unexpected audit entries: %+v", entries)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestPostgresStoreListWithQuery(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("new sqlmock failed: %v", err)
	}
	defer db.Close()

	query := ListQuery{
		ActorID: "user-001",
		Action:  "notification.read",
	}
	sqlText, args := buildAuditListQuery(query)
	rows := sqlmock.NewRows([]string{
		"id", "actor_id", "actor_name", "action", "resource_type", "resource_id", "result", "message", "details", "created_at",
	}).AddRow(
		"audit-902",
		"user-001",
		"海大同学",
		"notification.read",
		"notification_message",
		"notification-102",
		"success",
		"通知已读",
		`{"category":"system"}`,
		time.Date(2026, 5, 11, 17, 0, 0, 0, time.UTC),
	)

	mock.ExpectQuery(regexp.QuoteMeta(sqlText)).
		WithArgs(toDriverValues(args)...).
		WillReturnRows(rows)

	store := NewPostgresStore(db)
	entries, err := store.List(context.Background(), query)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(entries) != 1 || entries[0].ResourceID != "notification-102" {
		t.Fatalf("unexpected audit entries: %+v", entries)
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
