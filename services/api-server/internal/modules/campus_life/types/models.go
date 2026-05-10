package types

import "time"

type Pagination struct {
	Page     int
	PageSize int
}

type FeedQuery struct {
	Pagination
	FeedTypes []string
	Keyword   string
	UserRole  string
}

type MarketQuery struct {
	Pagination
	Category string
	Keyword  string
}

type ErrandQuery struct {
	Pagination
	Category string
	Keyword  string
	UserRole string
}

type ResourceQuery struct {
	Pagination
	Category string
	Keyword  string
}

type LostFoundQuery struct {
	Pagination
	Category string
	Keyword  string
	Type     string
}

type MarketPublishRequest struct {
	Title         string   `json:"title"`
	Desc          string   `json:"desc"`
	Price         string   `json:"price"`
	OriginalPrice string   `json:"original_price"`
	Category      string   `json:"category"`
	Condition     string   `json:"condition"`
	TradeMode     string   `json:"trade_mode"`
	Contact       string   `json:"contact"`
	Images        []string `json:"images"`
}

type FavoriteMarketRequest struct {
	ProductID string `json:"product_id"`
	Action    string `json:"action"`
}

type ErrandPublishRequest struct {
	Title      string   `json:"title"`
	Desc       string   `json:"desc"`
	Category   string   `json:"category"`
	RouteStart string   `json:"route_start"`
	RouteEnd   string   `json:"route_end"`
	Deadline   string   `json:"deadline"`
	Reward     string   `json:"reward"`
	Contact    string   `json:"contact"`
	Urgent     bool     `json:"urgent"`
	Images     []string `json:"images"`
}

type ErrandActionRequest struct {
	TaskID string `json:"task_id"`
}

type ResourcePublishRequest struct {
	Title      string   `json:"title"`
	Desc       string   `json:"desc"`
	Category   string   `json:"category"`
	CourseName string   `json:"course_name"`
	Contact    string   `json:"contact"`
	FilePaths  []string `json:"file_paths"`
}

type LostFoundPublishRequest struct {
	Type        string `json:"type"`
	Category    string `json:"category"`
	Title       string `json:"title"`
	Desc        string `json:"desc"`
	Location    string `json:"location"`
	EventTime   string `json:"event_time"`
	ItemFeature string `json:"item_feature"`
	Contact     string `json:"contact"`
	Reward      string `json:"reward"`
}

type ResourceFile struct {
	Name     string `json:"name"`
	Path     string `json:"path,omitempty"`
	URL      string `json:"url"`
	FileType string `json:"file_type"`
	FileSize string `json:"file_size"`
}

type MarketItem struct {
	ID               string
	Title            string
	Desc             string
	PublisherUserID  string
	Publisher        string
	PublisherInitial string
	Image            string
	CreatedAt        time.Time
	Likes            int
	LikedByUserIDs   map[string]bool
	Extra            MarketExtra
}

type MarketExtra struct {
	Category      string
	Price         string
	OriginalPrice string
	Condition     string
	TradeMode     string
	Contact       string
	Images        []string
}

type ErrandItem struct {
	ID               string
	Title            string
	Desc             string
	Category         string
	RouteStart       string
	RouteEnd         string
	Deadline         time.Time
	Reward           string
	Contact          string
	Urgent           bool
	Images           []string
	Status           string
	PublisherUserID  string
	Publisher        string
	PublisherInitial string
	AcceptorUserID   string
	CreatedAt        time.Time
}

type ResourceItem struct {
	ID               string
	Title            string
	Desc             string
	PublisherUserID  string
	Publisher        string
	PublisherInitial string
	CreatedAt        time.Time
	Extra            ResourceExtra
}

type ResourceExtra struct {
	Category    string
	CourseName  string
	Contact     string
	Files       []ResourceFile
	FileSize    string
	FileType    string
	DownloadURL string
	Likes       int
	Views       int
}

type LostFoundItem struct {
	ID               string
	Title            string
	Desc             string
	PublisherUserID  string
	Publisher        string
	PublisherInitial string
	CreatedAt        time.Time
	Extra            LostFoundExtra
}

type LostFoundExtra struct {
	Type        string
	Category    string
	Location    string
	EventTime   string
	ItemFeature string
	Contact     string
}
