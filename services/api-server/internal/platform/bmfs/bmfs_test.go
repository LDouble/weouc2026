package bmfs

import (
	"testing"
	"time"

	"github.com/liangluo/weouc2026/services/api-server/internal/platform/auth"
)

func ownerOnlyGuard(ctx ActionContext) bool {
	return ctx.IsOwner
}

func adminOnlyGuard(ctx ActionContext) bool {
	for _, r := range ctx.Principal.Roles {
		if r == "admin" {
			return true
		}
	}
	return false
}

func buildTestMachine() *Machine {
	return NewMachine("post").
		AddTransition("reviewing", "review_approve", "published", nil, nil).
		AddTransition("reviewing", "review_reject", "rejected", nil, nil).
		AddTransition("published", "unpublish", "draft", ownerOnlyGuard, nil).
		AddTransition("draft", "submit", "reviewing", nil, nil).
		AddTransition("rejected", "resubmit", "reviewing", nil, nil)
}

func TestNormalTransition(t *testing.T) {
	m := buildTestMachine()
	actx := ActionContext{
		Principal: auth.Principal{Authenticated: true, UserID: "u1"},
		IsOwner:   true,
		Now:       time.Now(),
	}

	result, err := m.Execute("reviewing", "review_approve", actx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.FromStatus != "reviewing" {
		t.Errorf("expected FromStatus reviewing, got %s", result.FromStatus)
	}
	if result.ToStatus != "published" {
		t.Errorf("expected ToStatus published, got %s", result.ToStatus)
	}
	if result.Action != "review_approve" {
		t.Errorf("expected Action review_approve, got %s", result.Action)
	}
}

func TestGuardRejection(t *testing.T) {
	m := buildTestMachine()
	actx := ActionContext{
		Principal: auth.Principal{Authenticated: true, UserID: "u1"},
		IsOwner:   false,
		Now:       time.Now(),
	}

	_, err := m.Execute("published", "unpublish", actx)
	if err == nil {
		t.Fatal("expected error for guard rejection, got nil")
	}
}

func TestGuardAllowsOwner(t *testing.T) {
	m := buildTestMachine()
	actx := ActionContext{
		Principal: auth.Principal{Authenticated: true, UserID: "u1"},
		IsOwner:   true,
		Now:       time.Now(),
	}

	result, err := m.Execute("published", "unpublish", actx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ToStatus != "draft" {
		t.Errorf("expected ToStatus draft, got %s", result.ToStatus)
	}
}

func TestAvailableActions(t *testing.T) {
	m := buildTestMachine()
	actx := ActionContext{
		Principal: auth.Principal{Authenticated: true, UserID: "u1"},
		IsOwner:   false,
		Now:       time.Now(),
	}

	actions := m.AvailableActions("reviewing", actx)
	if len(actions) != 2 {
		t.Fatalf("expected 2 actions, got %d", len(actions))
	}
	if !actions["can_review_approve"] {
		t.Error("expected can_review_approve to be true")
	}
	if !actions["can_review_reject"] {
		t.Error("expected can_review_reject to be true")
	}
}

func TestAvailableActionsWithGuard(t *testing.T) {
	m := buildTestMachine()
	actx := ActionContext{
		Principal: auth.Principal{Authenticated: true, UserID: "u1"},
		IsOwner:   false,
		Now:       time.Now(),
	}

	actions := m.AvailableActions("published", actx)
	if len(actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(actions))
	}
	if actions["can_unpublish"] {
		t.Error("expected can_unpublish to be false for non-owner")
	}
}

func TestAvailableActionsOwnerWithGuard(t *testing.T) {
	m := buildTestMachine()
	actx := ActionContext{
		Principal: auth.Principal{Authenticated: true, UserID: "u1"},
		IsOwner:   true,
		Now:       time.Now(),
	}

	actions := m.AvailableActions("published", actx)
	if len(actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(actions))
	}
	if !actions["can_unpublish"] {
		t.Error("expected can_unpublish to be true for owner")
	}
}

func TestInvalidActionFromCurrentState(t *testing.T) {
	m := buildTestMachine()
	actx := ActionContext{
		Principal: auth.Principal{Authenticated: true, UserID: "u1"},
		IsOwner:   true,
		Now:       time.Now(),
	}

	_, err := m.Execute("draft", "review_approve", actx)
	if err == nil {
		t.Fatal("expected error for invalid action from current state, got nil")
	}
}

func TestUnknownStatus(t *testing.T) {
	m := buildTestMachine()
	actx := ActionContext{
		Principal: auth.Principal{Authenticated: true, UserID: "u1"},
		IsOwner:   true,
		Now:       time.Now(),
	}

	_, err := m.Execute("nonexistent", "review_approve", actx)
	if err == nil {
		t.Fatal("expected error for unknown status, got nil")
	}
}

