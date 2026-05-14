package service

import (
	"context"
	"errors"
	"slices"
	"strings"
	"time"

	clrepo "github.com/liangluo/weouc2026/services/api-server/internal/modules/campus_life/repo"
	cltypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/campus_life/types"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/auth"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/bmfs"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/httpx"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Service) ListMeetups(ctx context.Context, principal auth.Principal, query cltypes.MeetupQuery) (map[string]any, error) {
	currentUserID, visibleStatuses, includeAllStatus := buildVisibilityFilter(principal)

	items, total, err := s.repository.ListByType(ctx, cltypes.ContentTypeMeetup, cltypes.ContentFilter{
		Pagination:        query.Pagination,
		Category:          query.Category,
		Keyword:           query.Keyword,
		CurrentUserID:     currentUserID,
		VisibleStatuses:   visibleStatuses,
		IncludeAllStatus:  includeAllStatus,
		ParticipantUserID: currentUserID,
	})
	if err != nil {
		return nil, httpx.Internal("读取组局列表失败", err)
	}

	list := make([]map[string]any, 0, len(items))
	now := time.Now().In(chinaLocation)
	machine := cltypes.MeetupStateMachine()
	for _, item := range items {
		mp, _ := unmarshalPayload[cltypes.MeetupPayload](item.TypePayload)
		payload := buildMeetupPayload(item, mp, principal, now)
		role := meetupUserRole(item, principal)
		isOwner := role == "publisher"
		actx := bmfs.ActionContext{Principal: principal, IsOwner: isOwner, UserRole: role, Now: now}
		canActions := machine.AvailableActions(item.Status, actx)
		payload["is_owner"] = isOwner
		payload["can_edit"] = canEditContent(isOwner, item.Status) && item.Status != cltypes.StatusCancelled
		payload["can_delete"] = canActions["can_cancel"]
		list = append(list, payload)
	}

	return listEnvelope(list, int(total), query.Pagination), nil
}

func (s *Service) GetMeetupDetail(ctx context.Context, principal auth.Principal, id string) (map[string]any, error) {
	item, err := s.repository.GetByID(ctx, id)
	if errors.Is(err, clrepo.ErrNotFound) {
		return nil, httpx.NotFound("组局不存在", nil)
	}
	if err != nil {
		return nil, httpx.Internal("读取组局详情失败", err)
	}
	if item.ContentType != cltypes.ContentTypeMeetup {
		return nil, httpx.NotFound("组局不存在", nil)
	}
	if err := ensureContentVisible(principal, item, "组局不存在"); err != nil {
		return nil, err
	}
	if !shouldExposeMeetupState(principal, item, "") {
		return nil, httpx.NotFound("组局不存在", nil)
	}

	mp, _ := unmarshalPayload[cltypes.MeetupPayload](item.TypePayload)
	now := time.Now().In(chinaLocation)
	payload := buildMeetupPayload(item, mp, principal, now)
	payload["can_view_contact"] = canViewContact(principal, item.PublisherUserID)
	role := meetupUserRole(item, principal)
	isOwner := role == "publisher"
	actx := bmfs.ActionContext{Principal: principal, IsOwner: isOwner, UserRole: role, Now: now}
	canActions := cltypes.MeetupStateMachine().AvailableActions(item.Status, actx)
	payload["can_edit"] = canEditContent(isOwner, item.Status) && item.Status != cltypes.StatusCancelled
	payload["can_delete"] = canActions["can_cancel"]
	return payload, nil
}

