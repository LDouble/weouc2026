package service

import (
	"errors"
	"slices"
	"strconv"
	"strings"
	"time"

	cltypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/campus_life/types"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/auth"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/httpx"
)

func buildVisibilityFilter(principal auth.Principal) (currentUserID string, visibleStatuses []string, includeAllStatus bool) {
	if canModerateCampusLife(principal) {
		return principal.UserID, nil, true
	}
	if !principal.Authenticated {
		return "", cltypes.VisibleStatuses, false
	}
	return principal.UserID, cltypes.VisibleStatuses, false
}

func simpleUserRole(publisherUserID string, principal auth.Principal) string {
	if !principal.Authenticated {
		return "viewer"
	}
	if publisherUserID == principal.UserID {
		return "publisher"
	}
	return "viewer"
}

func errandUserRole(item cltypes.CommunityContent, principal auth.Principal) string {
	if !principal.Authenticated {
		return "viewer"
	}
	if item.PublisherUserID == principal.UserID {
		return "publisher"
	}
	ep, _ := unmarshalPayload[cltypes.ErrandPayload](item.TypePayload)
	if ep.AcceptorUserID == principal.UserID {
		return "acceptor"
	}
	return "viewer"
}

func meetupUserRole(item cltypes.CommunityContent, principal auth.Principal) string {
	if !principal.Authenticated {
		return "viewer"
	}
	if item.PublisherUserID == principal.UserID {
		return "publisher"
	}
	mp, _ := unmarshalPayload[cltypes.MeetupPayload](item.TypePayload)
	if slices.Contains(mp.ParticipantUserIDs, principal.UserID) {
		return "participant"
	}
	return "viewer"
}

func canEditContent(isOwner bool, status string) bool {
	if !isOwner {
		return false
	}
	return status == cltypes.StatusPublished || status == cltypes.StatusReviewing
}

func canDeleteContent(isOwner bool, status string, extraConditions ...bool) bool {
	if !isOwner {
		return false
	}
	if status != cltypes.StatusPublished && status != cltypes.StatusReviewing && status != cltypes.StatusRejected {
		return false
	}
	for _, cond := range extraConditions {
		if !cond {
			return false
		}
	}
	return true
}