func TestMultipleTransitionsFromSameState(t *testing.T) {
	m := buildTestMachine()
	actx := ActionContext{
		Principal: auth.Principal{Authenticated: true, UserID: "u1"},
		IsOwner:   true,
		Now:       time.Now(),
	}

	result, err := m.Execute("reviewing", "review_reject", actx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ToStatus != "rejected" {
		t.Errorf("expected ToStatus rejected, got %s", result.ToStatus)
	}
}

func TestOnTransitionCallback(t *testing.T) {
	var capturedFrom, capturedTo, capturedAction string
	called := false

	m := NewMachine("order").
		AddTransition("created", "pay", "paid", nil, func(ctx ActionContext, from, to, action string) {
			called = true
			capturedFrom = from
			capturedTo = to
			capturedAction = action
		})

	actx := ActionContext{
		Principal: auth.Principal{Authenticated: true, UserID: "u1"},
		IsOwner:   true,
		Now:       time.Now(),
	}

	result, err := m.Execute("created", "pay", actx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !called {
		t.Fatal("expected OnTransition callback to be called")
	}
	if capturedFrom != "created" {
		t.Errorf("expected from created, got %s", capturedFrom)
	}
	if capturedTo != "paid" {
		t.Errorf("expected to paid, got %s", capturedTo)
	}
	if capturedAction != "pay" {
		t.Errorf("expected action pay, got %s", capturedAction)
	}
	if result.ToStatus != "paid" {
		t.Errorf("expected ToStatus paid, got %s", result.ToStatus)
	}
}

func TestExecuteResultCanActions(t *testing.T) {
	m := buildTestMachine()
	actx := ActionContext{
		Principal: auth.Principal{Authenticated: true, UserID: "u1"},
		IsOwner:   false,
		Now:       time.Now(),
	}

	result, err := m.Execute("reviewing", "review_approve", actx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.CanActions) != 1 {
		t.Fatalf("expected 1 can_action for published state, got %d", len(result.CanActions))
	}
	if result.CanActions["can_unpublish"] {
		t.Error("expected can_unpublish to be false for non-owner in published state")
	}
}

func TestMachineName(t *testing.T) {
	m := NewMachine("errand")
	if m.Name() != "errand" {
		t.Errorf("expected name errand, got %s", m.Name())
	}
}

func TestAddStateExplicitly(t *testing.T) {
	m := NewMachine("test").
		AddState("initial").
		AddTransition("initial", "go", "final", nil, nil)

	actx := ActionContext{
		Principal: auth.Principal{Authenticated: true, UserID: "u1"},
		Now:       time.Now(),
	}

	result, err := m.Execute("initial", "go", actx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ToStatus != "final" {
		t.Errorf("expected ToStatus final, got %s", result.ToStatus)
	}
}

func TestAdminGuardTransition(t *testing.T) {
	m := NewMachine("moderation").
		AddTransition("flagged", "dismiss", "published", adminOnlyGuard, nil)

	adminCtx := ActionContext{
		Principal: auth.Principal{Authenticated: true, UserID: "admin1", Roles: []string{"admin"}},
		Now:       time.Now(),
	}
	userCtx := ActionContext{
		Principal: auth.Principal{Authenticated: true, UserID: "user1", Roles: []string{"student"}},
		Now:       time.Now(),
	}

	result, err := m.Execute("flagged", "dismiss", adminCtx)
	if err != nil {
		t.Fatalf("expected admin to dismiss, got error: %v", err)
	}
	if result.ToStatus != "published" {
		t.Errorf("expected ToStatus published, got %s", result.ToStatus)
	}

	_, err = m.Execute("flagged", "dismiss", userCtx)
	if err == nil {
		t.Fatal("expected non-admin to be blocked by guard")
	}
}

func TestFullLifecycle(t *testing.T) {
	m := buildTestMachine()
	actx := ActionContext{
		Principal: auth.Principal{Authenticated: true, UserID: "u1"},
		IsOwner:   true,
		Now:       time.Now(),
	}

	steps := []struct {
		status string
		action string
		want   string
	}{
		{"draft", "submit", "reviewing"},
		{"reviewing", "review_reject", "rejected"},
		{"rejected", "resubmit", "reviewing"},
		{"reviewing", "review_approve", "published"},
		{"published", "unpublish", "draft"},
	}

	for _, step := range steps {
		result, err := m.Execute(step.status, step.action, actx)
		if err != nil {
			t.Fatalf("Execute(%q, %q) failed: %v", step.status, step.action, err)
		}
		if result.ToStatus != step.want {
			t.Errorf("Execute(%q, %q): expected ToStatus %q, got %q", step.status, step.action, step.want, result.ToStatus)
		}
	}
}
