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

func (s *Service) ListErrands(ctx context.Context, principal auth.Principal, query cltypes.ErrandQuery) (map[string]any, error) {
	currentUserID, visibleStatuses, includeAllStatus := buildVisibilityFilter(principal)

	items, total, err := s.repository.ListByType(ctx, cltypes.ContentTypeErrand, cltypes.ContentFilter{
		Pagination:       query.Pagination,
		Category:         query.Category,
		Keyword:          query.Keyword,
		CurrentUserID:    currentUserID,
		VisibleStatuses:  visibleStatuses,
		IncludeAllStatus: includeAllStatus,
		AcceptorUserID:   currentUserID,
	})
	if err != nil {
		return nil, httpx.Internal("读取跑腿列表失败", err)
	}

	list := make([]map[string]any, 0, len(items))
	for _, item := range items {
		ep, _ := unmarshalPayload[cltypes.ErrandPayload](item.TypePayload)
		role := errandUserRole(item, principal)
		isOwner := role == "publisher"
		list = append(list, map[string]any{
			"id":                item.ID.Hex(),
			"category":          ep.Category,
			"title":             item.Title,
			"desc":              item.Desc,
			"route_start":       ep.RouteStart,
			"route_end":         ep.RouteEnd,
			"deadline":          ep.Deadline.Format(time.RFC3339),
			"reward":            ep.Reward,
			"status":            item.Status,
			"user_role":         role,
			"is_owner":          isOwner,
			"is_accepted":       item.Status == cltypes.StatusAccepted,
			"views":             0,
			"publisher":         publisherName(item, principal),
			"publisher_initial": initialOf(publisherName(item, principal)),
			"created_at":        item.CreatedAt.Format(time.RFC3339),
			"can_edit":          canEditContent(isOwner, item.Status),
			"can_delete":        canDeleteContent(isOwner, item.Status),
			"can_accept":        !isOwner && item.Status == cltypes.StatusPublished && principal.Authenticated,
			"can_cancel_accept": role == "acceptor" && item.Status == cltypes.StatusAccepted,
		})
	}

	return listEnvelope(list, int(total), query.Pagination), nil
}

func (s *Service) GetErrandDetail(ctx context.Context, principal auth.Principal, id string) (map[string]any, error) {
	item, err := s.repository.GetByID(ctx, id)
	if errors.Is(err, clrepo.ErrNotFound) {
		return nil, httpx.NotFound("跑腿任务不存在", nil)
	}
	if err != nil {
		return nil, httpx.Internal("读取跑腿详情失败", err)
	}
	if item.ContentType != cltypes.ContentTypeErrand {
		return nil, httpx.NotFound("跑腿任务不存在", nil)
	}
	if err := ensureContentVisible(principal, item, "跑腿任务不存在"); err != nil {
		return nil, err
	}
	role := errandUserRole(item, principal)
	isOwner := role == "publisher"
	canView := canViewContact(principal, item.PublisherUserID)
	ep, _ := unmarshalPayload[cltypes.ErrandPayload](item.TypePayload)
	resolvedImages := resolveManagedURLs(ctx, s.storageProvider, item.Images)
	return map[string]any{
		"item": map[string]any{
			"id":                item.ID.Hex(),
			"category":          ep.Category,
			"title":             item.Title,
			"desc":              item.Desc,
			"route_start":       ep.RouteStart,
			"route_end":         ep.RouteEnd,
			"deadline":          ep.Deadline.Format(time.RFC3339),
			"reward":            ep.Reward,
			"contact":           visibleValue(canView, item.Contact),
			"status":            item.Status,
			"is_accepted":       item.Status == cltypes.StatusAccepted,
			"images":            resolvedImages,
			"publisher":         publisherName(item, principal),
			"publisher_initial": initialOf(publisherName(item, principal)),
			"created_at":        item.CreatedAt.Format(time.RFC3339),
			"extra": map[string]any{
				"category":    ep.Category,
				"route_start": ep.RouteStart,
				"route_end":   ep.RouteEnd,
				"deadline":    ep.Deadline.Format(time.RFC3339),
				"reward":      ep.Reward,
				"contact":     visibleValue(canView, item.Contact),
				"images":      resolvedImages,
				"status":      item.Status,
				"urgent":      ep.Urgent,
			},
		},
		"user_role":         role,
		"is_owner":          isOwner,
		"can_view_contact":  canView,
		"can_edit":          canEditContent(isOwner, item.Status),
		"can_delete":        canDeleteContent(isOwner, item.Status),
		"can_accept":        !isOwner && item.Status == cltypes.StatusPublished && principal.Authenticated,
		"can_cancel_accept": role == "acceptor" && item.Status == cltypes.StatusAccepted,
	}, nil
}

