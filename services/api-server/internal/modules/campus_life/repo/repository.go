package repo

import (
	"context"
	"errors"

	cltypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/campus_life/types"
)

var ErrNotFound = errors.New("item not found")

type Repository interface {
	ListMarkets(ctx context.Context) ([]cltypes.MarketItem, error)
	GetMarket(ctx context.Context, id string) (cltypes.MarketItem, error)
	SaveMarket(ctx context.Context, item cltypes.MarketItem) (cltypes.MarketItem, error)
	UpdateMarket(ctx context.Context, id string, mutate func(*cltypes.MarketItem) error) (cltypes.MarketItem, error)

	ListErrands(ctx context.Context) ([]cltypes.ErrandItem, error)
	GetErrand(ctx context.Context, id string) (cltypes.ErrandItem, error)
	SaveErrand(ctx context.Context, item cltypes.ErrandItem) (cltypes.ErrandItem, error)
	UpdateErrand(ctx context.Context, id string, mutate func(*cltypes.ErrandItem) error) (cltypes.ErrandItem, error)

	ListResources(ctx context.Context) ([]cltypes.ResourceItem, error)
	GetResource(ctx context.Context, id string) (cltypes.ResourceItem, error)
	SaveResource(ctx context.Context, item cltypes.ResourceItem) (cltypes.ResourceItem, error)
	UpdateResource(ctx context.Context, id string, mutate func(*cltypes.ResourceItem) error) (cltypes.ResourceItem, error)

	ListLostFound(ctx context.Context) ([]cltypes.LostFoundItem, error)
	GetLostFound(ctx context.Context, id string) (cltypes.LostFoundItem, error)
	SaveLostFound(ctx context.Context, item cltypes.LostFoundItem) (cltypes.LostFoundItem, error)
	UpdateLostFound(ctx context.Context, id string, mutate func(*cltypes.LostFoundItem) error) (cltypes.LostFoundItem, error)

	ListCarpools(ctx context.Context) ([]cltypes.CarpoolItem, error)
	GetCarpool(ctx context.Context, id string) (cltypes.CarpoolItem, error)
	SaveCarpool(ctx context.Context, item cltypes.CarpoolItem) (cltypes.CarpoolItem, error)
	UpdateCarpool(ctx context.Context, id string, mutate func(*cltypes.CarpoolItem) error) (cltypes.CarpoolItem, error)

	NextID(ctx context.Context, prefix string) (string, error)
}