func matchFeedType(allowed []string, value string) bool {
	if len(allowed) == 0 {
		return true
	}
	return slices.Contains(allowed, value)
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

func listEnvelope(list []map[string]any, total int, pagination cltypes.Pagination) map[string]any {
	page := pagination.Page
	if page <= 0 {
		page = 1
	}
	pageSize := pagination.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	return map[string]any{
		"list":     list,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	}
}

func paginateRows[T any](items []T, pagination cltypes.Pagination) []T {
	page := pagination.Page
	if page <= 0 {
		page = 1
	}
	pageSize := pagination.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	start := (page - 1) * pageSize
	if start >= len(items) {
		return []T{}
	}
	end := min(len(items), start+pageSize)
	return items[start:end]
}

func paginateMaps(items []map[string]any, pagination cltypes.Pagination) []map[string]any {
	return paginateRows(items, pagination)
}

func canViewContact(principal auth.Principal, ownerUserID string) bool {
	if !principal.Authenticated {
		return false
	}
	return principal.AcademicBound || principal.UserID == ownerUserID
}

func visibleValue(allowed bool, value string) string {
	if !allowed {
		return ""
	}
	return value
}

func displayName(principal auth.Principal) string {
	if principal.DisplayName != "" {
		return principal.DisplayName
	}
	return "校园用户"
}

func initialOf(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return "校"
	}
	runes := []rune(text)
	return string(runes[0])
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func detectFileType(filePath string) string {
	lower := strings.ToLower(filePath)
	switch {
	case strings.HasSuffix(lower, ".pdf"):
		return "application/pdf"
	case strings.HasSuffix(lower, ".doc"), strings.HasSuffix(lower, ".docx"):
		return "application/msword"
	case strings.HasSuffix(lower, ".xls"), strings.HasSuffix(lower, ".xlsx"):
		return "application/vnd.ms-excel"
	case strings.HasSuffix(lower, ".ppt"), strings.HasSuffix(lower, ".pptx"):
		return "application/vnd.ms-powerpoint"
	case strings.HasSuffix(lower, ".jpg"), strings.HasSuffix(lower, ".jpeg"), strings.HasSuffix(lower, ".png"), strings.HasSuffix(lower, ".webp"):
		return "image/png"
	default:
		return "application/octet-stream"
	}
}

func firstFileType(files []cltypes.ResourceFile) string {
	if len(files) == 0 {
		return ""
	}
	return files[0].FileType
}

func firstFileSize(files []cltypes.ResourceFile) string {
	if len(files) == 0 {
		return ""
	}
	return files[0].FileSize
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func parsePage(value string, defaultValue int) int {
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || parsed <= 0 {
		return defaultValue
	}
	return parsed
}

func isAppError(err error) bool {
	var appErr *httpx.AppError
	return errors.As(err, &appErr)
}

func publisherName(item cltypes.CommunityContent, principal auth.Principal) string {
	if principal.Authenticated && item.PublisherUserID == principal.UserID {
		return displayName(principal)
	}
	return "校园用户"
}

func shouldExposeContent(principal auth.Principal, item cltypes.CommunityContent) bool {
	if item.Status == cltypes.StatusPublished || item.Status == cltypes.StatusOpen || item.Status == cltypes.StatusAccepted || item.Status == cltypes.StatusFull || item.Status == cltypes.StatusResolved {
		return true
	}
	return canAccessPendingContent(principal, item.PublisherUserID, "")
}

func ensureContentVisible(principal auth.Principal, item cltypes.CommunityContent, notFoundMessage string) error {
	if shouldExposeContent(principal, item) {
		return nil
	}
	return httpx.NotFound(notFoundMessage, nil)
}

func canAccessPendingContent(principal auth.Principal, ownerUserID, userRole string) bool {
	if canModerateCampusLife(principal) {
		return true
	}
	if !principal.Authenticated {
		return false
	}
	if userRole == "publisher" {
		return principal.UserID == ownerUserID
	}
	return principal.UserID == ownerUserID
}

func shouldExposeMeetupState(principal auth.Principal, item cltypes.CommunityContent, userRole string) bool {
	if item.Status != cltypes.StatusCancelled {
		return true
	}
	if canModerateCampusLife(principal) {
		return true
	}
	if !principal.Authenticated {
		return false
	}
	if userRole == "publisher" {
		return item.PublisherUserID == principal.UserID
	}
	if item.PublisherUserID == principal.UserID {
		return true
	}
	mp, _ := unmarshalPayload[cltypes.MeetupPayload](item.TypePayload)
	return slices.Contains(mp.ParticipantUserIDs, principal.UserID)
}

func meetupJoinedCountFromPayload(mp cltypes.MeetupPayload) int {
	return len(mp.ParticipantUserIDs) + 1
}

func meetupRemainingSeatsFromPayload(mp cltypes.MeetupPayload) int {
	remaining := mp.MaxParticipants - meetupJoinedCountFromPayload(mp)
	if remaining < 0 {
		return 0
	}
	return remaining
}

func refreshMeetupPayloadStatus(mp cltypes.MeetupPayload) cltypes.MeetupPayload {
	return mp
}

func canModerateCampusLife(principal auth.Principal) bool {
	return principal.HasPermission(campusLifeModeratePermission)
}

func isSupportedReviewStatus(status string) bool {
	switch status {
	case cltypes.StatusReviewing, cltypes.StatusPublished, cltypes.StatusRejected, cltypes.StatusOffline:
		return true
	default:
		return false
	}
}

func matchReviewQuery(query cltypes.ReviewQuery, contentType, status string, values ...string) bool {
	if query.ContentType != "" && query.ContentType != contentType {
		return false
	}
	if query.ReviewStatus != "" && query.ReviewStatus != status {
		return false
	}
	return matchKeyword(query.Keyword, values...)
}

func parseCarpoolTravelAt(dateText, timeText string) (time.Time, error) {
	return time.ParseInLocation(
		"2006-01-02 15:04",
		strings.TrimSpace(dateText)+" "+strings.TrimSpace(timeText),
		chinaLocation,
	)
}

func normalizedCarpoolCategoryFromPayload(cp cltypes.CarpoolPayload, now time.Time) string {
	if cp.TravelAt.IsZero() {
		if isSupportedCarpoolCategory(cp.Category) {
			return cp.Category
		}
		return "today"
	}

	travelDate := startOfDay(cp.TravelAt.In(chinaLocation))
	today := startOfDay(now.In(chinaLocation))
	tomorrow := today.AddDate(0, 0, 1)
	if travelDate.Equal(today) {
		return "today"
	}
	if travelDate.Equal(tomorrow) {
		return "tomorrow"
	}
	if !travelDate.After(endOfWeek(today)) {
		return "week"
	}
	return "longterm"
}

func formatCarpoolTravelText(travelAt time.Time, now time.Time) string {
	if travelAt.IsZero() {
		return ""
	}
	travelLocal := travelAt.In(chinaLocation)
	travelDate := startOfDay(travelLocal)
	today := startOfDay(now.In(chinaLocation))
	tomorrow := today.AddDate(0, 0, 1)

	switch {
	case travelDate.Equal(today):
		return "今天 " + travelLocal.Format("15:04")
	case travelDate.Equal(tomorrow):
		return "明天 " + travelLocal.Format("15:04")
	default:
		return travelLocal.Format("1月2日 15:04")
	}
}

func defaultCarpoolType(category string) string {
	switch category {
	case "tomorrow":
		return "明日顺路"
	case "week":
		return "本周拼车"
	case "longterm":
		return "长期通勤"
	default:
		return "今日顺路"
	}
}

func isSupportedCarpoolCategory(category string) bool {
	switch category {
	case "today", "tomorrow", "week", "longterm":
		return true
	default:
		return false
	}
}

func sanitizeTags(values []string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		result = append(result, value)
	}
	return result
}

func startOfDay(value time.Time) time.Time {
	local := value.In(chinaLocation)
	return time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, chinaLocation)
}

