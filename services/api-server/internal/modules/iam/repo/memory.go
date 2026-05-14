package repo

import (
	"context"
	"errors"

	iamtypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/iam/types"
)

var ErrUserNotFound = errors.New("user not found")
var ErrSessionNotFound = errors.New("session not found")
var ErrCaptchaNotFound = errors.New("captcha not found")

type UserRepository interface {
	FindByID(ctx context.Context, userID string) (iamtypes.User, error)
	FindByOpenID(ctx context.Context, openID string) (iamtypes.User, error)
	FindByUsername(ctx context.Context, username string) (iamtypes.User, error)
	Save(ctx context.Context, user iamtypes.User) (iamtypes.User, error)
	Update(ctx context.Context, userID string, mutate func(*iamtypes.User) error) (iamtypes.User, error)
	NextID() string
}

type SessionRepository interface {
	Save(ctx context.Context, session iamtypes.Session) (iamtypes.Session, error)
	Find(ctx context.Context, token string) (iamtypes.Session, error)
	Delete(ctx context.Context, token string) error
}

type CaptchaRepository interface {
	Save(ctx context.Context, ticket iamtypes.CaptchaTicket) error
	Find(ctx context.Context, studentID string) (iamtypes.CaptchaTicket, error)
	Delete(ctx context.Context, studentID string) error
}

func cloneUser(user iamtypes.User) iamtypes.User {
	cloned := user
	cloned.Roles = append([]string(nil), user.Roles...)
	cloned.Permissions = append([]string(nil), user.Permissions...)
	if user.StudentProfile != nil {
		profile := *user.StudentProfile
		cloned.StudentProfile = &profile
	}

	return cloned
}
