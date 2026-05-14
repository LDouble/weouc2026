package types

import (
	"testing"

	"github.com/liangluo/weouc2026/services/api-server/internal/platform/auth"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/bmfs"
)

func moderatorPrincipal() auth.Principal {
	return auth.Principal{
		Authenticated: true,
		UserID:        "moderator-1",
		Permissions:   []string{"campus_life:moderate"},
	}
}

func ownerPrincipal() auth.Principal {
	return auth.Principal{
		Authenticated: true,
		UserID:        "owner-1",
	}
}

func viewerPrincipal() auth.Principal {
	return auth.Principal{
		Authenticated: true,
		UserID:        "viewer-1",
	}
}

func unauthenticatedPrincipal() auth.Principal {
	return auth.Principal{}
}

func moderatorActx(isOwner bool) bmfs.ActionContext {
	return bmfs.ActionContext{
		Principal: moderatorPrincipal(),
		IsOwner:   isOwner,
	}
}

func ownerActx() bmfs.ActionContext {
	return bmfs.ActionContext{
		Principal: ownerPrincipal(),
		IsOwner:   true,
	}
}

func viewerActx() bmfs.ActionContext {
	return bmfs.ActionContext{
		Principal: viewerPrincipal(),
		IsOwner:   false,
	}
}

func unauthenticatedActx() bmfs.ActionContext {
	return bmfs.ActionContext{
		Principal: unauthenticatedPrincipal(),
		IsOwner:   false,
	}
}

func acceptorActx() bmfs.ActionContext {
	return bmfs.ActionContext{
		Principal: viewerPrincipal(),
		IsOwner:   false,
		UserRole:  "acceptor",
	}
}

func participantActx() bmfs.ActionContext {
	return bmfs.ActionContext{
		Principal: viewerPrincipal(),
		IsOwner:   false,
		UserRole:  "participant",
	}
}

func TestErrandFullLifecycle(t *testing.T) {
	m := ErrandStateMachine()

	result, err := m.Execute(StatusReviewing, ActionReviewApprove, moderatorActx(false))
	if err != nil {
		t.Fatalf("review_approve: %v", err)
	}
	if result.ToStatus != StatusPublished {
		t.Fatalf("expected published, got %s", result.ToStatus)
	}

	result, err = m.Execute(StatusPublished, ActionAccept, viewerActx())
	if err != nil {
		t.Fatalf("accept: %v", err)
	}
	if result.ToStatus != StatusAccepted {
		t.Fatalf("expected accepted, got %s", result.ToStatus)
	}

	result, err = m.Execute(StatusAccepted, ActionCancelAccept, acceptorActx())
	if err != nil {
		t.Fatalf("cancel_accept: %v", err)
	}
	if result.ToStatus != StatusPublished {
		t.Fatalf("expected published, got %s", result.ToStatus)
	}
}

func TestErrandCancelFromReviewing(t *testing.T) {
	m := ErrandStateMachine()

	result, err := m.Execute(StatusReviewing, ActionCancel, ownerActx())
	if err != nil {
		t.Fatalf("cancel: %v", err)
	}
	if result.ToStatus != StatusCancelled {
		t.Fatalf("expected cancelled, got %s", result.ToStatus)
	}
}

func TestErrandRejectFromReviewing(t *testing.T) {
	m := ErrandStateMachine()

	result, err := m.Execute(StatusReviewing, ActionReviewReject, moderatorActx(false))
	if err != nil {
		t.Fatalf("review_reject: %v", err)
	}
	if result.ToStatus != StatusRejected {
		t.Fatalf("expected rejected, got %s", result.ToStatus)
	}

	result, err = m.Execute(StatusRejected, ActionReviewReapprove, ownerActx())
	if err != nil {
		t.Fatalf("review_reapprove: %v", err)
	}
	if result.ToStatus != StatusReviewing {
		t.Fatalf("expected reviewing, got %s", result.ToStatus)
	}
}

func TestErrandGuardNotOwner(t *testing.T) {
	m := ErrandStateMachine()

	_, err := m.Execute(StatusReviewing, ActionCancel, viewerActx())
	if err == nil {
		t.Fatal("expected error when non-owner cancels from reviewing")
	}
}

