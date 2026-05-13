package repo

import (
	"context"
	"errors"

	portaltypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/portal/types"
)

var ErrNotFound = errors.New("portal content not found")

type BannerListQuery struct {
	Keyword string
}

type NoticeListQuery struct {
	Keyword string
}

type Repository interface {
	ListBanners(ctx context.Context, query BannerListQuery) ([]portaltypes.BannerItem, error)
	GetBanner(ctx context.Context, id string) (portaltypes.BannerItem, error)
	SaveBanner(ctx context.Context, item portaltypes.BannerItem) (portaltypes.BannerItem, error)
	DeleteBanner(ctx context.Context, id string) error
	ListNotices(ctx context.Context, query NoticeListQuery) ([]portaltypes.NoticeItem, error)
	GetNotice(ctx context.Context, id string) (portaltypes.NoticeItem, error)
	SaveNotice(ctx context.Context, item portaltypes.NoticeItem) (portaltypes.NoticeItem, error)
	DeleteNotice(ctx context.Context, id string) error
	NextID(ctx context.Context, prefix string) (string, error)
}
