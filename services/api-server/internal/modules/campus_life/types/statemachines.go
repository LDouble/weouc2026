package types

import (
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/bmfs"
)

const (
	ActionReviewApprove   = "review_approve"
	ActionReviewReject    = "review_reject"
	ActionReviewReapprove = "review_reapprove"
	ActionCancel          = "cancel"
	ActionAccept          = "accept"
	ActionCancelAccept    = "cancel_accept"
	ActionJoin            = "join"
	ActionCancelJoin      = "cancel_join"
	ActionDelete          = "delete"
	ActionMarkResolved    = "mark_resolved"
	ActionOfflineByAdmin  = "offline_by_admin"
)

const campusLifeModeratePermission = "campus_life:moderate"

func guardIsModerator(actx bmfs.ActionContext) bool {
	return actx.Principal.HasPermission(campusLifeModeratePermission)
}

func guardIsOwner(actx bmfs.ActionContext) bool {
	return actx.IsOwner
}

func guardIsOwnerOrModerator(actx bmfs.ActionContext) bool {
	return actx.IsOwner || actx.Principal.HasPermission(campusLifeModeratePermission)
}

func guardIsAuthenticatedNotOwner(actx bmfs.ActionContext) bool {
	return actx.Principal.Authenticated && !actx.IsOwner
}

func ErrandStateMachine() *bmfs.Machine {
	m := bmfs.NewMachine("errand")
	m.AddTransition(StatusReviewing, ActionReviewApprove, StatusPublished, guardIsModerator, nil)
	m.AddTransition(StatusReviewing, ActionReviewReject, StatusRejected, guardIsModerator, nil)
	m.AddTransition(StatusReviewing, ActionCancel, StatusCancelled, guardIsOwner, nil)
	m.AddTransition(StatusPublished, ActionAccept, StatusAccepted, guardIsAuthenticatedNotOwner, nil)
	m.AddTransition(StatusPublished, ActionCancel, StatusCancelled, guardIsOwner, nil)
	m.AddTransition(StatusPublished, ActionOfflineByAdmin, StatusOffline, guardIsModerator, nil)
	m.AddTransition(StatusAccepted, ActionCancelAccept, StatusPublished, func(actx bmfs.ActionContext) bool {
		return actx.UserRole == "acceptor"
	}, nil)
	m.AddTransition(StatusAccepted, ActionOfflineByAdmin, StatusOffline, guardIsModerator, nil)
	m.AddTransition(StatusRejected, ActionCancel, StatusCancelled, guardIsOwner, nil)
	m.AddTransition(StatusRejected, ActionReviewReapprove, StatusReviewing, guardIsOwner, nil)
	return m
}

func MeetupStateMachine() *bmfs.Machine {
	m := bmfs.NewMachine("meetup")
	m.AddTransition(StatusReviewing, ActionReviewApprove, StatusOpen, guardIsModerator, nil)
	m.AddTransition(StatusReviewing, ActionReviewReject, StatusRejected, guardIsModerator, nil)
	m.AddTransition(StatusReviewing, ActionCancel, StatusCancelled, guardIsOwner, nil)
	m.AddTransition(StatusOpen, ActionJoin, StatusOpen, guardIsAuthenticatedNotOwner, nil)
	m.AddTransition(StatusOpen, ActionCancel, StatusCancelled, guardIsOwner, nil)
	m.AddTransition(StatusOpen, ActionOfflineByAdmin, StatusOffline, guardIsModerator, nil)
	m.AddTransition(StatusFull, ActionCancelJoin, StatusOpen, func(actx bmfs.ActionContext) bool {
		return actx.UserRole == "participant"
	}, nil)
	m.AddTransition(StatusFull, ActionCancel, StatusCancelled, guardIsOwner, nil)
	m.AddTransition(StatusFull, ActionOfflineByAdmin, StatusOffline, guardIsModerator, nil)
	m.AddTransition(StatusRejected, ActionCancel, StatusCancelled, guardIsOwner, nil)
	m.AddTransition(StatusRejected, ActionReviewReapprove, StatusReviewing, guardIsOwner, nil)
	return m
}

