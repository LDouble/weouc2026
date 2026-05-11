package types

import "time"

type User struct {
	ID             string
	OpenID         string
	Username       string
	PasswordHash   string
	Nickname       string
	AvatarURL      string
	Roles          []string
	Permissions    []string
	StudentProfile *StudentProfile
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type StudentProfile struct {
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatar_url"`
	StudentID string    `json:"student_id"`
	Major     string    `json:"major"`
	College   string    `json:"college"`
	Grade     string    `json:"grade"`
	IsBound   bool      `json:"is_bound"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Session struct {
	Token     string
	UserID    string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type CaptchaTicket struct {
	StudentID string
	Code      string
	SentAt    time.Time
	ExpiresAt time.Time
}

type WeChatLoginRequest struct {
	Code  string `json:"code"`
	AppID string `json:"app_id"`
}

type WeChatLoginResponse struct {
	Token    string         `json:"token"`
	OpenID   string         `json:"openid"`
	UserInfo WeChatUserInfo `json:"userInfo"`
}

type WeChatUserInfo struct {
	UserID    string `json:"userId"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatarUrl"`
}

type SendCaptchaRequest struct {
	StudentID string `json:"sid"`
}

type BindStudentRequest struct {
	StudentID string `json:"student_id"`
	Password  string `json:"password"`
	Captcha   string `json:"captcha"`
}

type UpdateStudentRequest struct {
	IsBound *bool `json:"is_bound"`
}

type AdminLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AdminLoginResponse struct {
	Token    string `json:"token"`
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Roles    []string `json:"roles"`
}
