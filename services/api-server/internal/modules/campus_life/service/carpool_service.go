package service

import (
	"context"
	"errors"
	"strings"
	"time"

	clrepo "github.com/liangluo/weouc2026/services/api-server/internal/modules/campus_life/repo"
	cltypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/campus_life/types"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/auth"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/bmfs"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/httpx"
)

func (s *Service) ListCarpools(ctx context.Context, principal auth.Principal, query cltypes.CarpoolQuery) (map[string]any, error) {
	currentUserID, visibleStatuses, includeAllStatus := buildVisibilityFilter(principal)

	items, _, err := s.repository.ListByType(ctx, cltypes.ContentTypeCarpool, cltypes.ContentFilter{
		Pagination:       query.Pagination,
		Keyword:          query.Keyword,
		CurrentUserID:    currentUserID,
		VisibleStatuses:  visibleStatuses,
		IncludeAllStatus: includeAllStatus,
	})
	if err != nil {
		return nil, httpx.Internal("读取拼车列表失败", err)
	}

	now := time.Now().In(chinaLocation)
	machine := cltypes.GetMachine(cltypes.ContentTypeCarpool)
	filtered := make([]map[string]any, 0)
	for _, item := range items {
		cp, _ := unmarshalPayload[cltypes.CarpoolPayload](item.TypePayload)
		if query.Category != "" && query.Category != "all" && normalizedCarpoolCategoryFromPayload(cp, now) != query.Category {
			continue
		}
		canView := canViewContact(principal, item.PublisherUserID)
		isOwner := item.PublisherUserID == principal.UserID
		actx := bmfs.ActionContext{
			Principal: principal,
			IsOwner:   isOwner,
			UserRole:  simpleUserRole(item.PublisherUserID, principal),
			Now:       now,
		}
		actions := machine.AvailableActions(item.Status, actx)
		payload := buildCarpoolPayload(item, cp, canView, now)
		payload["is_owner"] = isOwner
		payload["can_edit"] = canEditContent(isOwner, item.Status) && cp.TravelAt.After(now.UTC())
		payload["can_delete"] = actions["can_delete"]
		payload["can_join_carpool"] = !isOwner && item.Status == cltypes.StatusPublished && principal.Authenticated
		filtered = append(filtered, payload)
	}

	return listEnvelope(paginateMaps(filtered, query.Pagination), len(filtered), query.Pagination), nil
}

func (s *Service) GetCarpoolDetail(ctx context.Context, principal auth.Principal, id string) (map[string]any, error) {
	item, err := s.repository.GetByID(ctx, id)
	if errors.Is(err, clrepo.ErrNotFound) {
		return nil, httpx.NotFound("拼车行程不存在", nil)
	}
	if err != nil {
		return nil, httpx.Internal("读取拼车详情失败", err)
	}
	if item.ContentType != cltypes.ContentTypeCarpool {
		return nil, httpx.NotFound("拼车行程不存在", nil)
	}
	if err := ensureContentVisible(principal, item, "拼车行程不存在"); err != nil {
		return nil, err
	}

	canView := canViewContact(principal, item.PublisherUserID)
	role := simpleUserRole(item.PublisherUserID, principal)
	isOwner := role == "publisher"
	now := time.Now().In(chinaLocation)
	actx := bmfs.ActionContext{
		Principal: principal,
		IsOwner:   isOwner,
		UserRole:  role,
		Now:       now,
	}
	actions := cltypes.GetMachine(cltypes.ContentTypeCarpool).AvailableActions(item.Status, actx)
	cp, _ := unmarshalPayload[cltypes.CarpoolPayload](item.TypePayload)
	payload := buildCarpoolPayload(item, cp, canView, now)
	payload["can_view_contact"] = canView
	payload["user_role"] = role
	payload["is_owner"] = isOwner
	payload["can_edit"] = canEditContent(isOwner, item.Status) && cp.TravelAt.After(now.UTC())
	payload["can_delete"] = actions["can_delete"]
	payload["can_join_carpool"] = !isOwner && item.Status == cltypes.StatusPublished && principal.Authenticated
	return payload, nil
}