func (s *Service) PublishErrand(ctx context.Context, principal auth.Principal, request cltypes.ErrandPublishRequest) (map[string]any, error) {
	if strings.TrimSpace(request.Title) == "" || strings.TrimSpace(request.Contact) == "" {
		return nil, httpx.BadRequest("跑腿标题和联系方式不能为空", nil)
	}
	deadline, err := time.Parse(time.RFC3339, request.Deadline)
	if err != nil {
		return nil, httpx.BadRequest("deadline 必须为 RFC3339 时间", nil)
	}
	item := cltypes.CommunityContent{
		ContentType:     cltypes.ContentTypeErrand,
		Title:           strings.TrimSpace(request.Title),
		Desc:            strings.TrimSpace(request.Desc),
		Status:          cltypes.StatusReviewing,
		PublisherUserID: principal.UserID,
		Contact:         strings.TrimSpace(request.Contact),
		Images:          append([]string(nil), request.Images...),
		TypePayload: marshalPayload(cltypes.ErrandPayload{
			Category:   strings.TrimSpace(request.Category),
			RouteStart: strings.TrimSpace(request.RouteStart),
			RouteEnd:   strings.TrimSpace(request.RouteEnd),
			Deadline:   deadline,
			Reward:     strings.TrimSpace(request.Reward),
			Urgent:     request.Urgent,
		}),
		CreatedBy: principal.UserID,
		UpdatedBy: principal.UserID,
	}
	item, err = s.repository.Save(ctx, item)
	if err != nil {
		return nil, httpx.Internal("保存跑腿任务失败", err)
	}
	s.recordAudit(ctx, principal, "campus_life.errand.publish", "errand", item.ID.Hex(), "跑腿任务发布成功", map[string]any{
		"status":   item.Status,
		"category": request.Category,
	})
	return map[string]any{"id": item.ID.Hex()}, nil
}

func (s *Service) AcceptErrand(ctx context.Context, principal auth.Principal, taskID string) error {
	_, err := s.repository.Update(ctx, taskID, func(item *cltypes.CommunityContent) error {
		if item.ContentType != cltypes.ContentTypeErrand {
			return httpx.NotFound("跑腿任务不存在", nil)
		}
		if item.PublisherUserID == principal.UserID {
			return httpx.BadRequest("不能接自己发布的任务", nil)
		}
		switch item.Status {
		case cltypes.StatusPublished:
		case cltypes.StatusReviewing:
			return httpx.BadRequest("该任务仍在审核中，暂不可接单", nil)
		case cltypes.StatusRejected:
			return httpx.BadRequest("该任务审核未通过，无法接单", nil)
		case cltypes.StatusOffline:
			return httpx.BadRequest("该任务已下线，无法接单", nil)
		case cltypes.StatusCancelled:
			return httpx.BadRequest("该任务已取消", nil)
		case cltypes.StatusAccepted:
			return httpx.BadRequest("该任务已被接单", nil)
		default:
			return httpx.BadRequest("该任务当前状态无法接单", nil)
		}
		item.Status = cltypes.StatusAccepted
		item.UpdatedBy = principal.UserID
		ep, _ := unmarshalPayload[cltypes.ErrandPayload](item.TypePayload)
		ep.AcceptorUserID = principal.UserID
		item.TypePayload = marshalPayload(ep)
		return nil
	})
	if err != nil {
		if err == clrepo.ErrNotFound {
			return httpx.NotFound("跑腿任务不存在", nil)
		}
		if isAppError(err) {
			return err
		}
		return httpx.Internal("更新跑腿接单状态失败", err)
	}

	return nil
}

func (s *Service) CancelErrandPublish(ctx context.Context, principal auth.Principal, taskID string) error {
	_, err := s.repository.Update(ctx, taskID, func(item *cltypes.CommunityContent) error {
		if item.ContentType != cltypes.ContentTypeErrand {
			return httpx.NotFound("跑腿任务不存在", nil)
		}
		if item.PublisherUserID != principal.UserID {
			return httpx.Forbidden("只有发布者可以取消发布", nil)
		}
		item.Status = cltypes.StatusCancelled
		item.UpdatedBy = principal.UserID
		return nil
	})
	if err != nil {
		if err == clrepo.ErrNotFound {
			return httpx.NotFound("跑腿任务不存在", nil)
		}
		if isAppError(err) {
			return err
		}
		return httpx.Internal("取消跑腿发布失败", err)
	}

	return nil
}

func (s *Service) CancelErrandAccept(ctx context.Context, principal auth.Principal, taskID string) error {
	_, err := s.repository.Update(ctx, taskID, func(item *cltypes.CommunityContent) error {
		if item.ContentType != cltypes.ContentTypeErrand {
			return httpx.NotFound("跑腿任务不存在", nil)
		}
		ep, _ := unmarshalPayload[cltypes.ErrandPayload](item.TypePayload)
		if ep.AcceptorUserID != principal.UserID {
			return httpx.Forbidden("只有接单者可以取消接单", nil)
		}
		item.Status = cltypes.StatusPublished
		item.UpdatedBy = principal.UserID
		ep.AcceptorUserID = ""
		item.TypePayload = marshalPayload(ep)
		return nil
	})
	if err != nil {
		if err == clrepo.ErrNotFound {
			return httpx.NotFound("跑腿任务不存在", nil)
		}
		if isAppError(err) {
			return err
		}
		return httpx.Internal("取消跑腿接单失败", err)
	}

	return nil
}
