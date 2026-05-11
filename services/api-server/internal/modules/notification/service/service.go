package service

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	notificationrepo "github.com/liangluo/weouc2026/services/api-server/internal/modules/notification/repo"
	notificationtypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/notification/types"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/audit"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/auth"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/httpx"
)

type Service struct {
	repository notificationrepo.Repository
	recorder   audit.Recorder
}

func New(repository notificationrepo.Repository, recorder audit.Recorder) *Service {
	if repository == nil {
		repository = notificationrepo.NewInMemoryRepository()
	}
	return &Service{repository: repository, recorder: recorder}
}

func (s *Service) ListMessages(
	ctx context.Context,
	principal auth.Principal,
	query notificationtypes.MessageQuery,
) (map[string]any, error) {
	items, err := s.repository.ListMessages(ctx)
	if err != nil {
		return nil, httpx.Internal("读取通知列表失败", err)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.After(items[j].CreatedAt) })

	filtered := make([]notificationtypes.MessageItem, 0, len(items))
	for _, item := range items {
		if !isVisibleToUser(item, principal.UserID) {
			continue
		}
		if query.Category != "" && query.Category != item.Category {
			continue
		}
		if query.UnreadOnly && isReadByUser(item, principal.UserID) {
			continue
		}
		filtered = append(filtered, item)
	}

	offset, limit := normalizePagination(query.Page, query.PageSize)
	paged := paginateMessages(filtered, offset, limit)
	return map[string]any{
		"list":     messagePayloads(paged, principal.UserID),
		"total":    len(filtered),
		"page":     normalizePage(query.Page),
		"pageSize": normalizePageSize(query.PageSize),
	}, nil
}

func (s *Service) GetUnreadCount(ctx context.Context, principal auth.Principal) (map[string]any, error) {
	items, err := s.repository.ListMessages(ctx)
	if err != nil {
		return nil, httpx.Internal("读取未读通知数量失败", err)
	}

	count := 0
	for _, item := range items {
		if isVisibleToUser(item, principal.UserID) && !isReadByUser(item, principal.UserID) {
			count++
		}
	}
	return map[string]any{"count": count}, nil
}

func (s *Service) PublishMessage(
	ctx context.Context,
	principal auth.Principal,
	request notificationtypes.PublishRequest,
) (map[string]any, error) {
	title := strings.TrimSpace(request.Title)
	content := strings.TrimSpace(request.Content)
	if title == "" || content == "" {
		return nil, httpx.BadRequest("title 和 content 为必填项", nil)
	}

	scope := firstNonEmpty(strings.TrimSpace(request.TargetScope), "all")
	if scope != "all" && scope != "users" {
		return nil, httpx.BadRequest("target_scope 仅支持 all/users", nil)
	}

	targetUserIDs := sanitizeValues(request.TargetUserIDs)
	if scope == "users" && len(targetUserIDs) == 0 {
		return nil, httpx.BadRequest("target_scope=users 时必须提供 target_user_ids", nil)
	}

	id, err := s.repository.NextID(ctx, "notification")
	if err != nil {
		return nil, httpx.Internal("生成通知 ID 失败", err)
	}

	item := notificationtypes.MessageItem{
		ID:              id,
		Title:           title,
		Content:         content,
		Category:        firstNonEmpty(strings.TrimSpace(request.Category), "system"),
		TargetScope:     scope,
		TargetUserIDs:   targetUserIDs,
		ActionURL:       strings.TrimSpace(request.ActionURL),
		PublisherUserID: principal.UserID,
		Publisher:       displayName(principal),
		CreatedAt:       time.Now().UTC(),
		ReadByUserIDs:   map[string]time.Time{},
	}
	if _, err := s.repository.SaveMessage(ctx, item); err != nil {
		return nil, httpx.Internal("保存通知失败", err)
	}

	audit.RecordBestEffort(ctx, s.recorder, audit.Entry{
		ActorID:      principal.UserID,
		ActorName:    principal.DisplayName,
		Action:       "notification.publish",
		ResourceType: "notification_message",
		ResourceID:   item.ID,
		Message:      "站内通知发布成功",
		Details: map[string]any{
			"category":     item.Category,
			"target_scope": item.TargetScope,
		},
	})
	return map[string]any{"id": item.ID}, nil
}