func (s *Service) PublishMeetup(ctx context.Context, principal auth.Principal, request cltypes.MeetupPublishRequest) (map[string]any, error) {
	if strings.TrimSpace(request.Title) == "" || strings.TrimSpace(request.Location) == "" || strings.TrimSpace(request.Contact) == "" {
		return nil, httpx.BadRequest("标题、地点和联系方式不能为空", nil)
	}
	if request.MaxParticipants <= 1 {
		return nil, httpx.BadRequest("max_participants 至少为 2", nil)
	}

	startAt, err := time.Parse(time.RFC3339, strings.TrimSpace(request.StartAt))
	if err != nil {
		return nil, httpx.BadRequest("start_at 必须为 RFC3339 时间", nil)
	}
	deadlineAt := startAt
	if strings.TrimSpace(request.DeadlineAt) != "" {
		deadlineAt, err = time.Parse(time.RFC3339, strings.TrimSpace(request.DeadlineAt))
		if err != nil {
			return nil, httpx.BadRequest("deadline_at 必须为 RFC3339 时间", nil)
		}
	}
	if deadlineAt.After(startAt) {
		return nil, httpx.BadRequest("deadline_at 不能晚于 start_at", nil)
	}

	mp := cltypes.MeetupPayload{
		Category:           strings.TrimSpace(request.Category),
		Location:           strings.TrimSpace(request.Location),
		StartAt:            startAt,
		DeadlineAt:         deadlineAt,
		MaxParticipants:    request.MaxParticipants,
		FeeText:            strings.TrimSpace(request.FeeText),
		ParticipantUserIDs: []string{},
	}
	mp = refreshMeetupPayloadStatus(mp)

	item := cltypes.CommunityContent{
		ContentType:     cltypes.ContentTypeMeetup,
		Title:           strings.TrimSpace(request.Title),
		Desc:            firstNonEmpty(strings.TrimSpace(request.Desc), strings.TrimSpace(request.Title)),
		Status:          cltypes.StatusReviewing,
		PublisherUserID: principal.UserID,
		Contact:         strings.TrimSpace(request.Contact),
		Tags:            sanitizeTags(request.Tags),
		TypePayload:     marshalPayload(mp),
		CreatedBy:       principal.UserID,
		UpdatedBy:       principal.UserID,
	}
	item, err = s.repository.Save(ctx, item)
	if err != nil {
		return nil, httpx.Internal("保存组局失败", err)
	}

	s.recordAudit(ctx, principal, "campus_life.meetup.publish", "meetup", item.ID.Hex(), "组局发布成功", map[string]any{
		"status":           item.Status,
		"category":         request.Category,
		"max_participants": request.MaxParticipants,
	})

	return map[string]any{"id": item.ID.Hex()}, nil
}

func (s *Service) JoinMeetup(ctx context.Context, principal auth.Principal, meetupID string) error {
	var execResult *bmfs.ExecuteResult
	_, err := s.repository.Update(ctx, meetupID, func(item *cltypes.CommunityContent) error {
		if item.ContentType != cltypes.ContentTypeMeetup {
			return httpx.NotFound("组局不存在", nil)
		}
		if item.PublisherUserID == principal.UserID {
			return httpx.BadRequest("不能报名自己发起的组局", nil)
		}
		isOwner := item.PublisherUserID == principal.UserID
		role := meetupUserRole(*item, principal)
		actx := bmfs.ActionContext{Principal: principal, IsOwner: isOwner, UserRole: role, Now: time.Now()}
		result, err := cltypes.MeetupStateMachine().Execute(item.Status, cltypes.ActionJoin, actx)
		if err != nil {
			switch item.Status {
			case cltypes.StatusReviewing:
				return httpx.BadRequest("该组局仍在审核中，暂不可报名", nil)
			case cltypes.StatusRejected:
				return httpx.BadRequest("该组局审核未通过，无法报名", nil)
			case cltypes.StatusOffline:
				return httpx.BadRequest("该组局已下线，无法报名", nil)
			case cltypes.StatusCancelled:
				return httpx.BadRequest("该组局已取消", nil)
			case cltypes.StatusFull:
				return httpx.BadRequest("该组局人数已满", nil)
			default:
				return httpx.BadRequest("该组局当前状态无法报名", nil)
			}
		}
		mp, _ := unmarshalPayload[cltypes.MeetupPayload](item.TypePayload)
		now := time.Now().UTC()
		if !mp.DeadlineAt.IsZero() && mp.DeadlineAt.Before(now) {
			return httpx.BadRequest("该组局报名已截止", nil)
		}
		if !mp.StartAt.IsZero() && mp.StartAt.Before(now) {
			return httpx.BadRequest("该组局已开始，无法再报名", nil)
		}
		if slices.Contains(mp.ParticipantUserIDs, principal.UserID) {
			return httpx.BadRequest("你已经报名过该组局", nil)
		}
		execResult = result
		mp.ParticipantUserIDs = append(mp.ParticipantUserIDs, principal.UserID)
		mp = refreshMeetupPayloadStatus(mp)
		item.TypePayload = marshalPayload(mp)
		if mp.MaxParticipants > 0 && len(mp.ParticipantUserIDs)+1 >= mp.MaxParticipants {
			item.Status = cltypes.StatusFull
		} else {
			item.Status = result.ToStatus
		}
		item.UpdatedBy = principal.UserID
		return nil
	})
	if err != nil {
		if err == clrepo.ErrNotFound {
			return httpx.NotFound("组局不存在", nil)
		}
		if isAppError(err) {
			return err
		}
		return httpx.Internal("报名组局失败", err)
	}

	if execResult != nil {
		_ = s.repository.WriteTransitionLog(ctx, cltypes.StateTransitionLog{
			ID:          primitive.NewObjectID(),
			ContentType: cltypes.ContentTypeMeetup,
			ContentID:   meetupID,
			FromStatus:  execResult.FromStatus,
			ToStatus:    cltypes.StatusOpen,
			Action:      execResult.Action,
			ActorUserID: principal.UserID,
			CreatedAt:   time.Now(),
		})
	}

	return nil
}

