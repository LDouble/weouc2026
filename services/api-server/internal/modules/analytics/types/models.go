package types

type AuditLogQuery struct {
	Page         int
	PageSize     int
	ActorID      string
	Action       string
	ResourceType string
	ResourceID   string
}
