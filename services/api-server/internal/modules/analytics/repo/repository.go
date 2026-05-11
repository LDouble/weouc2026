package repo

import (
	"context"

	"github.com/liangluo/weouc2026/services/api-server/internal/platform/audit"
)

type Repository interface {
	ListAuditLogs(ctx context.Context) ([]audit.Entry, error)
}

type AuditRepository struct {
	store audit.Repository
}

func NewAuditRepository(store audit.Repository) *AuditRepository {
	return &AuditRepository{store: store}
}

func (r *AuditRepository) ListAuditLogs(ctx context.Context) ([]audit.Entry, error) {
	if r == nil || r.store == nil {
		return nil, nil
	}
	return r.store.List(ctx)
}
