package repo

import (
	"fmt"
	"strings"
	"time"

	cltypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/campus_life/types"
)

type sqlConditionBuilder struct {
	conditions []string
	args       []any
}

func (b *sqlConditionBuilder) addExact(column, value string) {
	value = strings.TrimSpace(value)
	if value == "" {
		return
	}
	b.args = append(b.args, value)
	b.conditions = append(b.conditions, fmt.Sprintf("%s = $%d", column, len(b.args)))
}

func (b *sqlConditionBuilder) addKeyword(keyword string, expressions ...string) {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" || len(expressions) == 0 {
		return
	}
	b.args = append(b.args, "%"+keyword+"%")
	placeholder := fmt.Sprintf("$%d", len(b.args))
	parts := make([]string, 0, len(expressions))
	for _, expression := range expressions {
		parts = append(parts, expression+` ILIKE `+placeholder)
	}
	b.conditions = append(b.conditions, `(`+strings.Join(parts, ` OR `)+`)`)
}

func (b *sqlConditionBuilder) addIn(column string, values []string) {
	values = compactStrings(values)
	if len(values) == 0 {
		return
	}
	placeholders := make([]string, 0, len(values))
	for _, value := range values {
		b.args = append(b.args, value)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(b.args)))
	}
	b.conditions = append(b.conditions, fmt.Sprintf("%s IN (%s)", column, strings.Join(placeholders, `, `)))
}

func (b *sqlConditionBuilder) addTimeGTE(column string, value time.Time) {
	if value.IsZero() {
		return
	}
	b.args = append(b.args, value)
	b.conditions = append(b.conditions, fmt.Sprintf("%s >= $%d", column, len(b.args)))
}

func (b *sqlConditionBuilder) addTimeLT(column string, value time.Time) {
	if value.IsZero() {
		return
	}
	b.args = append(b.args, value)
	b.conditions = append(b.conditions, fmt.Sprintf("%s < $%d", column, len(b.args)))
}

func (b *sqlConditionBuilder) addVisibility(reviewStatusColumn, ownerColumn string, query ContentVisibilityQuery) {
	if query.IncludeAllReviewStatuses {
		return
	}
	statuses := compactStrings(query.ReviewStatuses)
	ownerUserID := strings.TrimSpace(query.IncludeOwnerUserID)
	if len(statuses) == 0 && ownerUserID == "" {
		return
	}
	if len(statuses) == 0 {
		b.addExact(ownerColumn, ownerUserID)
		return
	}
	if ownerUserID == "" {
		b.addIn(reviewStatusColumn, statuses)
		return
	}

	placeholders := make([]string, 0, len(statuses))
	for _, status := range statuses {
		b.args = append(b.args, status)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(b.args)))
	}
	b.args = append(b.args, ownerUserID)
	ownerPlaceholder := fmt.Sprintf("$%d", len(b.args))
	b.conditions = append(
		b.conditions,
		fmt.Sprintf("((%s IN (%s)) OR %s = %s)", reviewStatusColumn, strings.Join(placeholders, `, `), ownerColumn, ownerPlaceholder),
	)
}

func (b *sqlConditionBuilder) addMeetupState(query MeetupStateQuery) {
	if query.IncludeAllStatuses {
		return
	}
	userID := strings.TrimSpace(query.IncludeCancelledForUserID)
	if userID == "" {
		b.conditions = append(b.conditions, `status <> 'cancelled'`)
		return
	}
	b.args = append(b.args, userID)
	placeholder := fmt.Sprintf("$%d", len(b.args))
	b.conditions = append(
		b.conditions,
		`(status <> 'cancelled' OR publisher_user_id = `+placeholder+` OR participant_user_ids ? `+placeholder+`)`,
	)
}

func (b *sqlConditionBuilder) build(base string) (string, []any) {
	if len(b.conditions) == 0 {
		return base, b.args
	}
	return base + ` WHERE ` + strings.Join(b.conditions, ` AND `), b.args
}

func buildMarketListSQL(query MarketListQuery) (string, []any) {
	var builder sqlConditionBuilder
	builder.addVisibility("review_status", "publisher_user_id", query.Visibility)
	builder.addExact(`extra ->> 'category'`, query.Category)
	builder.addExact(`publisher_user_id`, query.PublisherUserID)
	builder.addKeyword(query.Keyword, `title`, `description`)
	return builder.build(`SELECT ` + marketColumns + ` FROM campus_markets`)
}

