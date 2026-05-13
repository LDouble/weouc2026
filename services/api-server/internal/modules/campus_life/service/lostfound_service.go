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

func (s *Service) ListLostFound(ctx context.Context, principal auth.Principal, query cltypes.LostFoundQuery) (map[string]any, error) {
	currentUserID, visibleStatuses, includeAllStatus := buildVisibilityFilter(principal)

	items, total, err := s.repository.ListByType(ctx, cltypes.ContentTypeLostFound, cltypes.ContentFilter{
		Pagination:       query.Pagination,
		Category:         query.Category,
		Keyword:          query.Keyword,
		SubType:          query.Type,
		CurrentUserID:    currentUserID,
		VisibleStatuses:  visibleStatuses,
		IncludeAllStatus: includeAllStatus,
	})
	if err != nil {
		return nil, httpx.Internal("读取失物招领列表失败", err)
	}

	list := make([]map[string]any, 0, len(items))
	for _, item := range items {
		canView := canViewContact(principal, item.PublisherUserID)
		isOwner := item.PublisherUserID == principal.UserID
		lp, _ := unmarshalPayload[cltypes.LostFoundPayload](item.TypePayload)
		list = append(list, map[string]any{
			"id":                item.ID.Hex(),
			"title":             item.Title,
			"desc":              item.Desc,
			"publisher":         publisherName(item, principal),
			"publisher_initial": initialOf(publisherName(item, principal)),
			"created_at":        item.CreatedAt.Format(time.RFC3339),
			"status":            item.Status,
			"is_owner":          isOwner,
			"can_edit":          canEditContent(isOwner, item.Status),
			"can_delete":        canDeleteContent(isOwner, item.Status),
			"can_mark_resolved": isOwner && item.Status == cltypes.StatusPublished,
			"extra": map[string]any{
				"type":         lp.Type,
				"category":     lp.Category,
				"location":     lp.Location,
				"event_time":   lp.EventTime,
				"item_feature": lp.ItemFeature,
				"contact":      visibleValue(canView, item.Contact),
			},
		})
	}
	return listEnvelope(list, int(total), query.Pagination), nil
}

func (s *Service) GetLostFoundDetail(ctx context.Context, principal auth.Principal, id string) (map[string]any, error) {
	item, err := s.repository.GetByID(ctx, id)
	if errors.Is(err, clrepo.ErrNotFound) {
		return nil, httpx.NotFound("失物招领不存在", nil)
	}
	if err != nil {
		return nil, httpx.Internal("读取失物招领详情失败", err)
	}
	if item.ContentType != cltypes.ContentTypeLostFound {
		return nil, httpx.NotFound("失物招领不存在", nil)
	}
	if err := ensureContentVisible(principal, item, "失物招领不存在"); err != nil {
		return nil, err
	}
	canView := canViewContact(principal, item.PublisherUserID)
	role := simpleUserRole(item.PublisherUserID, principal)
	isOwner := role == "publisher"
	lp, _ := unmarshalPayload[cltypes.LostFoundPayload](item.TypePayload)
	return map[string]any{
		"id":                item.ID.Hex(),
		"title":             item.Title,
		"desc":              item.Desc,
		"publisher":         publisherName(item, principal),
		"publisher_initial": initialOf(publisherName(item, principal)),
		"created_at":        item.CreatedAt.Format(time.RFC3339),
		"status":            item.Status,
		"user_role":         role,
		"is_owner":          isOwner,
		"can_view_contact":  canView,
		"can_edit":          canEditContent(isOwner, item.Status),
		"can_delete":        canDeleteContent(isOwner, item.Status),
		"can_mark_resolved": isOwner && item.Status == cltypes.StatusPublished,
		"extra": map[string]any{
			"type":         lp.Type,
			"category":     lp.Category,
			"location":     lp.Location,
			"event_time":   lp.EventTime,
			"item_feature": lp.ItemFeature,
			"contact":      visibleValue(canView, item.Contact),
		},
	}, nil
}

func (s *Service) PublishLostFound(ctx context.Context, principal auth.Principal, request cltypes.LostFoundPublishRequest) (map[string]any, error) {
	if strings.TrimSpace(request.Title) == "" || strings.TrimSpace(request.Contact) == "" {
		return nil, httpx.BadRequest("标题和联系方式不能为空", nil)
	}
	item := cltypes.CommunityContent{
		ContentType:     cltypes.ContentTypeLostFound,
		Title:           strings.TrimSpace(request.Title),
		Desc:            strings.TrimSpace(request.Desc),
		Status:          cltypes.StatusReviewing,
		PublisherUserID: principal.UserID,
		Contact:         strings.TrimSpace(request.Contact),
		TypePayload: marshalPayload(cltypes.LostFoundPayload{
			Type:        strings.TrimSpace(request.Type),
			Category:    strings.TrimSpace(request.Category),
			Location:    strings.TrimSpace(request.Location),
			EventTime:   strings.TrimSpace(request.EventTime),
			ItemFeature: strings.TrimSpace(request.ItemFeature),
		}),
		CreatedBy: principal.UserID,
		UpdatedBy: principal.UserID,
	}
	item, err := s.repository.Save(ctx, item)
	if err != nil {
		return nil, httpx.Internal("保存失物招领失败", err)
	}
	s.recordAudit(ctx, principal, "campus_life.lost_found.publish", "lostFound", item.ID.Hex(), "失物招领发布成功", map[string]any{
		"status":   item.Status,
		"type":     request.Type,
		"category": request.Category,
	})
	return map[string]any{"id": item.ID.Hex()}, nil
}

func (s *Service) DeleteLostFound(ctx context.Context, principal auth.Principal, id string) error {
	_, err := s.repository.Update(ctx, id, func(item *cltypes.CommunityContent) error {
		if item.ContentType != cltypes.ContentTypeLostFound {
			return httpx.NotFound("失物招领不存在", nil)
		}
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
			return httpx.NotFound("失物招领不存在", nil)
		}
		if isAppError(err) {
			return err
		}
		return httpx.Internal("下架失物招领失败", err)
	}
	s.recordAudit(ctx, principal, "campus_life.lost_found.delete", "lostFound", id, "失物招领下架成功", nil)
	return nil
}

func (s *Service) MarkLostFoundResolved(ctx context.Context, principal auth.Principal, id string) error {
	_, err := s.repository.Update(ctx, id, func(item *cltypes.CommunityContent) error {
		if item.ContentType != cltypes.ContentTypeLostFound {
			return httpx.NotFound("失物招领不存在", nil)
		}
		if item.PublisherUserID != principal.UserID {
			return httpx.Forbidden("只有发布者可以标记已找到", nil)
		}
		if item.Status != cltypes.StatusPublished {
			return httpx.BadRequest("当前状态不允许标记已找到", nil)
		}
		item.Status = cltypes.StatusResolved
		item.UpdatedBy = principal.UserID
		return nil
	})
	if err != nil {
		if err == clrepo.ErrNotFound {
			return httpx.NotFound("失物招领不存在", nil)
		}
		if isAppError(err) {
			return err
		}
		return httpx.Internal("标记失物招领已找到失败", err)
	}
	s.recordAudit(ctx, principal, "campus_life.lost_found.resolve", "lostFound", id, "失物招领标记已找到", nil)
	return nil
}