func (s *Service) CancelMeetupJoin(ctx context.Context, principal auth.Principal, meetupID string) error {
	var execResult *bmfs.ExecuteResult
	_, err := s.repository.Update(ctx, meetupID, func(item *cltypes.CommunityContent) error {
		if item.ContentType != cltypes.ContentTypeMeetup {
			return httpx.NotFound("组局不存在", nil)
		}
		mp, _ := unmarshalPayload[cltypes.MeetupPayload](item.TypePayload)
		index := slices.Index(mp.ParticipantUserIDs, principal.UserID)
		if index < 0 {
			return httpx.BadRequest("你尚未报名该组局", nil)
		}
		isOwner := item.PublisherUserID == principal.UserID
		role := meetupUserRole(*item, principal)
		actx := bmfs.ActionContext{Principal: principal, IsOwner: isOwner, UserRole: role, Now: time.Now()}
		result, err := cltypes.MeetupStateMachine().Execute(item.Status, cltypes.ActionCancelJoin, actx)
		if err != nil {
			return httpx.BadRequest("该组局当前状态无法取消报名", nil)
		}
		execResult = result
		mp.ParticipantUserIDs = append(mp.ParticipantUserIDs[:index], mp.ParticipantUserIDs[index+1:]...)
		mp = refreshMeetupPayloadStatus(mp)
		item.TypePayload = marshalPayload(mp)
		item.Status = result.ToStatus
		item.UpdatedBy = principal.UserID
		return nil
	})
	if err != nil {
		if err == clrepo.ErrNotFound {
			return httpx.NotFound("组局不存在", nil)
		}
		if isAppError(err) {
			return err
		}
		return httpx.Internal("取消报名组局失败", err)
	}

	if execResult != nil {
		_ = s.repository.WriteTransitionLog(ctx, cltypes.StateTransitionLog{
			ID:          primitive.NewObjectID(),
			ContentType: cltypes.ContentTypeMeetup,
			ContentID:   meetupID,
			FromStatus:  execResult.FromStatus,
			ToStatus:    execResult.ToStatus,
			Action:      execResult.Action,
			ActorUserID: principal.UserID,
			CreatedAt:   time.Now(),
		})
	}

	return nil
}

func (s *Service) CancelMeetupPublish(ctx context.Context, principal auth.Principal, meetupID string) error {
	var execResult *bmfs.ExecuteResult
	_, err := s.repository.Update(ctx, meetupID, func(item *cltypes.CommunityContent) error {
		if item.ContentType != cltypes.ContentTypeMeetup {
			return httpx.NotFound("组局不存在", nil)
		}
		if item.PublisherUserID != principal.UserID {
			return httpx.Forbidden("只有发起人可以取消组局", nil)
		}
		isOwner := item.PublisherUserID == principal.UserID
		role := meetupUserRole(*item, principal)
		actx := bmfs.ActionContext{Principal: principal, IsOwner: isOwner, UserRole: role, Now: time.Now()}
		result, err := cltypes.MeetupStateMachine().Execute(item.Status, cltypes.ActionCancel, actx)
		if err != nil {
			return httpx.BadRequest("该组局当前状态无法取消", nil)
		}
		execResult = result
		item.Status = result.ToStatus
		item.UpdatedBy = principal.UserID
		return nil
	})
	if err != nil {
		if err == clrepo.ErrNotFound {
			return httpx.NotFound("组局不存在", nil)
		}
		if isAppError(err) {
			return err
		}
		return httpx.Internal("取消组局失败", err)
	}

	if execResult != nil {
		_ = s.repository.WriteTransitionLog(ctx, cltypes.StateTransitionLog{
			ID:          primitive.NewObjectID(),
			ContentType: cltypes.ContentTypeMeetup,
			ContentID:   meetupID,
			FromStatus:  execResult.FromStatus,
			ToStatus:    execResult.ToStatus,
			Action:      execResult.Action,
			ActorUserID: principal.UserID,
			CreatedAt:   time.Now(),
		})
	}

	return nil
}
