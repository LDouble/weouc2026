package service

import (
	"context"
	"time"

	cltypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/campus_life/types"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/auth"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/httpx"
)

func (s *Service) ListFeed(ctx context.Context, principal auth.Principal, query cltypes.FeedQuery) (map[string]any, error) {
	currentUserID, visibleStatuses, includeAllStatus := buildVisibilityFilter(principal)

	filter := cltypes.FeedFilter{
		Pagination:        query.Pagination,
		FeedTypes:         query.FeedTypes,
		Keyword:           query.Keyword,
		CurrentUserID:     currentUserID,
		VisibleStatuses:   visibleStatuses,
		IncludeAllStatus:  includeAllStatus,
		AcceptorUserID:    currentUserID,
		ParticipantUserID: currentUserID,
	}

	items, total, err := s.repository.ListForFeed(ctx, filter)
	if err != nil {
		return nil, httpx.Internal("读取动态失败", err)
	}

	list := make([]map[string]any, 0, len(items))
	for _, item := range items {
		switch item.ContentType {
		case cltypes.ContentTypeMarket:
			resolvedImages := resolveManagedURLs(ctx, s.storageProvider, item.Images)
			resolvedImage := resolveManagedURL(ctx, s.storageProvider, firstNonEmpty(item.Images...))
			if resolvedImage == "" {
				resolvedImage = firstNonEmpty(resolvedImages...)
			}
			list = append(list, map[string]any{
				"id":              item.ID.Hex(),
				"feed_type":       "market",
				"feed_type_label": "二手交易",
				"title":           item.Title,
				"desc":            item.Desc,
				"publisher":       publisherName(item, principal),
				"created_at":      item.CreatedAt.Format(time.RFC3339),
				"status":          item.Status,
				"image":           resolvedImage,
				"extra": map[string]any{
					"images":   resolvedImages,
					"likes":    item.Likes,
					"comments": 0,
				},
			})
		case cltypes.ContentTypeErrand:
			role := errandUserRole(item, principal)
			list = append(list, map[string]any{
				"id":              item.ID.Hex(),
				"feed_type":       "errand",
				"feed_type_label": "校园跑腿",
				"title":           item.Title,
				"desc":            item.Desc,
				"publisher":       publisherName(item, principal),
				"created_at":      item.CreatedAt.Format(time.RFC3339),
				"status":          item.Status,
				"extra": map[string]any{
					"likes":     0,
					"comments":  0,
					"user_role": role,
				},
			})
		case cltypes.ContentTypeResource:
			list = append(list, map[string]any{
				"id":              item.ID.Hex(),
				"feed_type":       "resource",
				"feed_type_label": "资料共享",
				"title":           item.Title,
				"desc":            item.Desc,
				"publisher":       publisherName(item, principal),
				"created_at":      item.CreatedAt.Format(time.RFC3339),
				"status":          item.Status,
			})
		case cltypes.ContentTypeLostFound:
			list = append(list, map[string]any{
				"id":              item.ID.Hex(),
				"feed_type":       "lostFound",
				"feed_type_label": "失物招领",
				"title":           item.Title,
				"desc":            item.Desc,
				"publisher":       publisherName(item, principal),
				"created_at":      item.CreatedAt.Format(time.RFC3339),
				"status":          item.Status,
			})
		case cltypes.ContentTypeCarpool:
			cp, _ := unmarshalPayload[cltypes.CarpoolPayload](item.TypePayload)
			now := time.Now().In(chinaLocation)
			list = append(list, map[string]any{
				"id":              item.ID.Hex(),
				"feed_type":       "carpool",
				"feed_type_label": "校园拼车",
				"title":           carpoolTitleFromPayload(cp, now),
				"desc":            carpoolFeedDescFromPayload(cp, now),
				"publisher":       publisherName(item, principal),
				"created_at":      item.CreatedAt.Format(time.RFC3339),
				"status":          item.Status,
				"extra": map[string]any{
					"category":   normalizedCarpoolCategoryFromPayload(cp, now),
					"from":       cp.From,
					"to":         cp.To,
					"time":       formatCarpoolTravelText(cp.TravelAt, now),
					"type":       cp.Type,
					"seats_text": cp.SeatsText,
					"price":      cp.Price,
					"tags":       append([]string(nil), item.Tags...),
					"comments":   0,
				},
			})
		case cltypes.ContentTypeMeetup:
			mp, _ := unmarshalPayload[cltypes.MeetupPayload](item.TypePayload)
			list = append(list, map[string]any{
				"id":              item.ID.Hex(),
				"feed_type":       "meetup",
				"feed_type_label": "校园组局",
				"title":           item.Title,
				"desc":            meetupFeedDescFromPayload(mp),
				"publisher":       publisherName(item, principal),
				"created_at":      item.CreatedAt.Format(time.RFC3339),
				"status":          item.Status,
				"extra": map[string]any{
					"category":        mp.Category,
					"location":        mp.Location,
					"start_at":        mp.StartAt.In(chinaLocation).Format(time.RFC3339),
					"joined_count":    meetupJoinedCountFromPayload(mp),
					"remaining_seats": meetupRemainingSeatsFromPayload(mp),
					"fee_text":        mp.FeeText,
					"tags":            append([]string(nil), item.Tags...),
					"comments":        0,
				},
			})
		}
	}

	return listEnvelope(list, int(total), query.Pagination), nil
}
