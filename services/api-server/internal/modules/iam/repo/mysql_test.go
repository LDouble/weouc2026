package repo

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	iamtypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/iam/types"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestMySQLUserRepositoryFindByID(t *testing.T) {
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

	gormDB, err := gorm.Open(gormmysql.New(gormmysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("open gorm failed: %v", err)
	}

	rows := sqlmock.NewRows([]string{
		"id", "open_id", "username", "password_hash", "nickname", "avatar_url", "roles", "permissions", "student_profile", "created_at", "updated_at",
	}).AddRow(
		"user-001",
		"openid-001",
		"student001",
		"",
		"海大同学",
		"https://example.com/avatar.png",
		`["student"]`,
		`["contact:view"]`,
		string(profileJSON),
		time.Date(2026, 5, 10, 8, 0, 0, 0, time.UTC),
		time.Date(2026, 5, 10, 8, 30, 0, 0, time.UTC),
	)

	mock.ExpectQuery("SELECT \\* FROM `iam_users` WHERE id = \\? LIMIT \\?").
		WithArgs("user-001", 1).
		WillReturnRows(rows)

	repository := NewMySQLUserRepository(gormDB)
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

func TestMySQLUserRepositoryUpdateReturnsNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("new sqlmock failed: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(gormmysql.New(gormmysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("open gorm failed: %v", err)
	}

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT \\* FROM `iam_users` WHERE id = \\? LIMIT \\? FOR UPDATE").
		WithArgs("missing-user", 1).
		WillReturnError(gorm.ErrRecordNotFound)
	mock.ExpectRollback()

	repository := NewMySQLUserRepository(gormDB)
	_, err = repository.Update(context.Background(), "missing-user", func(user *iamtypes.User) error {
		return nil
	})
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}