func (s *Service) PublishCarpool(ctx context.Context, principal auth.Principal, request cltypes.CarpoolPublishRequest) (map[string]any, error) {
	if strings.TrimSpace(request.From) == "" || strings.TrimSpace(request.To) == "" {
		return nil, httpx.BadRequest("出发地和目的地不能为空", nil)
	}
	if strings.TrimSpace(request.TravelDate) == "" || strings.TrimSpace(request.TravelTime) == "" {
		return nil, httpx.BadRequest("travel_date 和 travel_time 不能为空", nil)
	}
	if strings.TrimSpace(request.Contact) == "" {
		return nil, httpx.BadRequest("联系方式不能为空", nil)
	}

	travelAt, err := parseCarpoolTravelAt(request.TravelDate, request.TravelTime)
	if err != nil {
		return nil, httpx.BadRequest("travel_date/travel_time 格式错误", nil)
	}
	now := time.Now().In(chinaLocation)
	tempCP := cltypes.CarpoolPayload{TravelAt: travelAt}
	category := firstNonEmpty(strings.TrimSpace(request.Category), normalizedCarpoolCategoryFromPayload(tempCP, now))
	if !isSupportedCarpoolCategory(category) {
		category = normalizedCarpoolCategoryFromPayload(tempCP, now)
	}

	item := cltypes.CommunityContent{
		ContentType:     cltypes.ContentTypeCarpool,
		Title:           strings.TrimSpace(request.From) + " -> " + strings.TrimSpace(request.To),
		Desc:            strings.TrimSpace(request.Note),
		Status:          cltypes.StatusReviewing,
		PublisherUserID: principal.UserID,
		Contact:         strings.TrimSpace(request.Contact),
		Tags:            sanitizeTags(request.Tags),
		TypePayload: marshalPayload(cltypes.CarpoolPayload{
			Category:  category,
			From:      strings.TrimSpace(request.From),
			To:        strings.TrimSpace(request.To),
			TravelAt:  travelAt,
			Type:      firstNonEmpty(strings.TrimSpace(request.Type), defaultCarpoolType(category)),
			SeatsText: strings.TrimSpace(request.SeatsText),
			Price:     strings.TrimSpace(request.Price),
			Note:      strings.TrimSpace(request.Note),
		}),
		CreatedBy: principal.UserID,
		UpdatedBy: principal.UserID,
	}
	item, err = s.repository.Save(ctx, item)
	if err != nil {
		return nil, httpx.Internal("保存拼车信息失败", err)
	}

	s.recordAudit(ctx, principal, "campus_life.carpool.publish", "carpool", item.ID.Hex(), "拼车信息发布成功", map[string]any{
		"status":   item.Status,
		"category": category,
	})

	return map[string]any{"id": item.ID.Hex()}, nil
}

func (s *Service) DeleteCarpool(ctx context.Context, principal auth.Principal, id string) error {
	var result *bmfs.ExecuteResult
	_, err := s.repository.Update(ctx, id, func(item *cltypes.CommunityContent) error {
		if item.ContentType != cltypes.ContentTypeCarpool {
			return httpx.NotFound("拼车行程不存在", nil)
		}
		isOwner := item.PublisherUserID == principal.UserID
		actx := bmfs.ActionContext{
			Principal: principal,
			IsOwner:   isOwner,
			UserRole:  simpleUserRole(item.PublisherUserID, principal),
			Now:       time.Now(),
		}
		machine := cltypes.GetMachine(cltypes.ContentTypeCarpool)
		var execErr error
		result, execErr = machine.Execute(item.Status, cltypes.ActionDelete, actx)
		if execErr != nil {
			if !isOwner {
				return httpx.Forbidden("只有发布者可以取消发布", nil)
			}
			return httpx.BadRequest("当前状态不允许取消发布", nil)
		}
		item.Status = result.ToStatus
		item.UpdatedBy = principal.UserID
		return nil
	})
	if err != nil {
		if errors.Is(err, clrepo.ErrNotFound) {
			return httpx.NotFound("拼车行程不存在", nil)
		}
		if isAppError(err) {
			return err
		}
		return httpx.Internal("取消拼车发布失败", err)
	}
	_ = s.repository.WriteTransitionLog(ctx, cltypes.StateTransitionLog{
		ContentType: cltypes.ContentTypeCarpool,
		ContentID:   id,
		FromStatus:  result.FromStatus,
		ToStatus:    result.ToStatus,
		Action:      cltypes.ActionDelete,
		ActorUserID: principal.UserID,
	})
	s.recordAudit(ctx, principal, "campus_life.carpool.delete", "carpool", id, "拼车取消发布成功", nil)
	return nil
}
