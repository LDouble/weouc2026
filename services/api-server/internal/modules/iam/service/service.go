package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/liangluo/weouc2026/services/api-server/internal/modules/iam/config"
	"github.com/liangluo/weouc2026/services/api-server/internal/modules/iam/repo"
	iamtypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/iam/types"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/auth"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/httpx"
	"github.com/liangluo/weouc2026/services/api-server/internal/providers/academic_provider"
	"github.com/liangluo/weouc2026/services/api-server/internal/providers/wechat_provider"
)

var errInvalidToken = errors.New("invalid token")

type Service struct {
	config   config.ModuleConfig
	users    repo.UserRepository
	sessions repo.SessionRepository
	captchas repo.CaptchaRepository
	wechat   wechat_provider.Provider
	academic academic_provider.Provider
}

func New(
	cfg config.ModuleConfig,
	users repo.UserRepository,
	sessions repo.SessionRepository,
	captchas repo.CaptchaRepository,
	wechat wechat_provider.Provider,
	academic academic_provider.Provider,
) *Service {
	return &Service{
		config:   cfg,
		users:    users,
		sessions: sessions,
		captchas: captchas,
		wechat:   wechat,
		academic: academic,
	}
}

func (s *Service) LoginWithWeChat(ctx context.Context, request iamtypes.WeChatLoginRequest) (iamtypes.WeChatLoginResponse, error) {
	if strings.TrimSpace(request.Code) == "" || strings.TrimSpace(request.AppID) == "" {
		return iamtypes.WeChatLoginResponse{}, httpx.BadRequest("code 和 app_id 为必填项", nil)
	}

	identity, err := s.wechat.ExchangeCode(ctx, request.Code, request.AppID)
	if err != nil {
		return iamtypes.WeChatLoginResponse{}, httpx.Internal("微信登录失败", err)
	}

	user, exists := s.users.FindByOpenID(ctx, identity.OpenID)
	if !exists {
		now := time.Now().UTC()
		user = iamtypes.User{
			ID:        s.users.NextID(),
			OpenID:    identity.OpenID,
			Nickname:  identity.Nickname,
			AvatarURL: identity.AvatarURL,
			Roles:     append([]string(nil), s.config.DefaultRoles...),
			CreatedAt: now,
			UpdatedAt: now,
		}
	} else {
		user.Nickname = identity.Nickname
		user.AvatarURL = identity.AvatarURL
		user.UpdatedAt = time.Now().UTC()
	}

	user = s.users.Save(ctx, user)

	token, err := newRandomToken()
	if err != nil {
		return iamtypes.WeChatLoginResponse{}, httpx.Internal("生成登录凭证失败", err)
	}

	s.sessions.Save(ctx, iamtypes.Session{
		Token:     token,
		UserID:    user.ID,
		CreatedAt: time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(s.config.AccessTokenTTL),
	})

	return iamtypes.WeChatLoginResponse{
		Token:  token,
		OpenID: user.OpenID,
		UserInfo: iamtypes.WeChatUserInfo{
			UserID:    user.ID,
			Nickname:  user.Nickname,
			AvatarURL: user.AvatarURL,
		},
	}, nil
}

func (s *Service) ResolveToken(ctx context.Context, token string) (auth.Principal, error) {
	session, exists := s.sessions.Find(ctx, token)
	if !exists {
		return auth.Principal{}, errInvalidToken
	}
	if session.ExpiresAt.Before(time.Now().UTC()) {
		s.sessions.Delete(ctx, token)
		return auth.Principal{}, errInvalidToken
	}

	user, exists := s.users.FindByID(ctx, session.UserID)
	if !exists {
		return auth.Principal{}, errInvalidToken
	}

	permissions := append([]string(nil), user.Permissions...)
	academicBound := user.StudentProfile != nil && user.StudentProfile.IsBound
	if academicBound && !contains(permissions, s.config.StudentPermission) {
		permissions = append(permissions, s.config.StudentPermission)
	}

	return auth.Principal{
		Authenticated: true,
		UserID:        user.ID,
		DisplayName:   user.Nickname,
		Roles:         append([]string(nil), user.Roles...),
		Permissions:   permissions,
		AcademicBound: academicBound,
	}, nil
}