func TestErrandGuardNotModerator(t *testing.T) {
	m := ErrandStateMachine()

	_, err := m.Execute(StatusReviewing, ActionReviewApprove, ownerActx())
	if err == nil {
		t.Fatal("expected error when non-moderator approves review")
	}
}

func TestErrandGuardNotAuthenticatedNotOwner(t *testing.T) {
	m := ErrandStateMachine()

	_, err := m.Execute(StatusPublished, ActionAccept, ownerActx())
	if err == nil {
		t.Fatal("expected error when owner tries to accept own errand")
	}
}

func TestErrandAvailableActions(t *testing.T) {
	m := ErrandStateMachine()

	reviewingOwner := m.AvailableActions(StatusReviewing, ownerActx())
	if !reviewingOwner["can_cancel"] {
		t.Error("expected can_cancel=true for owner in reviewing")
	}
	if reviewingOwner["can_review_approve"] {
		t.Error("expected can_review_approve=false for owner (not moderator) in reviewing")
	}

	reviewingMod := m.AvailableActions(StatusReviewing, moderatorActx(false))
	if !reviewingMod["can_review_approve"] {
		t.Error("expected can_review_approve=true for moderator in reviewing")
	}
	if !reviewingMod["can_review_reject"] {
		t.Error("expected can_review_reject=true for moderator in reviewing")
	}
	if reviewingMod["can_cancel"] {
		t.Error("expected can_cancel=false for moderator (not owner) in reviewing")
	}

	publishedViewer := m.AvailableActions(StatusPublished, viewerActx())
	if !publishedViewer["can_accept"] {
		t.Error("expected can_accept=true for authenticated non-owner in published")
	}
	if publishedViewer["can_cancel"] {
		t.Error("expected can_cancel=false for non-owner in published")
	}

	publishedOwner := m.AvailableActions(StatusPublished, ownerActx())
	if publishedOwner["can_accept"] {
		t.Error("expected can_accept=false for owner in published")
	}
	if !publishedOwner["can_cancel"] {
		t.Error("expected can_cancel=true for owner in published")
	}

	acceptedAcceptor := m.AvailableActions(StatusAccepted, acceptorActx())
	if !acceptedAcceptor["can_cancel_accept"] {
		t.Error("expected can_cancel_accept=true for acceptor in accepted")
	}

	acceptedViewer := m.AvailableActions(StatusAccepted, viewerActx())
	if acceptedViewer["can_cancel_accept"] {
		t.Error("expected can_cancel_accept=false for non-acceptor in accepted")
	}
}

func TestMeetupFullLifecycle(t *testing.T) {
	m := MeetupStateMachine()

	result, err := m.Execute(StatusReviewing, ActionReviewApprove, moderatorActx(false))
	if err != nil {
		t.Fatalf("review_approve: %v", err)
	}
	if result.ToStatus != StatusOpen {
		t.Fatalf("expected open, got %s", result.ToStatus)
	}

	result, err = m.Execute(StatusOpen, ActionJoin, viewerActx())
	if err != nil {
		t.Fatalf("join: %v", err)
	}
	if result.ToStatus != StatusOpen {
		t.Fatalf("expected open after join, got %s", result.ToStatus)
	}

	result, err = m.Execute(StatusFull, ActionCancelJoin, participantActx())
	if err != nil {
		t.Fatalf("cancel_join: %v", err)
	}
	if result.ToStatus != StatusOpen {
		t.Fatalf("expected open after cancel_join, got %s", result.ToStatus)
	}
}

func TestMeetupCancelFromReviewing(t *testing.T) {
	m := MeetupStateMachine()

	result, err := m.Execute(StatusReviewing, ActionCancel, ownerActx())
	if err != nil {
		t.Fatalf("cancel: %v", err)
	}
	if result.ToStatus != StatusCancelled {
		t.Fatalf("expected cancelled, got %s", result.ToStatus)
	}
}

