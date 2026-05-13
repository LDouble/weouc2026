package repo

import (
	"context"
	"errors"
	"time"

	cltypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/campus_life/types"
)

var ErrNotFound = errors.New("item not found")

type ContentVisibilityQuery struct {
	ReviewStatuses           []string
	IncludeAllReviewStatuses bool
	IncludeOwnerUserID       string
}

type MarketListQuery struct {
	Visibility      ContentVisibilityQuery
	Category        string
	Keyword         string
	PublisherUserID string
}

type ErrandListQuery struct {
	Visibility      ContentVisibilityQuery
	Category        string
	Keyword         string
	PublisherUserID string
	AcceptorUserID  string
}

type ResourceListQuery struct {
	Visibility      ContentVisibilityQuery
	Category        string
	Keyword         string
	PublisherUserID string
}

type LostFoundListQuery struct {
	Visibility      ContentVisibilityQuery
	Category        string
	Type            string
	Keyword         string
	PublisherUserID string
}

type CarpoolListQuery struct {
	Visibility      ContentVisibilityQuery
	Keyword         string
	PublisherUserID string
	TravelAtFrom    time.Time
	TravelAtTo      time.Time
}

type MeetupStateQuery struct {
	IncludeAllStatuses        bool
	IncludeCancelledForUserID string
}

type MeetupListQuery struct {
	Visibility        ContentVisibilityQuery
	Category          string
	Keyword           string
	PublisherUserID   string
	ParticipantUserID string
	State             MeetupStateQuery
}

type Repository interface {
	ListMarkets(ctx context.Context, query MarketListQuery) ([]cltypes.MarketItem, error)
	GetMarket(ctx context.Context, id string) (cltypes.MarketItem, error)
	SaveMarket(ctx context.Context, item cltypes.MarketItem) (cltypes.MarketItem, error)
	UpdateMarket(ctx context.Context, id string, mutate func(*cltypes.MarketItem) error) (cltypes.MarketItem, error)

	ListErrands(ctx context.Context, query ErrandListQuery) ([]cltypes.ErrandItem, error)
	GetErrand(ctx context.Context, id string) (cltypes.ErrandItem, error)
	SaveErrand(ctx context.Context, item cltypes.ErrandItem) (cltypes.ErrandItem, error)
	UpdateErrand(ctx context.Context, id string, mutate func(*cltypes.ErrandItem) error) (cltypes.ErrandItem, error)

	ListResources(ctx context.Context, query ResourceListQuery) ([]cltypes.ResourceItem, error)
	GetResource(ctx context.Context, id string) (cltypes.ResourceItem, error)
	SaveResource(ctx context.Context, item cltypes.ResourceItem) (cltypes.ResourceItem, error)
	UpdateResource(ctx context.Context, id string, mutate func(*cltypes.ResourceItem) error) (cltypes.ResourceItem, error)

	ListLostFound(ctx context.Context, query LostFoundListQuery) ([]cltypes.LostFoundItem, error)
	GetLostFound(ctx context.Context, id string) (cltypes.LostFoundItem, error)
	SaveLostFound(ctx context.Context, item cltypes.LostFoundItem) (cltypes.LostFoundItem, error)
	UpdateLostFound(ctx context.Context, id string, mutate func(*cltypes.LostFoundItem) error) (cltypes.LostFoundItem, error)

	ListCarpools(ctx context.Context, query CarpoolListQuery) ([]cltypes.CarpoolItem, error)
	GetCarpool(ctx context.Context, id string) (cltypes.CarpoolItem, error)
	SaveCarpool(ctx context.Context, item cltypes.CarpoolItem) (cltypes.CarpoolItem, error)
	UpdateCarpool(ctx context.Context, id string, mutate func(*cltypes.CarpoolItem) error) (cltypes.CarpoolItem, error)

	ListMeetups(ctx context.Context, query MeetupListQuery) ([]cltypes.MeetupItem, error)
	GetMeetup(ctx context.Context, id string) (cltypes.MeetupItem, error)
	SaveMeetup(ctx context.Context, item cltypes.MeetupItem) (cltypes.MeetupItem, error)
	UpdateMeetup(ctx context.Context, id string, mutate func(*cltypes.MeetupItem) error) (cltypes.MeetupItem, error)

	NextID(ctx context.Context, prefix string) (string, error)
}
