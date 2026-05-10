package wechat_provider

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

type Identity struct {
	OpenID    string
	Nickname  string
	AvatarURL string
}

type Provider interface {
	ExchangeCode(ctx context.Context, code, appID string) (Identity, error)
}

type MockProvider struct{}

func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

func (p *MockProvider) ExchangeCode(_ context.Context, code, appID string) (Identity, error) {
	sum := sha1.Sum([]byte(code + ":" + appID))
	short := hex.EncodeToString(sum[:])[:10]

	return Identity{
		OpenID:    "mock-openid-" + short,
		Nickname:  "微信用户" + short[len(short)-4:],
		AvatarURL: fmt.Sprintf("https://example.com/avatar/%s.png", short),
	}, nil
}