func TestMeetupGuardNotOwner(t *testing.T) {
	m := MeetupStateMachine()

	_, err := m.Execute(StatusReviewing, ActionCancel, viewerActx())
	if err == nil {
		t.Fatal("expected error when non-owner cancels meetup from reviewing")
	}
}

func TestMeetupGuardNotModerator(t *testing.T) {
	m := MeetupStateMachine()

	_, err := m.Execute(StatusReviewing, ActionReviewApprove, ownerActx())
	if err == nil {
		t.Fatal("expected error when non-moderator approves meetup review")
	}
}

func TestMeetupJoinGuardNotAuthenticated(t *testing.T) {
	m := MeetupStateMachine()

	_, err := m.Execute(StatusOpen, ActionJoin, unauthenticatedActx())
	if err == nil {
		t.Fatal("expected error when unauthenticated user joins meetup")
	}
}

func TestMeetupCancelJoinGuardNotParticipant(t *testing.T) {
	m := MeetupStateMachine()

	_, err := m.Execute(StatusFull, ActionCancelJoin, viewerActx())
	if err == nil {
		t.Fatal("expected error when non-participant cancels join")
	}
}

func TestMarketFullLifecycle(t *testing.T) {
	m := MarketStateMachine()

	result, err := m.Execute(StatusReviewing, ActionReviewApprove, moderatorActx(false))
	if err != nil {
		t.Fatalf("review_approve: %v", err)
	}
	if result.ToStatus != StatusPublished {
		t.Fatalf("expected published, got %s", result.ToStatus)
	}

	result, err = m.Execute(StatusPublished, ActionDelete, ownerActx())
	if err != nil {
		t.Fatalf("delete: %v", err)
	}
	if result.ToStatus != StatusOffline {
		t.Fatalf("expected offline, got %s", result.ToStatus)
	}
}

func TestMarketRejectAndReapprove(t *testing.T) {
	m := MarketStateMachine()

	result, err := m.Execute(StatusReviewing, ActionReviewReject, moderatorActx(false))
	if err != nil {
		t.Fatalf("review_reject: %v", err)
	}
	if result.ToStatus != StatusRejected {
		t.Fatalf("expected rejected, got %s", result.ToStatus)
	}

	result, err = m.Execute(StatusRejected, ActionReviewReapprove, ownerActx())
	if err != nil {
		t.Fatalf("review_reapprove: %v", err)
	}
	if result.ToStatus != StatusReviewing {
		t.Fatalf("expected reviewing, got %s", result.ToStatus)
	}
}

func TestLostFoundMarkResolved(t *testing.T) {
	m := LostFoundStateMachine()

	result, err := m.Execute(StatusPublished, ActionMarkResolved, ownerActx())
	if err != nil {
		t.Fatalf("mark_resolved: %v", err)
	}
	if result.ToStatus != StatusResolved {
		t.Fatalf("expected resolved, got %s", result.ToStatus)
	}
}

func TestCarpoolDeleteFromPublished(t *testing.T) {
	m := CarpoolStateMachine()

	result, err := m.Execute(StatusPublished, ActionDelete, ownerActx())
	if err != nil {
		t.Fatalf("delete: %v", err)
	}
	if result.ToStatus != StatusOffline {
		t.Fatalf("expected offline, got %s", result.ToStatus)
	}
}

func TestResourceDeleteFromPublished(t *testing.T) {
	m := ResourceStateMachine()

	result, err := m.Execute(StatusPublished, ActionDelete, ownerActx())
	if err != nil {
		t.Fatalf("delete: %v", err)
	}
	if result.ToStatus != StatusOffline {
		t.Fatalf("expected offline, got %s", result.ToStatus)
	}
}

func TestGuardIsModerator(t *testing.T) {
	modActx := bmfs.ActionContext{Principal: moderatorPrincipal(), IsOwner: false}
	userActx := bmfs.ActionContext{Principal: ownerPrincipal(), IsOwner: false}

	if !guardIsModerator(modActx) {
		t.Error("expected moderator guard to pass for moderator principal")
	}
	if guardIsModerator(userActx) {
		t.Error("expected moderator guard to fail for regular user")
	}
}

