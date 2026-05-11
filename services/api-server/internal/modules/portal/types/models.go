package types

import "time"

type BannerQuery struct {
	Page     int
	PageSize int
	Keyword  string
}

type NoticeQuery struct {
	Page     int
	PageSize int
	Keyword  string
}

type BannerSaveRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	ActionURL   string `json:"action_url"`
	Sort        int    `json:"sort"`
}

type NoticePublishRequest struct {
	Title    string   `json:"title"`
	Summary  string   `json:"summary"`
	Content  string   `json:"content"`
	Audience string   `json:"audience"`
	Tags     []string `json:"tags"`
	Pinned   bool     `json:"pinned"`
}

type BannerItem struct {
	ID          string
	Title       string
	Description string
	ImageURL    string
	ActionURL   string
	Sort        int
	CreatedAt   time.Time
}

type NoticeItem struct {
	ID              string
	Title           string
	Summary         string
	Content         string
	Audience        string
	Tags            []string
	Pinned          bool
	PublisherUserID string
	Publisher       string
	PublishedAt     time.Time
	CreatedAt       time.Time
}