func endOfWeek(value time.Time) time.Time {
	local := startOfDay(value)
	daysUntilSunday := (7 - int(local.Weekday())) % 7
	return local.AddDate(0, 0, daysUntilSunday)
}

func carpoolTitleFromPayload(cp cltypes.CarpoolPayload, now time.Time) string {
	if cp.From != "" || cp.To != "" {
		return strings.TrimSpace(cp.From + " -> " + cp.To)
	}
	return defaultCarpoolType(normalizedCarpoolCategoryFromPayload(cp, now))
}

func carpoolFeedDescFromPayload(cp cltypes.CarpoolPayload, now time.Time) string {
	parts := make([]string, 0, 3)
	if text := formatCarpoolTravelText(cp.TravelAt, now); text != "" {
		parts = append(parts, text)
	}
	if cp.SeatsText != "" {
		parts = append(parts, cp.SeatsText)
	}
	if cp.Price != "" {
		parts = append(parts, cp.Price)
	}
	if len(parts) > 0 {
		return strings.Join(parts, " · ")
	}
	return cp.Note
}

func meetupFeedDescFromPayload(mp cltypes.MeetupPayload) string {
	parts := make([]string, 0, 4)
	if mp.Location != "" {
		parts = append(parts, mp.Location)
	}
	if !mp.StartAt.IsZero() {
		parts = append(parts, mp.StartAt.In(chinaLocation).Format("1月2日 15:04"))
	}
	if remaining := meetupRemainingSeatsFromPayload(mp); remaining > 0 {
		parts = append(parts, "剩余 "+strconv.Itoa(remaining)+" 位")
	} else {
		parts = append(parts, "人数已满")
	}
	if mp.FeeText != "" {
		parts = append(parts, mp.FeeText)
	}
	return strings.Join(parts, " · ")
}

