package repo

import (
	"context"
	"testing"

	"github.com/liangluo/weouc2026/services/api-server/internal/modules/system/types"
)

type fakeProbe struct {
	status types.DependencyStatus
}

func (p fakeProbe) Check(context.Context) types.DependencyStatus {
	return p.status
}

func TestRuntimeStatusRepositoryMarksServiceUnavailableWhenRequiredDependencyFails(t *testing.T) {
	repository := &RuntimeStatusRepository{
		probes: []DependencyProbe{
			fakeProbe{status: types.DependencyStatus{Name: "postgres", Status: "ready", Required: true}},
			fakeProbe{status: types.DependencyStatus{Name: "redis", Status: "not_ready", Required: true}},
			fakeProbe{status: types.DependencyStatus{Name: "object_storage", Status: "skipped", Required: false}},
		},
	}

	status := repository.ReadinessSnapshot(context.Background())
	if status.IsReady() {
		t.Fatalf("expected readiness to be false, got %+v", status)
	}
	if status.Status != "not_ready" {
		t.Fatalf("expected status not_ready, got %q", status.Status)
	}
}

func TestRuntimeStatusRepositoryAllowsSkippedOptionalDependencies(t *testing.T) {
	repository := &RuntimeStatusRepository{
		probes: []DependencyProbe{
			fakeProbe{status: types.DependencyStatus{Name: "postgres", Status: "skipped", Required: false}},
			fakeProbe{status: types.DependencyStatus{Name: "redis", Status: "skipped", Required: false}},
		},
	}

	status := repository.ReadinessSnapshot(context.Background())
	if !status.IsReady() {
		t.Fatalf("expected readiness to stay true, got %+v", status)
	}
}
