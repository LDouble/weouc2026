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

func (s *Service) ListMarket(ctx context.Context, principal auth.Principal, query cltypes.MarketQuery) (map[string]any, error) {
	currentUserID, visibleStatuses, includeAllStatus := buildVisibilityFilter(principal)

	items, total, err := s.repository.ListByType(ctx, cltypes.ContentTypeMarket, cltypes.ContentFilter{
		Pagination:      query.Pagination,
		Category:        query.Category,
		Keyword:         query.Keyword,
		CurrentUserID:   currentUserID,
		VisibleStatuses: visibleStatuses,
		IncludeAllStatus: includeAllStatus,
	})
	if err != nil {
		return nil, httpx.Internal("读取二手列表失败", err)
	}

	list := make([]map[string]any, 0, len(items))
	for _, item := range items {
		canView := canViewContact(principal, item.PublisherUserID)
		isOwner := item.PublisherUserID == principal.UserID
		mp, _ := unmarshalPayload[cltypes.MarketPayload](item.TypePayload)
		resolvedImages := resolveManagedURLs(ctx, s.storageProvider, item.Images)
		resolvedImage := resolveManagedURL(ctx, s.storageProvider, firstNonEmpty(item.Images...))
		if resolvedImage == "" {
			resolvedImage = firstNonEmpty(resolvedImages...)
		}
		list = append(list, map[string]any{
			"id":                item.ID.Hex(),
			"title":             item.Title,
			"desc":              item.Desc,
			"publisher":         publisherName(item, principal),
			"publisher_initial": initialOf(publisherName(item, principal)),
			"image":             resolvedImage,
			"likes":             item.Likes,
			"liked":             item.LikedByUserIDs[principal.UserID],
			"created_at":        item.CreatedAt.Format(time.RFC3339),
			"status":            item.Status,
			"is_owner":          isOwner,
			"can_edit":          canEditContent(isOwner, item.Status),
			"can_delete":        canDeleteContent(isOwner, item.Status),
			"can_favorite":      !isOwner && principal.Authenticated,
			"extra": map[string]any{
				"category":       mp.Category,
				"price":          mp.Price,
				"original_price": mp.OriginalPrice,
				"condition":      mp.Condition,
				"trade_mode":     mp.TradeMode,
				"contact":        visibleValue(canView, item.Contact),
				"images":         resolvedImages,
				"likes":          item.Likes,
				"is_favorited":   item.LikedByUserIDs[principal.UserID],
			},
		})
	}

	return listEnvelope(list, int(total), query.Pagination), nil
}

func (s *Service) GetMarketDetail(ctx context.Context, principal auth.Principal, id string) (map[string]any, error) {
	item, err := s.repository.GetByID(ctx, id)
	if errors.Is(err, clrepo.ErrNotFound) {
		return nil, httpx.NotFound("二手商品不存在", nil)
	}
	if err != nil {
		return nil, httpx.Internal("读取二手详情失败", err)
	}
	if item.ContentType != cltypes.ContentTypeMarket {
		return nil, httpx.NotFound("二手商品不存在", nil)
	}
	if err := ensureContentVisible(principal, item, "二手商品不存在"); err != nil {
		return nil, err
	}

	canView := canViewContact(principal, item.PublisherUserID)
	role := simpleUserRole(item.PublisherUserID, principal)
	isOwner := role == "publisher"
	mp, _ := unmarshalPayload[cltypes.MarketPayload](item.TypePayload)
	resolvedImages := resolveManagedURLs(ctx, s.storageProvider, item.Images)
	resolvedImage := resolveManagedURL(ctx, s.storageProvider, firstNonEmpty(item.Images...))
	if resolvedImage == "" {
		resolvedImage = firstNonEmpty(resolvedImages...)
	}
	return map[string]any{
		"id":                item.ID.Hex(),
		"title":             item.Title,
		"desc":              item.Desc,
		"publisher":         publisherName(item, principal),
		"publisher_initial": initialOf(publisherName(item, principal)),
		"image":             resolvedImage,
		"likes":             item.Likes,
		"liked":             item.LikedByUserIDs[principal.UserID],
		"created_at":        item.CreatedAt.Format(time.RFC3339),
		"status":            item.Status,
		"user_role":         role,
		"is_owner":          isOwner,
		"can_view_contact":  canView,
		"can_edit":          canEditContent(isOwner, item.Status),
		"can_delete":        canDeleteContent(isOwner, item.Status),
		"can_favorite":      !isOwner && principal.Authenticated,
		"extra": map[string]any{
			"category":       mp.Category,
			"price":          mp.Price,
			"original_price": mp.OriginalPrice,
			"condition":      mp.Condition,
			"trade_mode":     mp.TradeMode,
			"contact":        visibleValue(canView, item.Contact),
			"images":         resolvedImages,
			"likes":          item.Likes,
			"is_favorited":   item.LikedByUserIDs[principal.UserID],
		},
	}, nil
}

