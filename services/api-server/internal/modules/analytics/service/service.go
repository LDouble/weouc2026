package service

import (
	"context"
	"sort"
	"strings"
	"time"

	analyticsrepo "github.com/liangluo/weouc2026/services/api-server/internal/modules/analytics/repo"
	analyticstypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/analytics/types"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/audit"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/httpx"
)

type Service struct {
	repository analyticsrepo.Repository
}

func New(repository analyticsrepo.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetDashboard(ctx context.Context) (map[string]any, error) {
	entries, err := s.repository.ListAuditLogs(ctx)
	if err != nil {
		return nil, httpx.Internal("读取审计看板失败", err)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].CreatedAt.After(entries[j].CreatedAt) })

	actionCounts := make(map[string]int)
	total := len(entries)
	loginCount := 0
	bindCount := 0
	publishCount := 0
	reviewCount := 0
	todayActions := 0
	now := time.Now().UTC()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	for _, entry := range entries {
		actionCounts[entry.Action]++
		switch {
		case entry.Action == "auth.login":
			loginCount++
		case entry.Action == "auth.bind_student":
			bindCount++
		case strings.HasSuffix(entry.Action, ".publish"):
			publishCount++
		case entry.Action == "campus_life.review.update":
			reviewCount++
		}
		if entry.CreatedAt.After(todayStart) || entry.CreatedAt.Equal(todayStart) {
			todayActions++
		}
	}

	breakdown := make([]map[string]any, 0, len(actionCounts))
	for action, count := range actionCounts {
		breakdown = append(breakdown, map[string]any{
			"action": action,
			"count":  count,
		})
	}
	sort.Slice(breakdown, func(i, j int) bool {
		left := breakdown[i]["count"].(int)
		right := breakdown[j]["count"].(int)
		if left != right {
			return left > right
		}
		return breakdown[i]["action"].(string) < breakdown[j]["action"].(string)
	})

	recent := entries
	if len(recent) > 10 {
		recent = recent[:10]
	}

	return map[string]any{
		"summary": map[string]any{
			"total_audit_logs": total,
			"login_count":      loginCount,
			"bind_count":       bindCount,
			"publish_count":    publishCount,
			"review_count":     reviewCount,
			"today_actions":    todayActions,
		},
		"action_breakdown": breakdown,
		"recent_logs":      toPayloads(recent),
	}, nil
}

func (s *Service) ListAuditLogs(ctx context.Context, query analyticstypes.AuditLogQuery) (map[string]any, error) {
	entries, err := s.repository.ListAuditLogs(ctx)
	if err != nil {
		return nil, httpx.Internal("读取审计日志失败", err)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].CreatedAt.After(entries[j].CreatedAt) })

	filtered := make([]audit.Entry, 0, len(entries))
	for _, entry := range entries {
		if query.ActorID != "" && entry.ActorID != query.ActorID {
			continue
		}
		if query.Action != "" && entry.Action != query.Action {
			continue
		}
		if query.ResourceType != "" && entry.ResourceType != query.ResourceType {
			continue
		}
		if query.ResourceID != "" && entry.ResourceID != query.ResourceID {
			continue
		}
		filtered = append(filtered, entry)
	}

	offset, limit := normalizePagination(query.Page, query.PageSize)
	paged := paginateEntries(filtered, offset, limit)
	return map[string]any{
		"list":     toPayloads(paged),
		"total":    len(filtered),
		"page":     normalizePage(query.Page),
		"pageSize": normalizePageSize(query.PageSize),
	}, nil
}

func toPayloads(entries []audit.Entry) []map[string]any {
	result := make([]map[string]any, 0, len(entries))
	for _, entry := range entries {
		result = append(result, map[string]any{
			"id":            entry.ID,
			"actor_id":      entry.ActorID,
			"actor_name":    entry.ActorName,
			"action":        entry.Action,
			"resource_type": entry.ResourceType,
			"resource_id":   entry.ResourceID,
			"result":        entry.Result,
			"message":       entry.Message,
			"details":       entry.Details,
			"created_at":    entry.CreatedAt.Format(time.RFC3339),
		})
	}
	return result
}

func paginateEntries(entries []audit.Entry, offset, limit int) []audit.Entry {
	if offset >= len(entries) {
		return nil
	}
	end := offset + limit
	if end > len(entries) {
		end = len(entries)
	}
	result := make([]audit.Entry, 0, end-offset)
	for _, entry := range entries[offset:end] {
		result = append(result, entry)
	}
	return result
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