func MarketStateMachine() *bmfs.Machine {
	m := bmfs.NewMachine("market")
	m.AddTransition(StatusReviewing, ActionReviewApprove, StatusPublished, guardIsModerator, nil)
	m.AddTransition(StatusReviewing, ActionReviewReject, StatusRejected, guardIsModerator, nil)
	m.AddTransition(StatusReviewing, ActionDelete, StatusOffline, guardIsOwner, nil)
	m.AddTransition(StatusPublished, ActionDelete, StatusOffline, guardIsOwner, nil)
	m.AddTransition(StatusPublished, ActionOfflineByAdmin, StatusOffline, guardIsModerator, nil)
	m.AddTransition(StatusRejected, ActionDelete, StatusOffline, guardIsOwner, nil)
	m.AddTransition(StatusRejected, ActionReviewReapprove, StatusReviewing, guardIsOwner, nil)
	return m
}

func LostFoundStateMachine() *bmfs.Machine {
	m := bmfs.NewMachine("lost_found")
	m.AddTransition(StatusReviewing, ActionReviewApprove, StatusPublished, guardIsModerator, nil)
	m.AddTransition(StatusReviewing, ActionReviewReject, StatusRejected, guardIsModerator, nil)
	m.AddTransition(StatusReviewing, ActionDelete, StatusOffline, guardIsOwner, nil)
	m.AddTransition(StatusPublished, ActionMarkResolved, StatusResolved, guardIsOwner, nil)
	m.AddTransition(StatusPublished, ActionDelete, StatusOffline, guardIsOwner, nil)
	m.AddTransition(StatusPublished, ActionOfflineByAdmin, StatusOffline, guardIsModerator, nil)
	m.AddTransition(StatusRejected, ActionDelete, StatusOffline, guardIsOwner, nil)
	m.AddTransition(StatusRejected, ActionReviewReapprove, StatusReviewing, guardIsOwner, nil)
	return m
}

func CarpoolStateMachine() *bmfs.Machine {
	m := bmfs.NewMachine("carpool")
	m.AddTransition(StatusReviewing, ActionReviewApprove, StatusPublished, guardIsModerator, nil)
	m.AddTransition(StatusReviewing, ActionReviewReject, StatusRejected, guardIsModerator, nil)
	m.AddTransition(StatusReviewing, ActionDelete, StatusOffline, guardIsOwner, nil)
	m.AddTransition(StatusPublished, ActionDelete, StatusOffline, guardIsOwner, nil)
	m.AddTransition(StatusPublished, ActionOfflineByAdmin, StatusOffline, guardIsModerator, nil)
	m.AddTransition(StatusRejected, ActionDelete, StatusOffline, guardIsOwner, nil)
	m.AddTransition(StatusRejected, ActionReviewReapprove, StatusReviewing, guardIsOwner, nil)
	return m
}

func ResourceStateMachine() *bmfs.Machine {
	m := bmfs.NewMachine("resource")
	m.AddTransition(StatusReviewing, ActionReviewApprove, StatusPublished, guardIsModerator, nil)
	m.AddTransition(StatusReviewing, ActionReviewReject, StatusRejected, guardIsModerator, nil)
	m.AddTransition(StatusReviewing, ActionDelete, StatusOffline, guardIsOwner, nil)
	m.AddTransition(StatusPublished, ActionDelete, StatusOffline, guardIsOwner, nil)
	m.AddTransition(StatusPublished, ActionOfflineByAdmin, StatusOffline, guardIsModerator, nil)
	m.AddTransition(StatusRejected, ActionDelete, StatusOffline, guardIsOwner, nil)
	m.AddTransition(StatusRejected, ActionReviewReapprove, StatusReviewing, guardIsOwner, nil)
	return m
}

func GetMachine(contentType string) *bmfs.Machine {
	switch contentType {
	case ContentTypeErrand:
		return ErrandStateMachine()
	case ContentTypeMeetup:
		return MeetupStateMachine()
	case ContentTypeMarket:
		return MarketStateMachine()
	case ContentTypeLostFound:
		return LostFoundStateMachine()
	case ContentTypeCarpool:
		return CarpoolStateMachine()
	case ContentTypeResource:
		return ResourceStateMachine()
	default:
		return nil
	}
}