func (s *Service) MarkRead(ctx context.Context, principal auth.Principal, messageID string) error {
	item, err := s.repository.GetMessage(ctx, strings.TrimSpace(messageID))
	if errors.Is(err, notificationrepo.ErrNotFound) {
		return httpx.NotFound("通知不存在", nil)
	}
	if err != nil {
		return httpx.Internal("读取通知失败", err)
	}
	if !isVisibleToUser(item, principal.UserID) {
		return httpx.NotFound("通知不存在", nil)
	}
	if isReadByUser(item, principal.UserID) {
		return nil
	}

	_, err = s.repository.UpdateMessage(ctx, item.ID, func(current *notificationtypes.MessageItem) error {
		if current.ReadByUserIDs == nil {
			current.ReadByUserIDs = make(map[string]time.Time)
		}
		current.ReadByUserIDs[principal.UserID] = time.Now().UTC()
		return nil
	})
	if errors.Is(err, notificationrepo.ErrNotFound) {
		return httpx.NotFound("通知不存在", nil)
	}
	if err != nil {
		return httpx.Internal("更新通知已读状态失败", err)
	}

	audit.RecordBestEffort(ctx, s.recorder, audit.Entry{
		ActorID:      principal.UserID,
		ActorName:    principal.DisplayName,
		Action:       "notification.read",
		ResourceType: "notification_message",
		ResourceID:   item.ID,
		Message:      "通知标记为已读",
	})
	return nil
}

func messagePayloads(items []notificationtypes.MessageItem, userID string) []map[string]any {
	result := make([]map[string]any, 0, len(items))
	for _, item := range items {
		result = append(result, messagePayload(item, userID))
	}
	return result
}

func messagePayload(item notificationtypes.MessageItem, userID string) map[string]any {
	payload := map[string]any{
		"id":           item.ID,
		"title":        item.Title,
		"content":      item.Content,
		"category":     item.Category,
		"target_scope": item.TargetScope,
		"action_url":   item.ActionURL,
		"publisher":    item.Publisher,
		"created_at":   item.CreatedAt.Format(time.RFC3339),
		"read":         isReadByUser(item, userID),
	}
	if item.TargetScope == "users" {
		payload["target_user_ids"] = append([]string(nil), item.TargetUserIDs...)
	}
	if readAt, exists := item.ReadByUserIDs[userID]; exists {
		payload["read_at"] = readAt.Format(time.RFC3339)
	}
	return payload
}

func paginateMessages(items []notificationtypes.MessageItem, offset, limit int) []notificationtypes.MessageItem {
	if offset >= len(items) {
		return nil
	}
	end := offset + limit
	if end > len(items) {
		end = len(items)
	}
	result := make([]notificationtypes.MessageItem, 0, end-offset)
	for _, item := range items[offset:end] {
		result = append(result, item)
	}
	return result
}

func isVisibleToUser(item notificationtypes.MessageItem, userID string) bool {
	if item.TargetScope == "all" {
		return true
	}
	for _, targetUserID := range item.TargetUserIDs {
		if targetUserID == userID {
			return true
		}
	}
	return false
}

func isReadByUser(item notificationtypes.MessageItem, userID string) bool {
	if item.ReadByUserIDs == nil {
		return false
	}
	_, exists := item.ReadByUserIDs[userID]
	return exists
}

func sanitizeValues(values []string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func displayName(principal auth.Principal) string {
	if strings.TrimSpace(principal.DisplayName) != "" {
		return strings.TrimSpace(principal.DisplayName)
	}
	if strings.TrimSpace(principal.UserID) != "" {
		return strings.TrimSpace(principal.UserID)
	}
	return "系统管理员"
}

func normalizePage(page int) int {
	if page <= 0 {
		return 1
	}
	return page
}

func normalizePageSize(pageSize int) int {
	if pageSize <= 0 {
		return 20
	}
	return pageSize
}

func normalizePagination(page, pageSize int) (int, int) {
	page = normalizePage(page)
	pageSize = normalizePageSize(pageSize)
	return (page - 1) * pageSize, pageSize
}
