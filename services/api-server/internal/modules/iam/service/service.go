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
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/audit"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/auth"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/httpx"
	"github.com/liangluo/weouc2026/services/api-server/internal/providers/academic_provider"
	"github.com/liangluo/weouc2026/services/api-server/internal/providers/wechat_provider"
)

type Service struct {
	config   config.ModuleConfig
	users    repo.UserRepository
	sessions repo.SessionRepository
	captchas repo.CaptchaRepository
	wechat   wechat_provider.Provider
	academic academic_provider.Provider
	recorder audit.Recorder
}

func New(
	cfg config.ModuleConfig,
	users repo.UserRepository,
	sessions repo.SessionRepository,
	captchas repo.CaptchaRepository,
	wechat wechat_provider.Provider,
	academic academic_provider.Provider,
	recorder audit.Recorder,
) *Service {
	return &Service{
		config:   cfg,
		users:    users,
		sessions: sessions,
		captchas: captchas,
		wechat:   wechat,
		academic: academic,
		recorder: recorder,
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

	user, err := s.users.FindByOpenID(ctx, identity.OpenID)
	if errors.Is(err, repo.ErrUserNotFound) {
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
	} else if err != nil {
		return iamtypes.WeChatLoginResponse{}, httpx.Internal("加载用户资料失败", err)
	} else {
		user.Nickname = identity.Nickname
		user.AvatarURL = identity.AvatarURL
		user.UpdatedAt = time.Now().UTC()
	}

	user, err = s.users.Save(ctx, user)
	if err != nil {
		return iamtypes.WeChatLoginResponse{}, httpx.Internal("保存用户资料失败", err)
	}

	token, err := newRandomToken()
	if err != nil {
		return iamtypes.WeChatLoginResponse{}, httpx.Internal("生成登录凭证失败", err)
	}

	if _, err := s.sessions.Save(ctx, iamtypes.Session{
		Token:     token,
		UserID:    user.ID,
		CreatedAt: time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(s.config.AccessTokenTTL),
	}); err != nil {
		return iamtypes.WeChatLoginResponse{}, httpx.Internal("保存登录会话失败", err)
	}

	audit.RecordBestEffort(ctx, s.recorder, audit.Entry{
		ActorID:      user.ID,
		ActorName:    user.Nickname,
		Action:       "auth.login",
		ResourceType: "user",
		ResourceID:   user.ID,
		Message:      "微信登录成功",
		Details: map[string]any{
			"openid": user.OpenID,
		},
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
	session, err := s.sessions.Find(ctx, token)
	if errors.Is(err, repo.ErrSessionNotFound) {
		return auth.Principal{}, httpx.Unauthorized("登录状态已失效，请重新登录")
	}
	if err != nil {
		return auth.Principal{}, httpx.Internal("登录态解析失败", err)
	}
	if session.ExpiresAt.Before(time.Now().UTC()) {
		if err := s.sessions.Delete(ctx, token); err != nil {
			return auth.Principal{}, httpx.Internal("清理过期登录态失败", err)
		}
		return auth.Principal{}, httpx.Unauthorized("登录状态已失效，请重新登录")
	}

	user, err := s.users.FindByID(ctx, session.UserID)
	if errors.Is(err, repo.ErrUserNotFound) {
		return auth.Principal{}, httpx.Unauthorized("登录状态已失效，请重新登录")
	}
	if err != nil {
		return auth.Principal{}, httpx.Internal("加载当前用户失败", err)
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
	user, err := s.users.FindByID(ctx, principal.UserID)
	if errors.Is(err, repo.ErrUserNotFound) {
		return iamtypes.StudentProfile{}, httpx.Unauthorized("需要登录后访问")
	}
	if err != nil {
		return iamtypes.StudentProfile{}, httpx.Internal("加载当前资料失败", err)
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
	if err := s.captchas.Save(ctx, iamtypes.CaptchaTicket{
		StudentID: studentID,
		Code:      code,
		SentAt:    now,
		ExpiresAt: now.Add(s.config.CaptchaTTL),
	}); err != nil {
		return httpx.Internal("保存验证码失败", err)
	}

	return nil
}

func (s *Service) BindStudent(ctx context.Context, principal auth.Principal, request iamtypes.BindStudentRequest) (iamtypes.StudentProfile, error) {
	studentID := strings.TrimSpace(request.StudentID)
	password := strings.TrimSpace(request.Password)
	captcha := strings.TrimSpace(request.Captcha)
	if studentID == "" || password == "" || captcha == "" {
		return iamtypes.StudentProfile{}, httpx.BadRequest("student_id、password、captcha 为必填项", nil)
	}

	ticket, err := s.captchas.Find(ctx, studentID)
	if errors.Is(err, repo.ErrCaptchaNotFound) {
		return iamtypes.StudentProfile{}, httpx.BadRequest("验证码不存在或已过期", nil)
	}
	if err != nil {
		return iamtypes.StudentProfile{}, httpx.Internal("读取验证码失败", err)
	}
	if ticket.ExpiresAt.Before(time.Now().UTC()) {
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

	if err := s.captchas.Delete(ctx, studentID); err != nil {
		return iamtypes.StudentProfile{}, httpx.Internal("清理验证码失败", err)
	}

	audit.RecordBestEffort(ctx, s.recorder, audit.Entry{
		ActorID:      principal.UserID,
		ActorName:    principal.DisplayName,
		Action:       "auth.bind_student",
		ResourceType: "student_profile",
		ResourceID:   studentID,
		Message:      "教务绑定成功",
		Details: map[string]any{
			"college": profile.College,
			"major":   profile.Major,
			"grade":   profile.Grade,
		},
	})
	return profile, nil
}

func (s *Service) UnbindStudent(ctx context.Context, principal auth.Principal) error {
	var previousStudentID string
	_, err := s.users.Update(ctx, principal.UserID, func(user *iamtypes.User) error {
		if user.StudentProfile != nil {
			previousStudentID = user.StudentProfile.StudentID
		}
		user.StudentProfile = nil
		user.UpdatedAt = time.Now().UTC()
		return nil
	})
	if err != nil {
		return httpx.Internal("解绑失败", err)
	}

	audit.RecordBestEffort(ctx, s.recorder, audit.Entry{
		ActorID:      principal.UserID,
		ActorName:    principal.DisplayName,
		Action:       "auth.unbind_student",
		ResourceType: "student_profile",
		ResourceID:   previousStudentID,
		Message:      "教务解绑成功",
	})

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