func buildErrandListSQL(query ErrandListQuery) (string, []any) {
	var builder sqlConditionBuilder
	builder.addVisibility("review_status", "publisher_user_id", query.Visibility)
	builder.addExact(`category`, query.Category)
	builder.addExact(`publisher_user_id`, query.PublisherUserID)
	builder.addExact(`acceptor_user_id`, query.AcceptorUserID)
	builder.addKeyword(query.Keyword, `title`, `description`)
	return builder.build(`SELECT ` + errandColumns + ` FROM campus_errands`)
}

func buildResourceListSQL(query ResourceListQuery) (string, []any) {
	var builder sqlConditionBuilder
	builder.addVisibility("review_status", "publisher_user_id", query.Visibility)
	builder.addExact(`extra ->> 'category'`, query.Category)
	builder.addExact(`publisher_user_id`, query.PublisherUserID)
	builder.addKeyword(query.Keyword, `title`, `description`)
	return builder.build(`SELECT ` + resourceColumns + ` FROM campus_resources`)
}

func buildLostFoundListSQL(query LostFoundListQuery) (string, []any) {
	var builder sqlConditionBuilder
	builder.addVisibility("review_status", "publisher_user_id", query.Visibility)
	builder.addExact(`extra ->> 'category'`, query.Category)
	builder.addExact(`extra ->> 'type'`, query.Type)
	builder.addExact(`publisher_user_id`, query.PublisherUserID)
	builder.addKeyword(query.Keyword, `title`, `description`)
	return builder.build(`SELECT ` + lostFoundColumns + ` FROM campus_lost_founds`)
}

func buildCarpoolListSQL(query CarpoolListQuery) (string, []any) {
	var builder sqlConditionBuilder
	builder.addVisibility("review_status", "publisher_user_id", query.Visibility)
	builder.addExact(`publisher_user_id`, query.PublisherUserID)
	builder.addTimeGTE(`travel_at`, query.TravelAtFrom)
	builder.addTimeLT(`travel_at`, query.TravelAtTo)
	builder.addKeyword(query.Keyword, `route_from`, `route_to`, `note`, `publisher`, `type_label`)
	return builder.build(`SELECT ` + carpoolColumns + ` FROM campus_carpools`)
}

func buildMeetupListSQL(query MeetupListQuery) (string, []any) {
	var builder sqlConditionBuilder
	builder.addVisibility("review_status", "publisher_user_id", query.Visibility)
	builder.addMeetupState(query.State)
	builder.addExact(`category`, query.Category)
	builder.addExact(`publisher_user_id`, query.PublisherUserID)
	if participantUserID := strings.TrimSpace(query.ParticipantUserID); participantUserID != "" {
		builder.args = append(builder.args, participantUserID)
		builder.conditions = append(builder.conditions, fmt.Sprintf(`participant_user_ids ? $%d`, len(builder.args)))
	}
	builder.addKeyword(query.Keyword, `title`, `description`, `location`, `publisher`)
	return builder.build(`SELECT ` + meetupColumns + ` FROM campus_meetups`)
}

func matchContentVisibility(reviewStatus, ownerUserID string, query ContentVisibilityQuery) bool {
	if query.IncludeAllReviewStatuses {
		return true
	}
	if ownerUserID != "" && strings.TrimSpace(query.IncludeOwnerUserID) == ownerUserID {
		return true
	}
	statuses := compactStrings(query.ReviewStatuses)
	if len(statuses) == 0 {
		return true
	}
	return containsFold(statuses, reviewStatus)
}

func matchKeywordQuery(keyword string, values ...string) bool {
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

func matchMeetupState(item cltypes.MeetupItem, query MeetupStateQuery) bool {
	if query.IncludeAllStatuses {
		return true
	}
	if normalizeMeetupState(item.Status) != "cancelled" {
		return true
	}
	userID := strings.TrimSpace(query.IncludeCancelledForUserID)
	if userID == "" {
		return false
	}
	if item.PublisherUserID == userID {
		return true
	}
	for _, participantUserID := range item.ParticipantUserIDs {
		if participantUserID == userID {
			return true
		}
	}
	return false
}

func normalizeMeetupState(status string) string {
	if strings.EqualFold(strings.TrimSpace(status), "cancelled") {
		return "cancelled"
	}
	return "open"
}

func compactStrings(values []string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			result = append(result, value)
		}
	}
	return result
}

func containsFold(values []string, target string) bool {
	target = strings.TrimSpace(target)
	for _, value := range values {
		if strings.EqualFold(strings.TrimSpace(value), target) {
			return true
		}
	}
	return false
}