func buildCarpoolPayload(item cltypes.CommunityContent, cp cltypes.CarpoolPayload, canView bool, now time.Time) map[string]any {
	category := normalizedCarpoolCategoryFromPayload(cp, now)
	timeText := formatCarpoolTravelText(cp.TravelAt, now)
	return map[string]any{
		"id":                item.ID.Hex(),
		"category":          category,
		"from":              cp.From,
		"to":                cp.To,
		"time":              timeText,
		"type":              cp.Type,
		"seats_text":        cp.SeatsText,
		"price":             cp.Price,
		"note":              cp.Note,
		"tags":              append([]string(nil), item.Tags...),
		"contact":           visibleValue(canView, item.Contact),
		"status":            item.Status,
		"publisher":         "校园用户",
		"publisher_initial": "校",
		"created_at":        item.CreatedAt.Format(time.RFC3339),
		"extra": map[string]any{
			"category":   category,
			"from":       cp.From,
			"to":         cp.To,
			"time":       timeText,
			"type":       cp.Type,
			"seats_text": cp.SeatsText,
			"price":      cp.Price,
			"note":       cp.Note,
			"tags":       append([]string(nil), item.Tags...),
			"contact":    visibleValue(canView, item.Contact),
			"travel_at":  cp.TravelAt.In(chinaLocation).Format(time.RFC3339),
			"status":     item.Status,
		},
	}
}

func buildMeetupPayload(item cltypes.CommunityContent, mp cltypes.MeetupPayload, principal auth.Principal, now time.Time) map[string]any {
	canView := canViewContact(principal, item.PublisherUserID)
	userRole := meetupUserRole(item, principal)
	joinedCount := meetupJoinedCountFromPayload(mp)
	remainingSeats := meetupRemainingSeatsFromPayload(mp)
	canJoin := principal.Authenticated &&
		userRole == "viewer" &&
		item.Status == cltypes.StatusOpen &&
		remainingSeats > 0 &&
		(mp.DeadlineAt.IsZero() || mp.DeadlineAt.After(now.UTC())) &&
		(mp.StartAt.IsZero() || mp.StartAt.After(now.UTC()))

	return map[string]any{
		"id":                 item.ID.Hex(),
		"category":           mp.Category,
		"title":              item.Title,
		"desc":               item.Desc,
		"location":           mp.Location,
		"start_at":           mp.StartAt.In(chinaLocation).Format(time.RFC3339),
		"deadline_at":        mp.DeadlineAt.In(chinaLocation).Format(time.RFC3339),
		"max_participants":   mp.MaxParticipants,
		"joined_count":       joinedCount,
		"remaining_seats":    remainingSeats,
		"fee_text":           mp.FeeText,
		"tags":               append([]string(nil), item.Tags...),
		"contact":            visibleValue(canView, item.Contact),
		"status":             item.Status,
		"publisher":          publisherName(item, principal),
		"publisher_initial":  initialOf(publisherName(item, principal)),
		"created_at":         item.CreatedAt.Format(time.RFC3339),
		"user_role":          userRole,
		"joined":             userRole == "participant",
		"can_join":           canJoin,
		"can_cancel_join":    userRole == "participant" && item.Status != cltypes.StatusCancelled,
		"can_cancel_publish": userRole == "publisher" && item.Status != cltypes.StatusCancelled,
		"extra": map[string]any{
			"category":         mp.Category,
			"location":         mp.Location,
			"start_at":         mp.StartAt.In(chinaLocation).Format(time.RFC3339),
			"deadline_at":      mp.DeadlineAt.In(chinaLocation).Format(time.RFC3339),
			"max_participants": mp.MaxParticipants,
			"joined_count":     joinedCount,
			"remaining_seats":  remainingSeats,
			"fee_text":         mp.FeeText,
			"tags":             append([]string(nil), item.Tags...),
			"contact":          visibleValue(canView, item.Contact),
			"status":           item.Status,
		},
	}
}

func matchUserRole(publisherUserID, acceptorUserID string, principal auth.Principal, userRole string) bool {
	if userRole == "" {
		return true
	}
	if !principal.Authenticated {
		return false
	}
	switch userRole {
	case "publisher":
		return publisherUserID == principal.UserID
	case "acceptor":
		return acceptorUserID == principal.UserID
	default:
		return true
	}
}

func matchMeetupUserRole(item cltypes.CommunityContent, principal auth.Principal, userRole string) bool {
	if userRole == "" {
		return true
	}
	mp, _ := unmarshalPayload[cltypes.MeetupPayload](item.TypePayload)
	switch userRole {
	case "publisher":
		return principal.Authenticated && item.PublisherUserID == principal.UserID
	case "participant":
		return principal.Authenticated && slices.Contains(mp.ParticipantUserIDs, principal.UserID)
	default:
		return true
	}
}
