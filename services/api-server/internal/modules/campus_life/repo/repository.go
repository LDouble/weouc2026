package repo

import (
	"context"
	"errors"

	cltypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/campus_life/types"
)

var ErrNotFound = errors.New("item not found")

type Repository interface {
	Save(ctx context.Context, item cltypes.CommunityContent) (cltypes.CommunityContent, error)
	GetByID(ctx context.Context, id string) (cltypes.CommunityContent, error)
	Update(ctx context.Context, id string, mutate func(*cltypes.CommunityContent) error) (cltypes.CommunityContent, error)
	ListByType(ctx context.Context, contentType string, filter cltypes.ContentFilter) ([]cltypes.CommunityContent, int64, error)
	ListForFeed(ctx context.Context, filter cltypes.FeedFilter) ([]cltypes.CommunityContent, int64, error)
	WriteTransitionLog(ctx context.Context, log cltypes.StateTransitionLog) error
}
