package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	iamtypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/iam/types"
	"github.com/redis/go-redis/v9"
)

const (
	sessionKeyPrefix = "iam:session:"
	captchaKeyPrefix = "iam:captcha:"
)

type RedisSessionRepository struct {
	client *redis.Client
}

type RedisCaptchaRepository struct {
	client *redis.Client
}

func NewRedisSessionRepository(client *redis.Client) *RedisSessionRepository {
	return &RedisSessionRepository{client: client}
}

func NewRedisCaptchaRepository(client *redis.Client) *RedisCaptchaRepository {
	return &RedisCaptchaRepository{client: client}
}

func (r *RedisSessionRepository) Save(ctx context.Context, session iamtypes.Session) (iamtypes.Session, error) {
	if r.client == nil {
		return iamtypes.Session{}, fmt.Errorf("redis session repository client is nil")
	}

	payload, err := json.Marshal(session)
	if err != nil {
		return iamtypes.Session{}, fmt.Errorf("marshal session failed: %w", err)
	}
	if err := r.client.Set(ctx, sessionRedisKey(session.Token), payload, ttlUntil(session.ExpiresAt)).Err(); err != nil {
		return iamtypes.Session{}, fmt.Errorf("save session failed: %w", err)
	}

	return session, nil
}

func (r *RedisSessionRepository) Find(ctx context.Context, token string) (iamtypes.Session, error) {
	if r.client == nil {
		return iamtypes.Session{}, fmt.Errorf("redis session repository client is nil")
	}

	payload, err := r.client.Get(ctx, sessionRedisKey(token)).Bytes()
	if err == redis.Nil {
		return iamtypes.Session{}, ErrSessionNotFound
	}
	if err != nil {
		return iamtypes.Session{}, fmt.Errorf("find session failed: %w", err)
	}

	var session iamtypes.Session
	if err := json.Unmarshal(payload, &session); err != nil {
		return iamtypes.Session{}, fmt.Errorf("unmarshal session failed: %w", err)
	}

	return session, nil
}

func (r *RedisSessionRepository) Delete(ctx context.Context, token string) error {
	if r.client == nil {
		return fmt.Errorf("redis session repository client is nil")
	}

	if err := r.client.Del(ctx, sessionRedisKey(token)).Err(); err != nil {
		return fmt.Errorf("delete session failed: %w", err)
	}

	return nil
}

func (r *RedisCaptchaRepository) Save(ctx context.Context, ticket iamtypes.CaptchaTicket) error {
	if r.client == nil {
		return fmt.Errorf("redis captcha repository client is nil")
	}

	payload, err := json.Marshal(ticket)
	if err != nil {
		return fmt.Errorf("marshal captcha failed: %w", err)
	}
	if err := r.client.Set(ctx, captchaRedisKey(ticket.StudentID), payload, ttlUntil(ticket.ExpiresAt)).Err(); err != nil {
		return fmt.Errorf("save captcha failed: %w", err)
	}

	return nil
}

func (r *RedisCaptchaRepository) Find(ctx context.Context, studentID string) (iamtypes.CaptchaTicket, error) {
	if r.client == nil {
		return iamtypes.CaptchaTicket{}, fmt.Errorf("redis captcha repository client is nil")
	}

	payload, err := r.client.Get(ctx, captchaRedisKey(studentID)).Bytes()
	if err == redis.Nil {
		return iamtypes.CaptchaTicket{}, ErrCaptchaNotFound
	}
	if err != nil {
		return iamtypes.CaptchaTicket{}, fmt.Errorf("find captcha failed: %w", err)
	}

	var ticket iamtypes.CaptchaTicket
	if err := json.Unmarshal(payload, &ticket); err != nil {
		return iamtypes.CaptchaTicket{}, fmt.Errorf("unmarshal captcha failed: %w", err)
	}

	return ticket, nil
}

func (r *RedisCaptchaRepository) Delete(ctx context.Context, studentID string) error {
	if r.client == nil {
		return fmt.Errorf("redis captcha repository client is nil")
	}

	if err := r.client.Del(ctx, captchaRedisKey(studentID)).Err(); err != nil {
		return fmt.Errorf("delete captcha failed: %w", err)
	}

	return nil
}

func sessionRedisKey(token string) string {
	return sessionKeyPrefix + token
}

func captchaRedisKey(studentID string) string {
	return captchaKeyPrefix + studentID
}

func ttlUntil(expiresAt time.Time) time.Duration {
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return time.Second
	}

	return ttl
}