func (s *Service) GetStudentProfile(ctx context.Context, principal auth.Principal) (iamtypes.StudentProfile, error) {
	user, exists := s.users.FindByID(ctx, principal.UserID)
	if !exists {
		return iamtypes.StudentProfile{}, httpx.Unauthorized("需要登录后访问")
	}
	if user.StudentProfile == nil || !user.StudentProfile.IsBound {
		return iamtypes.StudentProfile{}, httpx.NotFound("当前账号尚未完成教务绑定", nil)
	}

	profile := *user.StudentProfile
	return profile, nil
}

func (s *Service) SendCaptcha(ctx context.Context, _ auth.Principal, studentID string) error {
	if strings.TrimSpace(studentID) == "" {
		return httpx.BadRequest("学号不能为空", nil)
	}

	code, err := s.academic.GenerateCaptcha(ctx, studentID)
	if err != nil {
		return httpx.Internal("验证码发送失败", err)
	}

	now := time.Now().UTC()
	s.captchas.Save(ctx, iamtypes.CaptchaTicket{
		StudentID: studentID,
		Code:      code,
		SentAt:    now,
		ExpiresAt: now.Add(s.config.CaptchaTTL),
	})

	return nil
}

func (s *Service) BindStudent(ctx context.Context, principal auth.Principal, request iamtypes.BindStudentRequest) (iamtypes.StudentProfile, error) {
	studentID := strings.TrimSpace(request.StudentID)
	password := strings.TrimSpace(request.Password)
	captcha := strings.TrimSpace(request.Captcha)
	if studentID == "" || password == "" || captcha == "" {
		return iamtypes.StudentProfile{}, httpx.BadRequest("student_id、password、captcha 为必填项", nil)
	}

	ticket, exists := s.captchas.Find(ctx, studentID)
	if !exists || ticket.ExpiresAt.Before(time.Now().UTC()) {
		return iamtypes.StudentProfile{}, httpx.BadRequest("验证码不存在或已过期", nil)
	}
	if ticket.Code != captcha {
		return iamtypes.StudentProfile{}, httpx.BadRequest("验证码错误", nil)
	}

	snapshot, err := s.academic.LoadStudentSnapshot(ctx, studentID, password)
	if err != nil {
		return iamtypes.StudentProfile{}, httpx.BadRequest("教务账号校验失败", nil)
	}

	profile, err := s.updateStudentProfile(ctx, principal.UserID, func(user *iamtypes.User) error {
		now := time.Now().UTC()
		user.StudentProfile = &iamtypes.StudentProfile{
			Name:      snapshot.Name,
			AvatarURL: user.AvatarURL,
			StudentID: studentID,
			Major:     snapshot.Major,
			College:   snapshot.College,
			Grade:     snapshot.Grade,
			IsBound:   true,
			UpdatedAt: now,
		}
		return nil
	})
	if err != nil {
		return iamtypes.StudentProfile{}, err
	}

	s.captchas.Delete(ctx, studentID)
	return profile, nil
}

func (s *Service) UnbindStudent(ctx context.Context, principal auth.Principal) error {
	_, err := s.users.Update(ctx, principal.UserID, func(user *iamtypes.User) error {
		user.StudentProfile = nil
		user.UpdatedAt = time.Now().UTC()
		return nil
	})
	if err != nil {
		return httpx.Internal("解绑失败", err)
	}

	return nil
}

func (s *Service) updateStudentProfile(
	ctx context.Context,
	userID string,
	mutate func(user *iamtypes.User) error,
) (iamtypes.StudentProfile, error) {
	var next iamtypes.StudentProfile
	_, err := s.users.Update(ctx, userID, func(user *iamtypes.User) error {
		if err := mutate(user); err != nil {
			return err
		}
		if user.StudentProfile == nil {
			return errors.New("student profile update resulted nil")
		}
		user.UpdatedAt = time.Now().UTC()
		next = *user.StudentProfile
		return nil
	})
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return iamtypes.StudentProfile{}, httpx.Unauthorized("需要登录后访问")
		}
		return iamtypes.StudentProfile{}, httpx.Internal("资料更新失败", err)
	}

	return next, nil
}

func newRandomToken() (string, error) {
	buffer := make([]byte, 32)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}

	return hex.EncodeToString(buffer), nil
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}

	return false
}
