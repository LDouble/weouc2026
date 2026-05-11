package repo

import (
	"context"
	"errors"

	portaltypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/portal/types"
)

var ErrNotFound = errors.New("portal content not found")

type Repository interface {
	ListBanners(ctx context.Context) ([]portaltypes.BannerItem, error)
	ListNotices(ctx context.Context) ([]portaltypes.NoticeItem, error)
	GetNotice(ctx context.Context, id string) (portaltypes.NoticeItem, error)
	SaveNotice(ctx context.Context, item portaltypes.NoticeItem) (portaltypes.NoticeItem, error)
	NextID(ctx context.Context, prefix string) (string, error)
}
