package repo

import (
	"context"
	"strings"

	"github.com/liangluo/weouc2026/services/api-server/internal/platform/audit"
)

type AuditLogListQuery struct {
	ActorID      string
	Action       string
	ResourceType string
	ResourceID   string
}

type Repository interface {
	ListAuditLogs(ctx context.Context, query AuditLogListQuery) ([]audit.Entry, error)
}

type AuditRepository struct {
	store audit.Repository
}

func NewAuditRepository(store audit.Repository) *AuditRepository {
	return &AuditRepository{store: store}
}

func (r *AuditRepository) ListAuditLogs(ctx context.Context, query AuditLogListQuery) ([]audit.Entry, error) {
	if r == nil || r.store == nil {
		return nil, nil
	}
	return r.store.List(ctx, audit.ListQuery{
		ActorID:      strings.TrimSpace(query.ActorID),
		Action:       strings.TrimSpace(query.Action),
		ResourceType: strings.TrimSpace(query.ResourceType),
		ResourceID:   strings.TrimSpace(query.ResourceID),
	})
}
