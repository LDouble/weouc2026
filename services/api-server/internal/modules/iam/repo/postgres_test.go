package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	iamtypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/iam/types"
)

func TestPostgresUserRepositoryFindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("new sqlmock failed: %v", err)
	}
	defer db.Close()

	profileJSON, err := json.Marshal(iamtypes.StudentProfile{
		Name:      "同学0001",
		AvatarURL: "https://example.com/avatar.png",
		StudentID: "20260001",
		Major:     "软件工程",
		College:   "信息工程学院",
		Grade:     "2024级",
		IsBound:   true,
		UpdatedAt: time.Date(2026, 5, 10, 9, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("marshal profile failed: %v", err)
	}

	rows := sqlmock.NewRows([]string{
		"id", "open_id", "nickname", "avatar_url", "roles", "permissions", "student_profile", "created_at", "updated_at",
	}).AddRow(
		"user-001",
		"openid-001",
		"海大同学",
		"https://example.com/avatar.png",
		`["student"]`,
		`["contact:view"]`,
		profileJSON,
		time.Date(2026, 5, 10, 8, 0, 0, 0, time.UTC),
		time.Date(2026, 5, 10, 8, 30, 0, 0, time.UTC),
	)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT ` + userColumns + ` FROM iam_users WHERE id = $1 LIMIT 1`)).
		WithArgs("user-001").
		WillReturnRows(rows)

	repository := NewPostgresUserRepository(db)
	user, err := repository.FindByID(context.Background(), "user-001")
	if err != nil {
		t.Fatalf("FindByID returned error: %v", err)
	}
	if user.OpenID != "openid-001" || user.StudentProfile == nil || user.StudentProfile.StudentID != "20260001" {
		t.Fatalf("unexpected user payload: %+v", user)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestPostgresUserRepositoryUpdateReturnsNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("new sqlmock failed: %v", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT ` + userColumns + ` FROM iam_users WHERE id = $1 FOR UPDATE`)).
		WithArgs("missing-user").
		WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()

	repository := NewPostgresUserRepository(db)
	_, err = repository.Update(context.Background(), "missing-user", func(user *iamtypes.User) error {
		return nil
	})
	if err != ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}
