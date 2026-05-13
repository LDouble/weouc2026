package service

import (
	"context"
	"errors"
	"strings"
	"time"

	clrepo "github.com/liangluo/weouc2026/services/api-server/internal/modules/campus_life/repo"
	cltypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/campus_life/types"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/auth"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/httpx"
)

func (s *Service) ListReviewQueue(ctx context.Context, query cltypes.ReviewQuery) (map[string]any, error) {
	type row struct {
		createdAt time.Time
		payload   map[string]any
	}

	rows := make([]row, 0)
	appendRow := func(contentType string, createdAt time.Time, status string, title, desc, publisher, id string, extra map[string]any) {
		if !matchReviewQuery(query, contentType, status, title, desc, publisher, id) {
			return
		}
		rows = append(rows, row{
			createdAt: createdAt,
			payload: map[string]any{
				"content_type":  contentType,
				"content_id":    id,
				"title":         title,
				"desc":          desc,
				"publisher":     publisher,
				"created_at":    createdAt.Format(time.RFC3339),
				"review_status": status,
				"extra":         extra,
			},
		})
	}

	allTypes := []string{
		cltypes.ContentTypeMarket,
		cltypes.ContentTypeErrand,
		cltypes.ContentTypeResource,
		cltypes.ContentTypeLostFound,
		cltypes.ContentTypeCarpool,
		cltypes.ContentTypeMeetup,
	}

	for _, contentType := range allTypes {
		items, _, err := s.repository.ListByType(ctx, contentType, cltypes.ContentFilter{Pagination: query.Pagination})
		if err != nil {
			return nil, httpx.Internal("读取审核列表失败", err)
		}
		for _, item := range items {
			switch item.ContentType {
			case cltypes.ContentTypeMarket:
				mp, _ := unmarshalPayload[cltypes.MarketPayload](item.TypePayload)
				appendRow("market", item.CreatedAt, item.Status, item.Title, item.Desc, "校园用户", item.ID.Hex(), map[string]any{
					"category":       mp.Category,
					"price":          mp.Price,
					"original_price": mp.OriginalPrice,
					"condition":      mp.Condition,
					"trade_mode":     mp.TradeMode,
				})
			case cltypes.ContentTypeErrand:
				ep, _ := unmarshalPayload[cltypes.ErrandPayload](item.TypePayload)
				appendRow("errand", item.CreatedAt, item.Status, item.Title, item.Desc, "校园用户", item.ID.Hex(), map[string]any{
					"category":    ep.Category,
					"status":      item.Status,
					"route_start": ep.RouteStart,
					"route_end":   ep.RouteEnd,
					"deadline":    ep.Deadline.Format(time.RFC3339),
					"reward":      ep.Reward,
				})
			case cltypes.ContentTypeResource:
				rp, _ := unmarshalPayload[cltypes.ResourcePayload](item.TypePayload)
				appendRow("resource", item.CreatedAt, item.Status, item.Title, item.Desc, "校园用户", item.ID.Hex(), map[string]any{
					"category":    rp.Category,
					"course_name": rp.CourseName,
					"file_type":   rp.FileType,
					"file_size":   rp.FileSize,
				})
			case cltypes.ContentTypeLostFound:
				lp, _ := unmarshalPayload[cltypes.LostFoundPayload](item.TypePayload)
				appendRow("lostFound", item.CreatedAt, item.Status, item.Title, item.Desc, "校园用户", item.ID.Hex(), map[string]any{
					"category":   lp.Category,
					"type":       lp.Type,
					"location":   lp.Location,
					"event_time": lp.EventTime,
				})
			case cltypes.ContentTypeCarpool:
				cp, _ := unmarshalPayload[cltypes.CarpoolPayload](item.TypePayload)
				now := time.Now().In(chinaLocation)
				appendRow("carpool", item.CreatedAt, item.Status, carpoolTitleFromPayload(cp, now), cp.Note, "校园用户", item.ID.Hex(), map[string]any{
					"category":   normalizedCarpoolCategoryFromPayload(cp, now),
					"from":       cp.From,
					"to":         cp.To,
					"time":       formatCarpoolTravelText(cp.TravelAt, now),
					"seats_text": cp.SeatsText,
				})
			case cltypes.ContentTypeMeetup:
				mp, _ := unmarshalPayload[cltypes.MeetupPayload](item.TypePayload)
				appendRow("meetup", item.CreatedAt, item.Status, item.Title, item.Desc, "校园用户", item.ID.Hex(), map[string]any{
					"category":         mp.Category,
					"location":         mp.Location,
					"start_at":         mp.StartAt.In(chinaLocation).Format(time.RFC3339),
					"max_participants": mp.MaxParticipants,
					"joined_count":     meetupJoinedCountFromPayload(mp),
					"status":           item.Status,
				})
			}
		}
	}

	items := paginateRows(rows, query.Pagination)
	list := make([]map[string]any, 0, len(items))
	for _, item := range items {
		list = append(list, item.payload)
	}

	return listEnvelope(list, len(rows), query.Pagination), nil
}

func (s *Service) UpdateReviewStatus(ctx context.Context, principal auth.Principal, request cltypes.ReviewUpdateRequest) error {
	contentType := strings.TrimSpace(request.ContentType)
	contentID := strings.TrimSpace(request.ContentID)
	reviewStatus := strings.TrimSpace(strings.ToLower(request.ReviewStatus))

	if contentType == "" || contentID == "" {
		return httpx.BadRequest("content_type 和 content_id 不能为空", nil)
	}
	if !isSupportedReviewStatus(reviewStatus) {
		return httpx.BadRequest("review_status 仅支持 reviewing/published/rejected/offline", nil)
	}

	_, err := s.repository.Update(ctx, contentID, func(item *cltypes.CommunityContent) error {
		switch contentType {
		case "market", "errand", "resource", "lostFound", "carpool":
			if reviewStatus == cltypes.StatusPublished {
				item.Status = cltypes.StatusPublished
			} else {
				item.Status = reviewStatus
			}
		case "meetup":
			if reviewStatus == cltypes.StatusPublished {
				mp, _ := unmarshalPayload[cltypes.MeetupPayload](item.TypePayload)
				mp = refreshMeetupPayloadStatus(mp)
				item.TypePayload = marshalPayload(mp)
				if mp.MaxParticipants > 0 && len(mp.ParticipantUserIDs)+1 >= mp.MaxParticipants {
					item.Status = cltypes.StatusFull
				} else {
					item.Status = cltypes.StatusOpen
				}
			} else {
				item.Status = reviewStatus
			}
		default:
			return httpx.BadRequest("content_type 仅支持 market/errand/resource/lostFound/carpool/meetup", nil)
		}
		item.UpdatedBy = principal.UserID
		return nil
	})

	if errors.Is(err, clrepo.ErrNotFound) {
		return httpx.NotFound("待审核内容不存在", nil)
	}
	if err != nil {
		if isAppError(err) {
			return err
		}
		return httpx.Internal("更新审核状态失败", err)
	}

	s.recordAudit(ctx, principal, "campus_life.review.update", contentType, contentID, "校园生活审核状态更新成功", map[string]any{
		"review_status": reviewStatus,
	})

	return nil
}
