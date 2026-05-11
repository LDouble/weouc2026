package repo

import (
	"context"
	"errors"

	notificationtypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/notification/types"
)

var ErrNotFound = errors.New("notification not found")

type Repository interface {
	ListMessages(ctx context.Context) ([]notificationtypes.MessageItem, error)
	GetMessage(ctx context.Context, id string) (notificationtypes.MessageItem, error)
	SaveMessage(ctx context.Context, item notificationtypes.MessageItem) (notificationtypes.MessageItem, error)
	UpdateMessage(ctx context.Context, id string, mutate func(*notificationtypes.MessageItem) error) (notificationtypes.MessageItem, error)
	NextID(ctx context.Context, prefix string) (string, error)
}
