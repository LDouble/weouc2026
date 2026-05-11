package repo

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	iamtypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/iam/types"
)

const userColumns = `
id,
open_id,
username,
password_hash,
nickname,
avatar_url,
roles,
permissions,
student_profile,
created_at,
updated_at`

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, userID string) (iamtypes.User, error) {
	if r.db == nil {
		return iamtypes.User{}, fmt.Errorf("postgres user repository db is nil")
	}

	row := r.db.QueryRowContext(ctx, `SELECT `+userColumns+` FROM iam_users WHERE id = $1 LIMIT 1`, userID)
	user, err := scanUser(row)
	if errors.Is(err, sql.ErrNoRows) {
		return iamtypes.User{}, ErrUserNotFound
	}
	if err != nil {
		return iamtypes.User{}, fmt.Errorf("find user by id failed: %w", err)
	}

	return user, nil
}

func (r *PostgresUserRepository) FindByOpenID(ctx context.Context, openID string) (iamtypes.User, error) {
	if r.db == nil {
		return iamtypes.User{}, fmt.Errorf("postgres user repository db is nil")
	}

	row := r.db.QueryRowContext(ctx, `SELECT `+userColumns+` FROM iam_users WHERE open_id = $1 LIMIT 1`, openID)
	user, err := scanUser(row)
	if errors.Is(err, sql.ErrNoRows) {
		return iamtypes.User{}, ErrUserNotFound
	}
	if err != nil {
		return iamtypes.User{}, fmt.Errorf("find user by open id failed: %w", err)
	}

	return user, nil
}

func (r *PostgresUserRepository) FindByUsername(ctx context.Context, username string) (iamtypes.User, error) {
	if r.db == nil {
		return iamtypes.User{}, fmt.Errorf("postgres user repository db is nil")
	}

	row := r.db.QueryRowContext(ctx, `SELECT `+userColumns+` FROM iam_users WHERE username = $1 LIMIT 1`, username)
	user, err := scanUser(row)
	if errors.Is(err, sql.ErrNoRows) {
		return iamtypes.User{}, ErrUserNotFound
	}
	if err != nil {
		return iamtypes.User{}, fmt.Errorf("find user by username failed: %w", err)
	}

	return user, nil
}

func (r *PostgresUserRepository) Save(ctx context.Context, user iamtypes.User) (iamtypes.User, error) {
	if r.db == nil {
		return iamtypes.User{}, fmt.Errorf("postgres user repository db is nil")
	}

	saved, err := saveUser(ctx, r.db, user)
	if err != nil {
		return iamtypes.User{}, fmt.Errorf("save user failed: %w", err)
	}

	return saved, nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, userID string, mutate func(*iamtypes.User) error) (iamtypes.User, error) {
	if r.db == nil {
		return iamtypes.User{}, fmt.Errorf("postgres user repository db is nil")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return iamtypes.User{}, fmt.Errorf("begin user update tx failed: %w", err)
	}
	defer tx.Rollback()

	current, err := scanUser(tx.QueryRowContext(ctx, `SELECT `+userColumns+` FROM iam_users WHERE id = $1 FOR UPDATE`, userID))
	if errors.Is(err, sql.ErrNoRows) {
		return iamtypes.User{}, ErrUserNotFound
	}
	if err != nil {
		return iamtypes.User{}, fmt.Errorf("load user for update failed: %w", err)
	}
	if err := mutate(&current); err != nil {
		return iamtypes.User{}, err
	}

	updated, err := saveUser(ctx, tx, current)
	if err != nil {
		return iamtypes.User{}, fmt.Errorf("update user failed: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return iamtypes.User{}, fmt.Errorf("commit user update failed: %w", err)
	}

	return updated, nil
}

func (r *PostgresUserRepository) NextID() string {
	buffer := make([]byte, 8)
	if _, err := rand.Read(buffer); err != nil {
		return "user-fallback"
	}

	return "user-" + hex.EncodeToString(buffer)
}

type queryRower interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func saveUser(ctx context.Context, querier queryRower, user iamtypes.User) (iamtypes.User, error) {
	roles, err := json.Marshal(user.Roles)
	if err != nil {
		return iamtypes.User{}, fmt.Errorf("marshal roles failed: %w", err)
	}
	permissions, err := json.Marshal(user.Permissions)
	if err != nil {
		return iamtypes.User{}, fmt.Errorf("marshal permissions failed: %w", err)
	}

	var studentProfile any
	if user.StudentProfile != nil {
		studentProfile, err = json.Marshal(user.StudentProfile)
		if err != nil {
			return iamtypes.User{}, fmt.Errorf("marshal student profile failed: %w", err)
		}
	}

	row := querier.QueryRowContext(
		ctx,
		`INSERT INTO iam_users (
			id,
			open_id,
			username,
			password_hash,
			nickname,
			avatar_url,
			roles,
			permissions,
			student_profile,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7::jsonb, $8::jsonb, $9::jsonb, $10, $11)
		ON CONFLICT (id) DO UPDATE SET
			open_id = EXCLUDED.open_id,
			username = EXCLUDED.username,
			password_hash = EXCLUDED.password_hash,
			nickname = EXCLUDED.nickname,
			avatar_url = EXCLUDED.avatar_url,
			roles = EXCLUDED.roles,
			permissions = EXCLUDED.permissions,
			student_profile = EXCLUDED.student_profile,
			updated_at = EXCLUDED.updated_at
		RETURNING `+userColumns,
		user.ID,
		user.OpenID,
		user.Username,
		user.PasswordHash,
		user.Nickname,
		user.AvatarURL,
		string(roles),
		string(permissions),
		studentProfile,
		user.CreatedAt,
		user.UpdatedAt,
	)

	saved, err := scanUser(row)
	if err != nil {
		return iamtypes.User{}, err
	}

	return saved, nil
}

type userScanner interface {
	Scan(dest ...any) error
}

func scanUser(scanner userScanner) (iamtypes.User, error) {
	var user iamtypes.User
	var rolesRaw []byte
	var permissionsRaw []byte
	var studentProfileRaw []byte

	if err := scanner.Scan(
		&user.ID,
		&user.OpenID,
		&user.Username,
		&user.PasswordHash,
		&user.Nickname,
		&user.AvatarURL,
		&rolesRaw,
		&permissionsRaw,
		&studentProfileRaw,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return iamtypes.User{}, err
	}

	if len(rolesRaw) > 0 {
		if err := json.Unmarshal(rolesRaw, &user.Roles); err != nil {
			return iamtypes.User{}, fmt.Errorf("unmarshal roles failed: %w", err)
		}
	}
	if len(permissionsRaw) > 0 {
		if err := json.Unmarshal(permissionsRaw, &user.Permissions); err != nil {
			return iamtypes.User{}, fmt.Errorf("unmarshal permissions failed: %w", err)
		}
	}
	if len(studentProfileRaw) > 0 {
		var profile iamtypes.StudentProfile
		if err := json.Unmarshal(studentProfileRaw, &profile); err != nil {
			return iamtypes.User{}, fmt.Errorf("unmarshal student profile failed: %w", err)
		}
		user.StudentProfile = &profile
	}

	return user, nil
}
