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

type CarpoolQuery struct {
	Pagination
	Category string
	Keyword  string
}

type MeetupQuery struct {
	Pagination
	Category string
	Keyword  string
	UserRole string
}

type ReviewQuery struct {
	Pagination
	ContentType  string
	ReviewStatus string
	Keyword      string
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

type CarpoolPublishRequest struct {
	Category   string   `json:"category"`
	From       string   `json:"from"`
	To         string   `json:"to"`
	TravelDate string   `json:"travel_date"`
	TravelTime string   `json:"travel_time"`
	Time       string   `json:"time"`
	Type       string   `json:"type"`
	SeatsText  string   `json:"seats_text"`
	Price      string   `json:"price"`
	Note       string   `json:"note"`
	Tags       []string `json:"tags"`
	Contact    string   `json:"contact"`
}

type MeetupPublishRequest struct {
	Category        string   `json:"category"`
	Title           string   `json:"title"`
	Desc            string   `json:"desc"`
	Location        string   `json:"location"`
	StartAt         string   `json:"start_at"`
	DeadlineAt      string   `json:"deadline_at"`
	MaxParticipants int      `json:"max_participants"`
	FeeText         string   `json:"fee_text"`
	Tags            []string `json:"tags"`
	Contact         string   `json:"contact"`
}

type ReviewUpdateRequest struct {
	ContentType  string `json:"content_type"`
	ContentID    string `json:"content_id"`
	ReviewStatus string `json:"review_status"`
}

type MeetupActionRequest struct {
	MeetupID string `json:"meetup_id"`
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
	ReviewStatus     string
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
	ReviewStatus     string
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
	ReviewStatus     string
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
	ReviewStatus     string
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

type CarpoolItem struct {
	ID               string
	Category         string
	From             string
	To               string
	TravelAt         time.Time
	Type             string
	SeatsText        string
	Price            string
	Note             string
	Tags             []string
	Contact          string
	ReviewStatus     string
	PublisherUserID  string
	Publisher        string
	PublisherInitial string
	CreatedAt        time.Time
}

type MeetupItem struct {
	ID                 string
	Category           string
	Title              string
	Desc               string
	Location           string
	StartAt            time.Time
	DeadlineAt         time.Time
	MaxParticipants    int
	FeeText            string
	Tags               []string
	Contact            string
	Status             string
	ReviewStatus       string
	PublisherUserID    string
	Publisher          string
	PublisherInitial   string
	ParticipantUserIDs []string
	CreatedAt          time.Time
}
