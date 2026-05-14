package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ContentTypeMarket    = "market"
	ContentTypeErrand    = "errand"
	ContentTypeResource  = "resource"
	ContentTypeLostFound = "lost_found"
	ContentTypeCarpool   = "carpool"
	ContentTypeMeetup    = "meetup"
)

const (
	StatusReviewing = "reviewing"
	StatusPublished = "published"
	StatusRejected  = "rejected"
	StatusOffline   = "offline"
	StatusCancelled = "cancelled"
	StatusAccepted  = "accepted"
	StatusOpen      = "open"
	StatusFull      = "full"
	StatusResolved  = "resolved"
)

type CommunityContent struct {
	ID              primitive.ObjectID `bson:"_id"`
	ContentType     string             `bson:"content_type"`
	Title           string             `bson:"title"`
	Desc            string             `bson:"desc"`
	Status          string             `bson:"status"`
	PublisherUserID string             `bson:"publisher_user_id"`
	Contact         string             `bson:"contact"`
	Images          []string           `bson:"images"`
	Tags            []string           `bson:"tags"`
	CreatedAt       time.Time          `bson:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at"`
	CreatedBy       string             `bson:"created_by"`
	UpdatedBy       string             `bson:"updated_by"`
	DeletedAt       *time.Time         `bson:"deleted_at"`
	TypePayload     bson.M             `bson:"type_payload"`
	ExtJSON         map[string]any     `bson:"ext_json,omitempty"`
	LikedByUserIDs  map[string]bool    `bson:"liked_by_user_ids,omitempty"`
	Likes           int                `bson:"likes,omitempty"`
}

type MarketPayload struct {
	Category      string `json:"category" bson:"category"`
	Price         string `json:"price" bson:"price"`
	OriginalPrice string `json:"original_price" bson:"original_price"`
	Condition     string `json:"condition" bson:"condition"`
	TradeMode     string `json:"trade_mode" bson:"trade_mode"`
}

type ErrandPayload struct {
	Category       string    `json:"category" bson:"category"`
	RouteStart     string    `json:"route_start" bson:"route_start"`
	RouteEnd       string    `json:"route_end" bson:"route_end"`
	Deadline       time.Time `json:"deadline" bson:"deadline"`
	Reward         string    `json:"reward" bson:"reward"`
	Urgent         bool      `json:"urgent" bson:"urgent"`
	AcceptorUserID string    `json:"acceptor_user_id" bson:"acceptor_user_id"`
}

type ResourcePayload struct {
	Category    string         `json:"category" bson:"category"`
	CourseName  string         `json:"course_name" bson:"course_name"`
	Files       []ResourceFile `json:"files" bson:"files"`
	FileSize    string         `json:"file_size" bson:"file_size"`
	FileType    string         `json:"file_type" bson:"file_type"`
	DownloadURL string         `json:"download_url" bson:"download_url"`
	Likes       int            `json:"likes" bson:"likes"`
	Views       int            `json:"views" bson:"views"`
}

type LostFoundPayload struct {
	Type        string `json:"type" bson:"type"`
	Category    string `json:"category" bson:"category"`
	Location    string `json:"location" bson:"location"`
	EventTime   string `json:"event_time" bson:"event_time"`
	ItemFeature string `json:"item_feature" bson:"item_feature"`
}

type CarpoolPayload struct {
	Category  string    `json:"category" bson:"category"`
	From      string    `json:"from" bson:"from"`
	To        string    `json:"to" bson:"to"`
	TravelAt  time.Time `json:"travel_at" bson:"travel_at"`
	Type      string    `json:"type" bson:"type"`
	SeatsText string    `json:"seats_text" bson:"seats_text"`
	Price     string    `json:"price" bson:"price"`
	Note      string    `json:"note" bson:"note"`
}

type MeetupPayload struct {
	Category           string    `json:"category" bson:"category"`
	Location           string    `json:"location" bson:"location"`
	StartAt            time.Time `json:"start_at" bson:"start_at"`
	DeadlineAt         time.Time `json:"deadline_at" bson:"deadline_at"`
	MaxParticipants    int       `json:"max_participants" bson:"max_participants"`
	FeeText            string    `json:"fee_text" bson:"fee_text"`
	ParticipantUserIDs []string  `json:"participant_user_ids" bson:"participant_user_ids"`
}

type ResourceFile struct {
	Name     string `json:"name"`
	Path     string `json:"path,omitempty"`
	URL      string `json:"url"`
	FileType string `json:"file_type"`
	FileSize string `json:"file_size"`
}

type StateTransitionLog struct {
	ID          primitive.ObjectID `bson:"_id"`
	ContentType string             `bson:"content_type"`
	ContentID   string             `bson:"content_id"`
	FromStatus  string             `bson:"from_status"`
	ToStatus    string             `bson:"to_status"`
	Action      string             `bson:"action"`
	ActorUserID string             `bson:"actor_user_id"`
	CreatedAt   time.Time          `bson:"created_at"`
}

type Pagination struct {
	Page     int
	PageSize int
}

var VisibleStatuses = []string{
	StatusPublished,
	StatusOpen,
	StatusAccepted,
	StatusFull,
	StatusResolved,
}

type ContentFilter struct {
	Pagination
	ContentType       string
	Statuses          []string
	Keyword           string
	Category          string
	SubType           string
	PublisherUserID   string
	AcceptorUserID    string
	ParticipantUserID string
	VisibleStatuses   []string
	IncludeAllStatus  bool
	CurrentUserID     string
}

type FeedFilter struct {
	Pagination
	FeedTypes         []string
	Keyword           string
	PublisherUserID   string
	AcceptorUserID    string
	ParticipantUserID string
	VisibleStatuses   []string
	IncludeAllStatus  bool
	CurrentUserID     string
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
	ContentType string `json:"content_type"`
	ContentID   string `json:"content_id"`
	Action      string `json:"action"`
}

type MeetupActionRequest struct {
	MeetupID string `json:"meetup_id"`
}

type FeedQuery struct {
	Pagination
	FeedTypes []string
	Keyword   string
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
}

type ReviewQuery struct {
	Pagination
	ContentType  string
	ReviewStatus string
	Keyword      string
}
