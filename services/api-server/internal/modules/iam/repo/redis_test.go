package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	miniredis "github.com/alicebob/miniredis/v2"
	iamtypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/iam/types"
	"github.com/redis/go-redis/v9"
)

func TestRedisSessionRepositorySaveFindDelete(t *testing.T) {
	server, err := miniredis.Run()
	if err != nil {
		t.Fatalf("start miniredis failed: %v", err)
	}
	defer server.Close()

	client := redis.NewClient(&redis.Options{Addr: server.Addr()})
	defer client.Close()

	repository := NewRedisSessionRepository(client)
	session := iamtypes.Session{
		Token:     "token-001",
		UserID:    "user-001",
		CreatedAt: time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(10 * time.Minute),
	}

	if _, err := repository.Save(context.Background(), session); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}
	loaded, err := repository.Find(context.Background(), "token-001")
	if err != nil {
		t.Fatalf("Find returned error: %v", err)
	}
	if loaded.UserID != "user-001" {
		t.Fatalf("unexpected session payload: %+v", loaded)
	}
	if err := repository.Delete(context.Background(), "token-001"); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	if _, err := repository.Find(context.Background(), "token-001"); !errors.Is(err, ErrSessionNotFound) {
		t.Fatalf("expected ErrSessionNotFound after delete, got %v", err)
	}
}

func TestRedisCaptchaRepositorySaveFindDelete(t *testing.T) {
	server, err := miniredis.Run()
	if err != nil {
		t.Fatalf("start miniredis failed: %v", err)
	}
	defer server.Close()

	client := redis.NewClient(&redis.Options{Addr: server.Addr()})
	defer client.Close()

	repository := NewRedisCaptchaRepository(client)
	ticket := iamtypes.CaptchaTicket{
		StudentID: "20260001",
		Code:      "123456",
		SentAt:    time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(5 * time.Minute),
	}

	if err := repository.Save(context.Background(), ticket); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}
	loaded, err := repository.Find(context.Background(), "20260001")
	if err != nil {
		t.Fatalf("Find returned error: %v", err)
	}
	if loaded.Code != "123456" {
		t.Fatalf("unexpected captcha payload: %+v", loaded)
	}
	if err := repository.Delete(context.Background(), "20260001"); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	if _, err := repository.Find(context.Background(), "20260001"); !errors.Is(err, ErrCaptchaNotFound) {
		t.Fatalf("expected ErrCaptchaNotFound after delete, got %v", err)
	}
}
