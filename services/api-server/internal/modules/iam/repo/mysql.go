package repo

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	iamtypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/iam/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MySQLUserRepository struct {
	db *gorm.DB
}

type mysqlUserModel struct {
	ID                 string    `gorm:"column:id;type:varchar(64);primaryKey"`
	OpenID             string    `gorm:"column:open_id;type:varchar(128);index:idx_iam_users_open_id,unique"`
	Username           string    `gorm:"column:username;type:varchar(128);index:idx_iam_users_username,unique"`
	PasswordHash       string    `gorm:"column:password_hash;type:text"`
	Nickname           string    `gorm:"column:nickname;type:varchar(128)"`
	AvatarURL          string    `gorm:"column:avatar_url;type:text"`
	RolesJSON          string    `gorm:"column:roles;type:json"`
	PermissionsJSON    string    `gorm:"column:permissions;type:json"`
	StudentProfileJSON string    `gorm:"column:student_profile;type:json"`
	CreatedAtRaw       time.Time `gorm:"column:created_at;not null"`
	UpdatedAtRaw       time.Time `gorm:"column:updated_at;not null"`
}

func (mysqlUserModel) TableName() string {
	return "iam_users"
}

func AutoMigrateMySQL(ctx context.Context, db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("mysql user migration db is nil")
	}

	return db.WithContext(ctx).AutoMigrate(&mysqlUserModel{})
}

func NewMySQLUserRepository(db *gorm.DB) *MySQLUserRepository {
	return &MySQLUserRepository{db: db}
}

func (r *MySQLUserRepository) FindByID(ctx context.Context, userID string) (iamtypes.User, error) {
	return r.findOne(ctx, "id = ?", userID)
}

func (r *MySQLUserRepository) FindByOpenID(ctx context.Context, openID string) (iamtypes.User, error) {
	return r.findOne(ctx, "open_id = ?", openID)
}

func (r *MySQLUserRepository) FindByUsername(ctx context.Context, username string) (iamtypes.User, error) {
	return r.findOne(ctx, "username = ?", username)
}

func (r *MySQLUserRepository) Save(ctx context.Context, user iamtypes.User) (iamtypes.User, error) {
	if r.db == nil {
		return iamtypes.User{}, fmt.Errorf("mysql user repository db is nil")
	}

	model, err := newMySQLUserModel(user)
	if err != nil {
		return iamtypes.User{}, err
	}

	if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
		return iamtypes.User{}, fmt.Errorf("save user failed: %w", err)
	}

	return model.toDomain()
}

func (r *MySQLUserRepository) Update(ctx context.Context, userID string, mutate func(*iamtypes.User) error) (iamtypes.User, error) {
	if r.db == nil {
		return iamtypes.User{}, fmt.Errorf("mysql user repository db is nil")
	}

	var updated iamtypes.User
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var model mysqlUserModel
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", userID).Take(&model).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrUserNotFound
			}
			return fmt.Errorf("load user for update failed: %w", err)
		}

		current, err := model.toDomain()
		if err != nil {
			return err
		}
		if err := mutate(&current); err != nil {
			return err
		}

		nextModel, err := newMySQLUserModel(current)
		if err != nil {
			return err
		}
		if err := tx.Save(&nextModel).Error; err != nil {
			return fmt.Errorf("update user failed: %w", err)
		}

		updated, err = nextModel.toDomain()
		return err
	})
	if err != nil {
		return iamtypes.User{}, err
	}

	return updated, nil
}

func (r *MySQLUserRepository) NextID() string {
	buffer := make([]byte, 8)
	if _, err := rand.Read(buffer); err != nil {
		return "user-fallback"
	}

	return "user-" + hex.EncodeToString(buffer)
}

func (r *MySQLUserRepository) findOne(ctx context.Context, query string, args ...any) (iamtypes.User, error) {
	if r.db == nil {
		return iamtypes.User{}, fmt.Errorf("mysql user repository db is nil")
	}

	var model mysqlUserModel
	if err := r.db.WithContext(ctx).Where(query, args...).Take(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return iamtypes.User{}, ErrUserNotFound
		}
		return iamtypes.User{}, fmt.Errorf("find user failed: %w", err)
	}

	return model.toDomain()
}

func newMySQLUserModel(user iamtypes.User) (mysqlUserModel, error) {
	roles, err := json.Marshal(user.Roles)
	if err != nil {
		return mysqlUserModel{}, fmt.Errorf("marshal roles failed: %w", err)
	}
	permissions, err := json.Marshal(user.Permissions)
	if err != nil {
		return mysqlUserModel{}, fmt.Errorf("marshal permissions failed: %w", err)
	}

	model := mysqlUserModel{
		ID:              user.ID,
		OpenID:          user.OpenID,
		Username:        user.Username,
		PasswordHash:    user.PasswordHash,
		Nickname:        user.Nickname,
		AvatarURL:       user.AvatarURL,
		RolesJSON:       string(roles),
		PermissionsJSON: string(permissions),
		CreatedAtRaw:    user.CreatedAt,
		UpdatedAtRaw:    user.UpdatedAt,
	}

	if user.StudentProfile != nil {
		profile, err := json.Marshal(user.StudentProfile)
		if err != nil {
			return mysqlUserModel{}, fmt.Errorf("marshal student profile failed: %w", err)
		}
		model.StudentProfileJSON = string(profile)
	}

	return model, nil
}

func (m mysqlUserModel) toDomain() (iamtypes.User, error) {
	user := iamtypes.User{
		ID:           m.ID,
		OpenID:       m.OpenID,
		Username:     m.Username,
		PasswordHash: m.PasswordHash,
		Nickname:     m.Nickname,
		AvatarURL:    m.AvatarURL,
		CreatedAt:    m.CreatedAtRaw,
		UpdatedAt:    m.UpdatedAtRaw,
	}

	if strings.TrimSpace(m.RolesJSON) != "" {
		if err := json.Unmarshal([]byte(m.RolesJSON), &user.Roles); err != nil {
			return iamtypes.User{}, fmt.Errorf("unmarshal roles failed: %w", err)
		}
	}
	if strings.TrimSpace(m.PermissionsJSON) != "" {
		if err := json.Unmarshal([]byte(m.PermissionsJSON), &user.Permissions); err != nil {
			return iamtypes.User{}, fmt.Errorf("unmarshal permissions failed: %w", err)
		}
	}
	if strings.TrimSpace(m.StudentProfileJSON) != "" && m.StudentProfileJSON != "null" {
		var profile iamtypes.StudentProfile
		if err := json.Unmarshal([]byte(m.StudentProfileJSON), &profile); err != nil {
			return iamtypes.User{}, fmt.Errorf("unmarshal student profile failed: %w", err)
		}
		user.StudentProfile = &profile
	}

	return user, nil
}
