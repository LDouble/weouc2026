package migrate

import (
	"context"
	"database/sql"
	"testing"
	"testing/fstest"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRunFSAppliesUnseenMigrations(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("new sqlmock failed: %v", err)
	}
	defer db.Close()

	migrations := fstest.MapFS{
		"20260510_0001_test.sql": &fstest.MapFile{
			Data: []byte("CREATE TABLE test_users (id TEXT PRIMARY KEY);"),
		},
	}

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS schema_migrations").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SELECT 1 FROM schema_migrations WHERE version = \\$1").
		WithArgs("20260510_0001_test.sql").
		WillReturnError(sql.ErrNoRows)
	mock.ExpectBegin()
	mock.ExpectExec("CREATE TABLE test_users").
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("INSERT INTO schema_migrations").
		WithArgs("20260510_0001_test.sql", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err := runFS(context.Background(), db, migrations); err != nil {
		t.Fatalf("runFS returned error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestRunFSSkipsAppliedMigrations(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("new sqlmock failed: %v", err)
	}
	defer db.Close()

	migrations := fstest.MapFS{
		"20260510_0001_test.sql": &fstest.MapFile{
			Data: []byte("CREATE TABLE test_users (id TEXT PRIMARY KEY);"),
		},
	}

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS schema_migrations").WillReturnResult(sqlmock.NewResult(0, 0))
	rows := sqlmock.NewRows([]string{"1"}).AddRow(1)
	mock.ExpectQuery("SELECT 1 FROM schema_migrations WHERE version = \\$1").
		WithArgs("20260510_0001_test.sql").
		WillReturnRows(rows)

	if err := runFS(context.Background(), db, migrations); err != nil {
		t.Fatalf("runFS returned error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}
