package service

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	portalrepo "github.com/liangluo/weouc2026/services/api-server/internal/modules/portal/repo"
	portaltypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/portal/types"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/auth"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/httpx"
)

type Service struct {
	repository portalrepo.Repository
}

func New(repository portalrepo.Repository) *Service {
	if repository == nil {
		repository = portalrepo.NewInMemoryRepository()
	}
	return &Service{repository: repository}
}

func (s *Service) GetHome(ctx context.Context) (map[string]any, error) {
	banners, err := s.repository.ListBanners(ctx)
	if err != nil {
		return nil, httpx.Internal("读取门户轮播失败", err)
	}
	sort.Slice(banners, func(i, j int) bool { return banners[i].Sort < banners[j].Sort })

	notices, err := s.repository.ListNotices(ctx)
	if err != nil {
		return nil, httpx.Internal("读取门户公告失败", err)
	}
	sortNotices(notices)
	if len(notices) > 5 {
		notices = notices[:5]
	}

	return map[string]any{
		"banners": bannerPayloads(banners),
		"notices": noticePayloads(notices, false),
	}, nil
}

func (s *Service) ListNotices(ctx context.Context, query portaltypes.NoticeQuery) (map[string]any, error) {
	notices, err := s.repository.ListNotices(ctx)
	if err != nil {
		return nil, httpx.Internal("读取公告列表失败", err)
	}
	sortNotices(notices)

	filtered := make([]portaltypes.NoticeItem, 0, len(notices))
	for _, item := range notices {
		if !matchKeyword(query.Keyword, item.Title, item.Summary, item.Content, item.Publisher) {
			continue
		}
		filtered = append(filtered, item)
	}

	offset, limit := normalizePagination(query.Page, query.PageSize)
	paged := paginateNotices(filtered, offset, limit)
	return map[string]any{
		"list":     noticePayloads(paged, false),
		"total":    len(filtered),
		"page":     normalizePage(query.Page),
		"pageSize": normalizePageSize(query.PageSize),
	}, nil
}

func (s *Service) GetNotice(ctx context.Context, id string) (map[string]any, error) {
	item, err := s.repository.GetNotice(ctx, strings.TrimSpace(id))
	if errors.Is(err, portalrepo.ErrNotFound) {
		return nil, httpx.NotFound("公告不存在", nil)
	}
	if err != nil {
		return nil, httpx.Internal("读取公告详情失败", err)
	}
	return noticePayload(item, true), nil
}

func (s *Service) PublishNotice(
	ctx context.Context,
	principal auth.Principal,
	request portaltypes.NoticePublishRequest,
) (map[string]any, error) {
	title := strings.TrimSpace(request.Title)
	content := strings.TrimSpace(request.Content)
	if title == "" || content == "" {
		return nil, httpx.BadRequest("title 和 content 为必填项", nil)
	}

	id, err := s.repository.NextID(ctx, "notice")
	if err != nil {
		return nil, httpx.Internal("生成公告 ID 失败", err)
	}

	now := time.Now().UTC()
	item := portaltypes.NoticeItem{
		ID:              id,
		Title:           title,
		Summary:         firstNonEmpty(strings.TrimSpace(request.Summary), summarize(content)),
		Content:         content,
		Audience:        firstNonEmpty(strings.TrimSpace(request.Audience), "all"),
		Tags:            sanitizeTags(request.Tags),
		Pinned:          request.Pinned,
		PublisherUserID: principal.UserID,
		Publisher:       displayName(principal),
		PublishedAt:     now,
		CreatedAt:       now,
	}
	if _, err := s.repository.SaveNotice(ctx, item); err != nil {
		return nil, httpx.Internal("保存公告失败", err)
	}
	return map[string]any{"id": item.ID}, nil
}

func bannerPayloads(items []portaltypes.BannerItem) []map[string]any {
	result := make([]map[string]any, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]any{
			"id":          item.ID,
			"title":       item.Title,
			"description": item.Description,
			"image_url":   item.ImageURL,
			"action_url":  item.ActionURL,
			"sort":        item.Sort,
			"created_at":  item.CreatedAt.Format(time.RFC3339),
		})
	}
	return result
}

func noticePayloads(items []portaltypes.NoticeItem, includeContent bool) []map[string]any {
	result := make([]map[string]any, 0, len(items))
	for _, item := range items {
		result = append(result, noticePayload(item, includeContent))
	}
	return result
}

func noticePayload(item portaltypes.NoticeItem, includeContent bool) map[string]any {
	payload := map[string]any{
		"id":           item.ID,
		"title":        item.Title,
		"summary":      item.Summary,
		"audience":     item.Audience,
		"tags":         append([]string(nil), item.Tags...),
		"pinned":       item.Pinned,
		"publisher":    item.Publisher,
		"published_at": item.PublishedAt.Format(time.RFC3339),
		"created_at":   item.CreatedAt.Format(time.RFC3339),
	}
	if includeContent {
		payload["content"] = item.Content
	}
	return payload
}

func paginateNotices(items []portaltypes.NoticeItem, offset, limit int) []portaltypes.NoticeItem {
	if offset >= len(items) {
		return nil
	}
	end := offset + limit
	if end > len(items) {
		end = len(items)
	}
	result := make([]portaltypes.NoticeItem, 0, end-offset)
	for _, item := range items[offset:end] {
		result = append(result, item)
	}
	return result
}

func sortNotices(items []portaltypes.NoticeItem) {
	sort.Slice(items, func(i, j int) bool {
		if items[i].Pinned != items[j].Pinned {
			return items[i].Pinned
		}
		return items[i].PublishedAt.After(items[j].PublishedAt)
	})
}

func matchKeyword(keyword string, values ...string) bool {
	keyword = strings.TrimSpace(strings.ToLower(keyword))
	if keyword == "" {
		return true
	}
	for _, value := range values {
		if strings.Contains(strings.ToLower(value), keyword) {
			return true
		}
	}
	return false
}

func summarize(content string) string {
	runes := []rune(strings.TrimSpace(content))
	if len(runes) <= 48 {
		return string(runes)
	}
	return string(runes[:48]) + "..."
}

func sanitizeTags(tags []string) []string {
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		value := strings.TrimSpace(tag)
		if value != "" {
			result = append(result, value)
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