func (s *Service) PublishMarket(ctx context.Context, principal auth.Principal, request cltypes.MarketPublishRequest) (map[string]any, error) {
	if strings.TrimSpace(request.Title) == "" || strings.TrimSpace(request.Desc) == "" {
		return nil, httpx.BadRequest("标题和描述不能为空", nil)
	}
	item := cltypes.CommunityContent{
		ContentType:     cltypes.ContentTypeMarket,
		Title:           strings.TrimSpace(request.Title),
		Desc:            strings.TrimSpace(request.Desc),
		Status:          cltypes.StatusReviewing,
		PublisherUserID: principal.UserID,
		Contact:         strings.TrimSpace(request.Contact),
		Images:          append([]string(nil), request.Images...),
		TypePayload: marshalPayload(cltypes.MarketPayload{
			Category:      strings.TrimSpace(request.Category),
			Price:         strings.TrimSpace(request.Price),
			OriginalPrice: strings.TrimSpace(request.OriginalPrice),
			Condition:     strings.TrimSpace(request.Condition),
			TradeMode:     strings.TrimSpace(request.TradeMode),
		}),
		LikedByUserIDs: map[string]bool{},
		CreatedBy:      principal.UserID,
		UpdatedBy:      principal.UserID,
	}
	item, err := s.repository.Save(ctx, item)
	if err != nil {
		return nil, httpx.Internal("保存二手信息失败", err)
	}

	s.recordAudit(ctx, principal, "campus_life.market.publish", "market", item.ID.Hex(), "二手信息发布成功", map[string]any{
		"status":   item.Status,
		"category": request.Category,
	})

	return map[string]any{"id": item.ID.Hex()}, nil
}

func (s *Service) FavoriteMarket(ctx context.Context, principal auth.Principal, request cltypes.FavoriteMarketRequest) error {
	if request.ProductID == "" || principal.UserID == "" {
		return httpx.BadRequest("缺少商品或用户信息", nil)
	}

	_, err := s.repository.Update(ctx, request.ProductID, func(item *cltypes.CommunityContent) error {
		if item.LikedByUserIDs == nil {
			item.LikedByUserIDs = map[string]bool{}
		}
		switch request.Action {
		case "add":
			if !item.LikedByUserIDs[principal.UserID] {
				item.LikedByUserIDs[principal.UserID] = true
				item.Likes++
			}
		case "remove":
			if item.LikedByUserIDs[principal.UserID] {
				delete(item.LikedByUserIDs, principal.UserID)
				if item.Likes > 0 {
					item.Likes--
				}
			}
		default:
			return httpx.BadRequest("action 仅支持 add/remove", nil)
		}
		return nil
	})
	if err != nil {
		if err == clrepo.ErrNotFound {
			return httpx.NotFound("二手商品不存在", nil)
		}
		if isAppError(err) {
			return err
		}
		return httpx.Internal("更新二手收藏状态失败", err)
	}

	return nil
}

func (s *Service) DeleteMarket(ctx context.Context, principal auth.Principal, id string) error {
	_, err := s.repository.Update(ctx, id, func(item *cltypes.CommunityContent) error {
		if item.PublisherUserID != principal.UserID {
			return httpx.Forbidden("只有发布者可以下架", nil)
		}
		if item.Status != cltypes.StatusPublished && item.Status != cltypes.StatusReviewing && item.Status != cltypes.StatusRejected {
			return httpx.BadRequest("当前状态不允许下架", nil)
		}
		item.Status = cltypes.StatusOffline
		item.UpdatedBy = principal.UserID
		return nil
	})
	if err != nil {
		if errors.Is(err, clrepo.ErrNotFound) {
			return httpx.NotFound("二手商品不存在", nil)
		}
		if isAppError(err) {
			return err
		}
		return httpx.Internal("下架二手商品失败", err)
	}
	s.recordAudit(ctx, principal, "campus_life.market.delete", "market", id, "二手商品下架成功", nil)
	return nil
}
