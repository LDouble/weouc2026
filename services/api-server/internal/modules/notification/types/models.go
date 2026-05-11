package types

import "time"

type MessageQuery struct {
	Page       int
	PageSize   int
	Category   string
	UnreadOnly bool
}

type PublishRequest struct {
	Title         string   `json:"title"`
	Content       string   `json:"content"`
	Category      string   `json:"category"`
	TargetScope   string   `json:"target_scope"`
	TargetUserIDs []string `json:"target_user_ids"`
	ActionURL     string   `json:"action_url"`
}

type MarkReadRequest struct {
	MessageID string `json:"message_id"`
}

type MessageItem struct {
	ID              string
	Title           string
	Content         string
	Category        string
	TargetScope     string
	TargetUserIDs   []string
	ActionURL       string
	PublisherUserID string
	Publisher       string
	CreatedAt       time.Time
	ReadByUserIDs   map[string]time.Time
}