func TestGuardIsOwner(t *testing.T) {
	ownerCtx := bmfs.ActionContext{Principal: ownerPrincipal(), IsOwner: true}
	nonOwnerCtx := bmfs.ActionContext{Principal: viewerPrincipal(), IsOwner: false}

	if !guardIsOwner(ownerCtx) {
		t.Error("expected owner guard to pass for owner")
	}
	if guardIsOwner(nonOwnerCtx) {
		t.Error("expected owner guard to fail for non-owner")
	}
}

func TestGuardIsAuthenticatedNotOwner(t *testing.T) {
	authNonOwner := bmfs.ActionContext{Principal: viewerPrincipal(), IsOwner: false}
	authOwner := bmfs.ActionContext{Principal: ownerPrincipal(), IsOwner: true}
	unauth := bmfs.ActionContext{Principal: unauthenticatedPrincipal(), IsOwner: false}

	if !guardIsAuthenticatedNotOwner(authNonOwner) {
		t.Error("expected guard to pass for authenticated non-owner")
	}
	if guardIsAuthenticatedNotOwner(authOwner) {
		t.Error("expected guard to fail for owner")
	}
	if guardIsAuthenticatedNotOwner(unauth) {
		t.Error("expected guard to fail for unauthenticated user")
	}
}

func TestGetMachine(t *testing.T) {
	cases := []struct {
		contentType string
		machineName string
	}{
		{ContentTypeErrand, "errand"},
		{ContentTypeMeetup, "meetup"},
		{ContentTypeMarket, "market"},
		{ContentTypeLostFound, "lost_found"},
		{ContentTypeCarpool, "carpool"},
		{ContentTypeResource, "resource"},
	}

	for _, tc := range cases {
		m := GetMachine(tc.contentType)
		if m == nil {
			t.Errorf("GetMachine(%q) returned nil", tc.contentType)
			continue
		}
		if m.Name() != tc.machineName {
			t.Errorf("GetMachine(%q): expected name %q, got %q", tc.contentType, tc.machineName, m.Name())
		}
	}
}

func TestGetMachineUnknownType(t *testing.T) {
	m := GetMachine("unknown")
	if m != nil {
		t.Error("expected nil for unknown content type")
	}
}

func TestErrandCanActionsFromPublished(t *testing.T) {
	m := ErrandStateMachine()

	authNonOwner := m.AvailableActions(StatusPublished, viewerActx())
	if !authNonOwner["can_accept"] {
		t.Error("expected can_accept=true for authenticated non-owner")
	}
	if authNonOwner["can_cancel"] {
		t.Error("expected can_cancel=false for non-owner")
	}

	owner := m.AvailableActions(StatusPublished, ownerActx())
	if owner["can_accept"] {
		t.Error("expected can_accept=false for owner")
	}
	if !owner["can_cancel"] {
		t.Error("expected can_cancel=true for owner")
	}
}

func TestMeetupCanActionsFromOpen(t *testing.T) {
	m := MeetupStateMachine()

	authNonOwner := m.AvailableActions(StatusOpen, viewerActx())
	if !authNonOwner["can_join"] {
		t.Error("expected can_join=true for authenticated non-owner")
	}
	if authNonOwner["can_cancel"] {
		t.Error("expected can_cancel=false for non-owner")
	}

	owner := m.AvailableActions(StatusOpen, ownerActx())
	if owner["can_join"] {
		t.Error("expected can_join=false for owner")
	}
	if !owner["can_cancel"] {
		t.Error("expected can_cancel=true for owner")
	}
}

func TestMarketCanActionsFromPublished(t *testing.T) {
	m := MarketStateMachine()

	owner := m.AvailableActions(StatusPublished, ownerActx())
	if !owner["can_delete"] {
		t.Error("expected can_delete=true for owner")
	}

	mod := m.AvailableActions(StatusPublished, moderatorActx(false))
	if !mod["can_offline_by_admin"] {
		t.Error("expected can_offline_by_admin=true for moderator")
	}
	if mod["can_delete"] {
		t.Error("expected can_delete=false for moderator (not owner)")
	}
}
