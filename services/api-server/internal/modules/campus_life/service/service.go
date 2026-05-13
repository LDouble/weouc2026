package service

import (
	"context"
	"errors"
	"path"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	clrepo "github.com/liangluo/weouc2026/services/api-server/internal/modules/campus_life/repo"
	cltypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/campus_life/types"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/audit"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/auth"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/httpx"
	"github.com/liangluo/weouc2026/services/api-server/internal/providers/storage_provider"
)

type Service struct {
	repository      clrepo.Repository
	storageProvider storage_provider.Provider
	recorder        audit.Recorder
}

var chinaLocation = time.FixedZone("Asia/Shanghai", 8*3600)

const (
	statusReviewing = "reviewing"
	statusPublished = "published"
	statusRejected  = "rejected"
	statusOffline   = "offline"
	statusCancelled = "cancelled"
	statusFull      = "full"
	statusAccepted  = "accepted"
	statusOpen      = "open"
	statusResolved  = "resolved"

	campusLifeModeratePermission = "campus_life:moderate"
)

func New(repository clrepo.Repository, storageProvider storage_provider.Provider, recorder audit.Recorder) *Service {
	return &Service{
		repository:      repository,
		storageProvider: storageProvider,
		recorder:        recorder,
	}
}

func (s *Service) ListFeed(ctx context.Context, principal auth.Principal, query cltypes.FeedQuery) (map[string]any, error) {
	type row struct {
		id        string
		feedType  string
		createdAt time.Time
		payload   map[string]any
	}

	rows := make([]row, 0)
	markets, err := s.repository.ListMarkets(ctx)
	if err != nil {
		return nil, httpx.Internal("读取二手动态失败", err)
	}
	for _, item := range markets {
		if !matchUserRole(item.PublisherUserID, "", principal, query.UserRole) {
			continue
		}
		if !shouldExposeContent(principal, item.PublisherUserID, item.ReviewStatus, query.UserRole) {
			continue
		}
		if !matchFeedType(query.FeedTypes, "market") || !matchKeyword(query.Keyword, item.Title, item.Desc) {
			continue
		}
		resolvedImages := resolveManagedURLs(ctx, s.storageProvider, item.Extra.Images)
		resolvedImage := resolveManagedURL(ctx, s.storageProvider, item.Image)
		if resolvedImage == "" {
			resolvedImage = firstNonEmpty(resolvedImages...)
		}
		rows = append(rows, row{
			id:        item.ID,
			feedType:  "market",
			createdAt: item.CreatedAt,
			payload: map[string]any{
				"id":              item.ID,
				"feed_type":       "market",
				"feed_type_label": "二手交易",
				"title":           item.Title,
				"desc":            item.Desc,
				"publisher":       item.Publisher,
				"created_at":      item.CreatedAt.Format(time.RFC3339),
				"status":          normalizeReviewStatus(item.ReviewStatus),
				"image":           resolvedImage,
				"extra": map[string]any{
					"images":   resolvedImages,
					"likes":    item.Likes,
					"comments": 0,
				},
			},
		})
	}
	errands, err := s.repository.ListErrands(ctx)
	if err != nil {
		return nil, httpx.Internal("读取跑腿动态失败", err)
	}
	for _, item := range errands {
		role := errandUserRole(item, principal)
		if !matchUserRole(item.PublisherUserID, item.AcceptorUserID, principal, query.UserRole) {
			continue
		}
		if !shouldExposeContent(principal, item.PublisherUserID, item.ReviewStatus, query.UserRole) {
			continue
		}
		if !matchFeedType(query.FeedTypes, "errand") || !matchKeyword(query.Keyword, item.Title, item.Desc) {
			continue
		}
		rows = append(rows, row{
			id:        item.ID,
			feedType:  "errand",
			createdAt: item.CreatedAt,
			payload: map[string]any{
				"id":              item.ID,
				"feed_type":       "errand",
				"feed_type_label": "校园跑腿",
				"title":           item.Title,
				"desc":            item.Desc,
				"publisher":       item.Publisher,
				"created_at":      item.CreatedAt.Format(time.RFC3339),
				"status":          mergeErrandStatus(item.Status, item.ReviewStatus),
				"extra": map[string]any{
					"likes":     0,
					"comments":  0,
					"user_role": role,
				},
			},
		})
	}
	resources, err := s.repository.ListResources(ctx)
	if err != nil {
		return nil, httpx.Internal("读取资料动态失败", err)
	}
	for _, item := range resources {
		if !matchUserRole(item.PublisherUserID, "", principal, query.UserRole) {
			continue
		}
		if !shouldExposeContent(principal, item.PublisherUserID, item.ReviewStatus, query.UserRole) {
			continue
		}
		if !matchFeedType(query.FeedTypes, "resource") || !matchKeyword(query.Keyword, item.Title, item.Desc) {
			continue
		}
		rows = append(rows, row{
			id:        item.ID,
			feedType:  "resource",
			createdAt: item.CreatedAt,
			payload: map[string]any{
				"id":              item.ID,
				"feed_type":       "resource",
				"feed_type_label": "资料共享",
				"title":           item.Title,
				"desc":            item.Desc,
				"publisher":       item.Publisher,
				"created_at":      item.CreatedAt.Format(time.RFC3339),
				"status":          normalizeReviewStatus(item.ReviewStatus),
			},
		})
	}
	lostFounds, err := s.repository.ListLostFound(ctx)
	if err != nil {
		return nil, httpx.Internal("读取失物招领动态失败", err)
	}
	for _, item := range lostFounds {
		if !matchUserRole(item.PublisherUserID, "", principal, query.UserRole) {
			continue
		}
		if !shouldExposeContent(principal, item.PublisherUserID, item.ReviewStatus, query.UserRole) {
			continue
		}
		if !matchFeedType(query.FeedTypes, "lostFound") || !matchKeyword(query.Keyword, item.Title, item.Desc) {
			continue
		}
		rows = append(rows, row{
			id:        item.ID,
			feedType:  "lostFound",
			createdAt: item.CreatedAt,
			payload: map[string]any{
				"id":              item.ID,
				"feed_type":       "lostFound",
				"feed_type_label": "失物招领",
				"title":           item.Title,
				"desc":            item.Desc,
				"publisher":       item.Publisher,
				"created_at":      item.CreatedAt.Format(time.RFC3339),
				"status":          normalizeReviewStatus(item.ReviewStatus),
			},
		})
	}
	now := time.Now().In(chinaLocation)
	carpools, err := s.repository.ListCarpools(ctx)
	if err != nil {
		return nil, httpx.Internal("读取拼车动态失败", err)
	}
	for _, item := range carpools {
		if !matchUserRole(item.PublisherUserID, "", principal, query.UserRole) {
			continue
		}
		if !shouldExposeContent(principal, item.PublisherUserID, item.ReviewStatus, query.UserRole) {
			continue
		}
		if !matchFeedType(query.FeedTypes, "carpool") ||
			!matchKeyword(query.Keyword, item.From, item.To, item.Note, item.Type) {
			continue
		}
		rows = append(rows, row{
			id:        item.ID,
			feedType:  "carpool",
			createdAt: item.CreatedAt,
			payload: map[string]any{
				"id":              item.ID,
				"feed_type":       "carpool",
				"feed_type_label": "校园拼车",
				"title":           carpoolTitle(item),
				"desc":            carpoolFeedDesc(item),
				"publisher":       item.Publisher,
				"created_at":      item.CreatedAt.Format(time.RFC3339),
				"status":          normalizeReviewStatus(item.ReviewStatus),
				"extra": map[string]any{
					"category":   normalizedCarpoolCategory(item, now),
					"from":       item.From,
					"to":         item.To,
					"time":       formatCarpoolTravelText(item.TravelAt, now),
					"type":       item.Type,
					"seats_text": item.SeatsText,
					"price":      item.Price,
					"tags":       append([]string(nil), item.Tags...),
					"comments":   0,
				},
			},
		})
	}
	meetups, err := s.repository.ListMeetups(ctx)
	if err != nil {
		return nil, httpx.Internal("读取组局动态失败", err)
	}
	for _, item := range meetups {
		if !matchMeetupUserRole(item, principal, query.UserRole) {
			continue
		}
		if !shouldExposeContent(principal, item.PublisherUserID, item.ReviewStatus, query.UserRole) {
			continue
		}
		if !shouldExposeMeetupState(principal, item, query.UserRole) {
			continue
		}
		if !matchFeedType(query.FeedTypes, "meetup") || !matchKeyword(query.Keyword, item.Title, item.Desc, item.Location) {
			continue
		}
		rows = append(rows, row{
			id:        item.ID,
			feedType:  "meetup",
			createdAt: item.CreatedAt,
			payload: map[string]any{
				"id":              item.ID,
				"feed_type":       "meetup",
				"feed_type_label": "校园组局",
				"title":           item.Title,
				"desc":            meetupFeedDesc(item),
				"publisher":       item.Publisher,
				"created_at":      item.CreatedAt.Format(time.RFC3339),
				"status":          mergeMeetupStatus(item.Status, item.ReviewStatus),
				"extra": map[string]any{
					"category":        item.Category,
					"location":        item.Location,
					"start_at":        item.StartAt.In(chinaLocation).Format(time.RFC3339),
					"joined_count":    meetupJoinedCount(item),
					"remaining_seats": meetupRemainingSeats(item),
					"fee_text":        item.FeeText,
					"tags":            append([]string(nil), item.Tags...),
					"comments":        0,
				},
			},
		})
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].createdAt.After(rows[j].createdAt)
	})

	items := paginateRows(rows, query.Pagination)
	list := make([]map[string]any, 0, len(items))
	for _, item := range items {
		list = append(list, item.payload)
	}

	return listEnvelope(list, len(rows), query.Pagination), nil
}

func (s *Service) ListMarket(ctx context.Context, principal auth.Principal, query cltypes.MarketQuery) (map[string]any, error) {
	items, err := s.repository.ListMarkets(ctx)
	if err != nil {
		return nil, httpx.Internal("读取二手列表失败", err)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt.After(items[j].CreatedAt)
	})

	filtered := make([]map[string]any, 0)
	for _, item := range items {
		if !shouldExposeContent(principal, item.PublisherUserID, item.ReviewStatus, "") {
			continue
		}
		if query.Category != "" && query.Category != item.Extra.Category {
			continue
		}
		if !matchKeyword(query.Keyword, item.Title, item.Desc) {
			continue
		}
		canViewContact := canViewContact(principal, item.PublisherUserID)
		isOwner := item.PublisherUserID == principal.UserID
		resolvedImages := resolveManagedURLs(ctx, s.storageProvider, item.Extra.Images)
		resolvedImage := resolveManagedURL(ctx, s.storageProvider, item.Image)
		if resolvedImage == "" {
			resolvedImage = firstNonEmpty(resolvedImages...)
		}
		filtered = append(filtered, map[string]any{
			"id":                item.ID,
			"title":             item.Title,
			"desc":              item.Desc,
			"publisher":         item.Publisher,
			"publisher_initial": item.PublisherInitial,
			"image":             resolvedImage,
			"likes":             item.Likes,
			"liked":             item.LikedByUserIDs[principal.UserID],
			"created_at":        item.CreatedAt.Format(time.RFC3339),
			"status":            normalizeReviewStatus(item.ReviewStatus),
			"is_owner":          isOwner,
			"can_edit":          canEditContent(isOwner, item.ReviewStatus),
			"can_delete":        canDeleteContent(isOwner, item.ReviewStatus),
			"can_favorite":      !isOwner && principal.Authenticated,
			"extra": map[string]any{
				"category":       item.Extra.Category,
				"price":          item.Extra.Price,
				"original_price": item.Extra.OriginalPrice,
				"condition":      item.Extra.Condition,
				"trade_mode":     item.Extra.TradeMode,
				"contact":        visibleValue(canViewContact, item.Extra.Contact),
				"images":         resolvedImages,
				"likes":          item.Likes,
				"is_favorited":   item.LikedByUserIDs[principal.UserID],
			},
		})
	}

	return listEnvelope(paginateMaps(filtered, query.Pagination), len(filtered), query.Pagination), nil
}

func (s *Service) GetMarketDetail(ctx context.Context, principal auth.Principal, id string) (map[string]any, error) {
	item, err := s.repository.GetMarket(ctx, id)
	if errors.Is(err, clrepo.ErrNotFound) {
		return nil, httpx.NotFound("二手商品不存在", nil)
	}
	if err != nil {
		return nil, httpx.Internal("读取二手详情失败", err)
	}
	if err := ensureContentVisible(principal, item.PublisherUserID, item.ReviewStatus, "二手商品不存在"); err != nil {
		return nil, err
	}

	canView := canViewContact(principal, item.PublisherUserID)
	role := simpleUserRole(item.PublisherUserID, principal)
	isOwner := role == "publisher"
	resolvedImages := resolveManagedURLs(ctx, s.storageProvider, item.Extra.Images)
	resolvedImage := resolveManagedURL(ctx, s.storageProvider, item.Image)
	if resolvedImage == "" {
		resolvedImage = firstNonEmpty(resolvedImages...)
	}
	return map[string]any{
		"id":                item.ID,
		"title":             item.Title,
		"desc":              item.Desc,
		"publisher":         item.Publisher,
		"publisher_initial": item.PublisherInitial,
		"image":             resolvedImage,
		"likes":             item.Likes,
		"liked":             item.LikedByUserIDs[principal.UserID],
		"created_at":        item.CreatedAt.Format(time.RFC3339),
		"status":            normalizeReviewStatus(item.ReviewStatus),
		"user_role":         role,
		"is_owner":          isOwner,
		"can_view_contact":  canView,
		"can_edit":          canEditContent(isOwner, item.ReviewStatus),
		"can_delete":        canDeleteContent(isOwner, item.ReviewStatus),
		"can_favorite":      !isOwner && principal.Authenticated,
		"extra": map[string]any{
			"category":       item.Extra.Category,
			"price":          item.Extra.Price,
			"original_price": item.Extra.OriginalPrice,
			"condition":      item.Extra.Condition,
			"trade_mode":     item.Extra.TradeMode,
			"contact":        visibleValue(canView, item.Extra.Contact),
			"images":         resolvedImages,
			"likes":          item.Likes,
			"is_favorited":   item.LikedByUserIDs[principal.UserID],
		},
	}, nil
}

func (s *Service) PublishMarket(ctx context.Context, principal auth.Principal, request cltypes.MarketPublishRequest) (map[string]any, error) {
	if strings.TrimSpace(request.Title) == "" || strings.TrimSpace(request.Desc) == "" {
		return nil, httpx.BadRequest("标题和描述不能为空", nil)
	}
	id, err := s.repository.NextID(ctx, "market")
	if err != nil {
		return nil, httpx.Internal("生成二手信息 ID 失败", err)
	}
	item := cltypes.MarketItem{
		ID:               id,
		Title:            strings.TrimSpace(request.Title),
		Desc:             strings.TrimSpace(request.Desc),
		ReviewStatus:     statusReviewing,
		PublisherUserID:  principal.UserID,
		Publisher:        displayName(principal),
		PublisherInitial: initialOf(displayName(principal)),
		Image:            firstNonEmpty(request.Images...),
		CreatedAt:        time.Now().UTC(),
		LikedByUserIDs:   map[string]bool{},
		Extra: cltypes.MarketExtra{
			Category:      strings.TrimSpace(request.Category),
			Price:         strings.TrimSpace(request.Price),
			OriginalPrice: strings.TrimSpace(request.OriginalPrice),
			Condition:     strings.TrimSpace(request.Condition),
			TradeMode:     strings.TrimSpace(request.TradeMode),
			Contact:       strings.TrimSpace(request.Contact),
			Images:        append([]string(nil), request.Images...),
		},
	}
	item, err = s.repository.SaveMarket(ctx, item)
	if err != nil {
		return nil, httpx.Internal("保存二手信息失败", err)
	}

	s.recordAudit(ctx, principal, "campus_life.market.publish", "market", item.ID, "二手信息发布成功", map[string]any{
		"review_status": item.ReviewStatus,
		"category":      item.Extra.Category,
	})

	return map[string]any{"id": item.ID}, nil
}

func (s *Service) FavoriteMarket(ctx context.Context, principal auth.Principal, request cltypes.FavoriteMarketRequest) error {
	if request.ProductID == "" || principal.UserID == "" {
		return httpx.BadRequest("缺少商品或用户信息", nil)
	}

	_, err := s.repository.UpdateMarket(ctx, request.ProductID, func(item *cltypes.MarketItem) error {
		if item.LikedByUserIDs == nil {
			item.LikedByUserIDs = map[string]bool{}
		}
		switch request.Action {
		case "add":
			if !item.LikedByUserIDs[principal.UserID] {
				item.LikedByUserIDs[principal.UserID] = true
				item.Likes++
			}
		case "remove":
			if item.LikedByUserIDs[principal.UserID] {
				delete(item.LikedByUserIDs, principal.UserID)
				if item.Likes > 0 {
					item.Likes--
				}
			}
		default:
			return httpx.BadRequest("action 仅支持 add/remove", nil)
		}
		return nil
	})
	if err != nil {
		if err == clrepo.ErrNotFound {
			return httpx.NotFound("二手商品不存在", nil)
		}
		if isAppError(err) {
			return err
		}
		return httpx.Internal("更新二手收藏状态失败", err)
	}

	return nil
}

func (s *Service) DeleteMarket(ctx context.Context, principal auth.Principal, id string) error {
	_, err := s.repository.UpdateMarket(ctx, id, func(item *cltypes.MarketItem) error {
		if item.PublisherUserID != principal.UserID {
			return httpx.Forbidden("只有发布者可以下架", nil)
		}
		normalized := normalizeReviewStatus(item.ReviewStatus)
		if normalized != statusPublished && normalized != statusReviewing && normalized != statusRejected {
			return httpx.BadRequest("当前状态不允许下架", nil)
		}
		item.ReviewStatus = statusOffline
		return nil
	})
	if err != nil {
		if errors.Is(err, clrepo.ErrNotFound) {
			return httpx.NotFound("二手商品不存在", nil)
		}
		if isAppError(err) {
			return err
		}
		return httpx.Internal("下架二手商品失败", err)
	}
	s.recordAudit(ctx, principal, "campus_life.market.delete", "market", id, "二手商品下架成功", nil)
	return nil
}

func (s *Service) ListErrands(ctx context.Context, principal auth.Principal, query cltypes.ErrandQuery) (map[string]any, error) {
	items, err := s.repository.ListErrands(ctx)
	if err != nil {
		return nil, httpx.Internal("读取跑腿列表失败", err)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.After(items[j].CreatedAt) })
	filtered := make([]map[string]any, 0)
	for _, item := range items {
		if !shouldExposeContent(principal, item.PublisherUserID, item.ReviewStatus, query.UserRole) {
			continue
		}
		if query.Category != "" && query.Category != item.Category {
			continue
		}
		if !matchKeyword(query.Keyword, item.Title, item.Desc) {
			continue
		}
		role := errandUserRole(item, principal)
		if query.UserRole != "" && role != query.UserRole {
			continue
		}
		status := mergeErrandStatus(item.Status, item.ReviewStatus)
		isOwner := role == "publisher"
		filtered = append(filtered, map[string]any{
			"id":                item.ID,
			"category":          item.Category,
			"title":             item.Title,
			"desc":              item.Desc,
			"route_start":       item.RouteStart,
			"route_end":         item.RouteEnd,
			"deadline":          item.Deadline.Format(time.RFC3339),
			"reward":            item.Reward,
			"status":            status,
			"user_role":         role,
			"is_owner":          isOwner,
			"is_accepted":       status == statusAccepted,
			"views":             0,
			"publisher":         item.Publisher,
			"publisher_initial": item.PublisherInitial,
			"created_at":        item.CreatedAt.Format(time.RFC3339),
			"can_edit":          canEditContent(isOwner, item.ReviewStatus),
			"can_delete":        canDeleteContent(isOwner, item.ReviewStatus),
			"can_accept":        !isOwner && status == statusPublished && principal.Authenticated,
			"can_cancel_accept": role == "acceptor" && status == statusAccepted,
		})
	}

	return listEnvelope(paginateMaps(filtered, query.Pagination), len(filtered), query.Pagination), nil
}

func (s *Service) GetErrandDetail(ctx context.Context, principal auth.Principal, id string) (map[string]any, error) {
	item, err := s.repository.GetErrand(ctx, id)
	if errors.Is(err, clrepo.ErrNotFound) {
		return nil, httpx.NotFound("跑腿任务不存在", nil)
	}
	if err != nil {
		return nil, httpx.Internal("读取跑腿详情失败", err)
	}
	if err := ensureContentVisible(principal, item.PublisherUserID, item.ReviewStatus, "跑腿任务不存在"); err != nil {
		return nil, err
	}
	role := errandUserRole(item, principal)
	isOwner := role == "publisher"
	canView := canViewContact(principal, item.PublisherUserID)
	status := mergeErrandStatus(item.Status, item.ReviewStatus)
	resolvedImages := resolveManagedURLs(ctx, s.storageProvider, item.Images)
	return map[string]any{
		"item": map[string]any{
			"id":                item.ID,
			"category":          item.Category,
			"title":             item.Title,
			"desc":              item.Desc,
			"route_start":       item.RouteStart,
			"route_end":         item.RouteEnd,
			"deadline":          item.Deadline.Format(time.RFC3339),
			"reward":            item.Reward,
			"contact":           visibleValue(canView, item.Contact),
			"status":            status,
			"is_accepted":       status == statusAccepted,
			"images":            resolvedImages,
			"publisher":         item.Publisher,
			"publisher_initial": item.PublisherInitial,
			"created_at":        item.CreatedAt.Format(time.RFC3339),
			"extra": map[string]any{
				"category":    item.Category,
				"route_start": item.RouteStart,
				"route_end":   item.RouteEnd,
				"deadline":    item.Deadline.Format(time.RFC3339),
				"reward":      item.Reward,
				"contact":     visibleValue(canView, item.Contact),
				"images":      resolvedImages,
				"status":      status,
				"urgent":      item.Urgent,
			},
		},
		"user_role":           role,
		"is_owner":            isOwner,
		"can_view_contact":    canView,
		"can_edit":            canEditContent(isOwner, item.ReviewStatus),
		"can_delete":          canDeleteContent(isOwner, item.ReviewStatus),
		"can_accept":          !isOwner && status == statusPublished && principal.Authenticated,
		"can_cancel_accept":   role == "acceptor" && status == statusAccepted,
	}, nil
}

func (s *Service) PublishErrand(ctx context.Context, principal auth.Principal, request cltypes.ErrandPublishRequest) (map[string]any, error) {
	if strings.TrimSpace(request.Title) == "" || strings.TrimSpace(request.Contact) == "" {
		return nil, httpx.BadRequest("跑腿标题和联系方式不能为空", nil)
	}
	deadline, err := time.Parse(time.RFC3339, request.Deadline)
	if err != nil {
		return nil, httpx.BadRequest("deadline 必须为 RFC3339 时间", nil)
	}
	id, err := s.repository.NextID(ctx, "errand")
	if err != nil {
		return nil, httpx.Internal("生成跑腿任务 ID 失败", err)
	}
	item := cltypes.ErrandItem{
		ID:               id,
		Title:            strings.TrimSpace(request.Title),
		Desc:             strings.TrimSpace(request.Desc),
		ReviewStatus:     statusReviewing,
		Category:         strings.TrimSpace(request.Category),
		RouteStart:       strings.TrimSpace(request.RouteStart),
		RouteEnd:         strings.TrimSpace(request.RouteEnd),
		Deadline:         deadline,
		Reward:           strings.TrimSpace(request.Reward),
		Contact:          strings.TrimSpace(request.Contact),
		Urgent:           request.Urgent,
		Images:           append([]string(nil), request.Images...),
		Status:           "published",
		PublisherUserID:  principal.UserID,
		Publisher:        displayName(principal),
		PublisherInitial: initialOf(displayName(principal)),
		CreatedAt:        time.Now().UTC(),
	}
	item, err = s.repository.SaveErrand(ctx, item)
	if err != nil {
		return nil, httpx.Internal("保存跑腿任务失败", err)
	}
	s.recordAudit(ctx, principal, "campus_life.errand.publish", "errand", item.ID, "跑腿任务发布成功", map[string]any{
		"review_status": item.ReviewStatus,
		"category":      item.Category,
	})
	return map[string]any{"id": item.ID}, nil
}

func (s *Service) AcceptErrand(ctx context.Context, principal auth.Principal, taskID string) error {
	_, err := s.repository.UpdateErrand(ctx, taskID, func(item *cltypes.ErrandItem) error {
		if item.PublisherUserID == principal.UserID {
			return httpx.BadRequest("不能接自己发布的任务", nil)
		}
		status := mergeErrandStatus(item.Status, item.ReviewStatus)
		if status != statusPublished {
			switch status {
			case statusReviewing:
				return httpx.BadRequest("该任务仍在审核中，暂不可接单", nil)
			case statusRejected:
				return httpx.BadRequest("该任务审核未通过，无法接单", nil)
			case statusOffline:
				return httpx.BadRequest("该任务已下线，无法接单", nil)
			case statusCancelled:
				return httpx.BadRequest("该任务已取消", nil)
			case statusAccepted:
				return httpx.BadRequest("该任务已被接单", nil)
			}
		}
		item.Status = statusAccepted
		item.AcceptorUserID = principal.UserID
		return nil
	})
	if err != nil {
		if err == clrepo.ErrNotFound {
			return httpx.NotFound("跑腿任务不存在", nil)
		}
		if isAppError(err) {
			return err
		}
		return httpx.Internal("更新跑腿接单状态失败", err)
	}

	return nil
}

func (s *Service) CancelErrandPublish(ctx context.Context, principal auth.Principal, taskID string) error {
	_, err := s.repository.UpdateErrand(ctx, taskID, func(item *cltypes.ErrandItem) error {
		if item.PublisherUserID != principal.UserID {
			return httpx.Forbidden("只有发布者可以取消发布", nil)
		}
		item.Status = "cancelled"
		return nil
	})
	if err != nil {
		if err == clrepo.ErrNotFound {
			return httpx.NotFound("跑腿任务不存在", nil)
		}
		if isAppError(err) {
			return err
		}
		return httpx.Internal("取消跑腿发布失败", err)
	}

	return nil
}

func (s *Service) CancelErrandAccept(ctx context.Context, principal auth.Principal, taskID string) error {
	_, err := s.repository.UpdateErrand(ctx, taskID, func(item *cltypes.ErrandItem) error {
		if item.AcceptorUserID != principal.UserID {
			return httpx.Forbidden("只有接单者可以取消接单", nil)
		}
		item.Status = "published"
		item.AcceptorUserID = ""
		return nil
	})
	if err != nil {
		if err == clrepo.ErrNotFound {
			return httpx.NotFound("跑腿任务不存在", nil)
		}
		if isAppError(err) {
			return err
		}
		return httpx.Internal("取消跑腿接单失败", err)
	}

	return nil
}

func (s *Service) ListResources(ctx context.Context, principal auth.Principal, query cltypes.ResourceQuery) (map[string]any, error) {
	items, err := s.repository.ListResources(ctx)
	if err != nil {
		return nil, httpx.Internal("读取资料列表失败", err)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.After(items[j].CreatedAt) })
	filtered := make([]map[string]any, 0)
	for _, item := range items {
		if normalizeReviewStatus(item.ReviewStatus) != statusPublished {
			continue
		}
		if query.Category != "" && query.Category != item.Extra.Category {
			continue
		}
		if !matchKeyword(query.Keyword, item.Title, item.Desc) {
			continue
		}
		isOwner := item.PublisherUserID == principal.UserID
		canView := canViewContact(principal, item.PublisherUserID)
		resolvedFiles := resolveResourceFiles(ctx, s.storageProvider, item.Extra.Files)
		filtered = append(filtered, map[string]any{
			"id":                item.ID,
			"title":             item.Title,
			"desc":              item.Desc,
			"publisher":         item.Publisher,
			"publisher_initial": item.PublisherInitial,
			"created_at":        item.CreatedAt.Format(time.RFC3339),
			"status":            normalizeReviewStatus(item.ReviewStatus),
			"is_owner":          isOwner,
			"can_edit":          canEditContent(isOwner, item.ReviewStatus),
			"can_delete":        canDeleteContent(isOwner, item.ReviewStatus),
			"can_download":      !isOwner && canView && principal.Authenticated,
			"extra": map[string]any{
				"category":     item.Extra.Category,
				"course_name":  item.Extra.CourseName,
				"contact":      visibleValue(canView, item.Extra.Contact),
				"files":        resolvedFiles,
				"file_size":    item.Extra.FileSize,
				"file_type":    item.Extra.FileType,
				"download_url": firstResourceURL(resolvedFiles, resolveManagedURL(ctx, s.storageProvider, item.Extra.DownloadURL)),
				"likes":        item.Extra.Likes,
				"views":        item.Extra.Views,
			},
		})
	}
	return listEnvelope(paginateMaps(filtered, query.Pagination), len(filtered), query.Pagination), nil
}

func (s *Service) GetResourceDetail(ctx context.Context, principal auth.Principal, id string) (map[string]any, error) {
	item, err := s.repository.GetResource(ctx, id)
	if errors.Is(err, clrepo.ErrNotFound) {
		return nil, httpx.NotFound("资料不存在", nil)
	}
	if err != nil {
		return nil, httpx.Internal("读取资料详情失败", err)
	}
	if err := ensureContentVisible(principal, item.PublisherUserID, item.ReviewStatus, "资料不存在"); err != nil {
		return nil, err
	}
	canView := canViewContact(principal, item.PublisherUserID)
	role := simpleUserRole(item.PublisherUserID, principal)
	isOwner := role == "publisher"
	resolvedFiles := resolveResourceFiles(ctx, s.storageProvider, item.Extra.Files)
	return map[string]any{
		"id":                item.ID,
		"title":             item.Title,
		"desc":              item.Desc,
		"publisher":         item.Publisher,
		"publisher_initial": item.PublisherInitial,
		"created_at":        item.CreatedAt.Format(time.RFC3339),
		"status":            normalizeReviewStatus(item.ReviewStatus),
		"user_role":         role,
		"is_owner":          isOwner,
		"can_view_contact":  canView,
		"can_edit":          canEditContent(isOwner, item.ReviewStatus),
		"can_delete":        canDeleteContent(isOwner, item.ReviewStatus),
		"can_download":      principal.Authenticated && canView,
		"extra": map[string]any{
			"category":     item.Extra.Category,
			"course_name":  item.Extra.CourseName,
			"contact":      visibleValue(canView, item.Extra.Contact),
			"files":        resolvedFiles,
			"file_size":    item.Extra.FileSize,
			"file_type":    item.Extra.FileType,
			"download_url": firstResourceURL(resolvedFiles, resolveManagedURL(ctx, s.storageProvider, item.Extra.DownloadURL)),
			"likes":        item.Extra.Likes,
			"views":        item.Extra.Views,
		},
	}, nil
}

func (s *Service) PublishResource(ctx context.Context, principal auth.Principal, request cltypes.ResourcePublishRequest) (map[string]any, error) {
	if strings.TrimSpace(request.Title) == "" || len(request.FilePaths) == 0 {
		return nil, httpx.BadRequest("资料标题和文件路径不能为空", nil)
	}
	files := make([]cltypes.ResourceFile, 0, len(request.FilePaths))
	for _, filePath := range request.FilePaths {
		filePath = strings.TrimSpace(filePath)
		if filePath == "" {
			continue
		}
		files = append(files, cltypes.ResourceFile{
			Name:     path.Base(filePath),
			Path:     strings.TrimPrefix(filePath, "/"),
			FileType: detectFileType(filePath),
			FileSize: "1.0MB",
		})
	}
	id, err := s.repository.NextID(ctx, "resource")
	if err != nil {
		return nil, httpx.Internal("生成资料 ID 失败", err)
	}
	item := cltypes.ResourceItem{
		ID:               id,
		Title:            strings.TrimSpace(request.Title),
		Desc:             firstNonEmpty(strings.TrimSpace(request.Desc), strings.TrimSpace(request.Title)),
		ReviewStatus:     statusReviewing,
		PublisherUserID:  principal.UserID,
		Publisher:        displayName(principal),
		PublisherInitial: initialOf(displayName(principal)),
		CreatedAt:        time.Now().UTC(),
		Extra: cltypes.ResourceExtra{
			Category:   strings.TrimSpace(request.Category),
			CourseName: strings.TrimSpace(request.CourseName),
			Contact:    strings.TrimSpace(request.Contact),
			Files:      files,
			FileSize:   firstFileSize(files),
			FileType:   firstFileType(files),
			Views:      0,
			Likes:      0,
		},
	}
	item, err = s.repository.SaveResource(ctx, item)
	if err != nil {
		return nil, httpx.Internal("保存资料失败", err)
	}
	s.recordAudit(ctx, principal, "campus_life.resource.publish", "resource", item.ID, "资料发布成功", map[string]any{
		"review_status": item.ReviewStatus,
		"category":      item.Extra.Category,
		"course_name":   item.Extra.CourseName,
	})
	return map[string]any{"id": item.ID}, nil
}

func (s *Service) DeleteResource(ctx context.Context, principal auth.Principal, id string) error {
	_, err := s.repository.UpdateResource(ctx, id, func(item *cltypes.ResourceItem) error {
		if item.PublisherUserID != principal.UserID {
			return httpx.Forbidden("只有发布者可以下架", nil)
		}
		normalized := normalizeReviewStatus(item.ReviewStatus)
		if normalized != statusPublished && normalized != statusReviewing && normalized != statusRejected {
			return httpx.BadRequest("当前状态不允许下架", nil)
		}
		item.ReviewStatus = statusOffline
		return nil
	})
	if err != nil {
		if errors.Is(err, clrepo.ErrNotFound) {
			return httpx.NotFound("资料不存在", nil)
		}
		if isAppError(err) {
			return err
		}
		return httpx.Internal("下架资料失败", err)
	}
	s.recordAudit(ctx, principal, "campus_life.resource.delete", "resource", id, "资料下架成功", nil)
	return nil
}

func (s *Service) ListLostFound(ctx context.Context, principal auth.Principal, query cltypes.LostFoundQuery) (map[string]any, error) {
	items, err := s.repository.ListLostFound(ctx)
	if err != nil {
		return nil, httpx.Internal("读取失物招领列表失败", err)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.After(items[j].CreatedAt) })
	filtered := make([]map[string]any, 0)
	for _, item := range items {
		if !shouldExposeContent(principal, item.PublisherUserID, item.ReviewStatus, "") {
			continue
		}
		if query.Category != "" && query.Category != item.Extra.Category {
			continue
		}
		if query.Type != "" && query.Type != item.Extra.Type {
			continue
		}
		if !matchKeyword(query.Keyword, item.Title, item.Desc) {
			continue
		}
		canView := canViewContact(principal, item.PublisherUserID)
		isOwner := item.PublisherUserID == principal.UserID
		filtered = append(filtered, map[string]any{
			"id":                item.ID,
			"title":             item.Title,
			"desc":              item.Desc,
			"publisher":         item.Publisher,
			"publisher_initial": item.PublisherInitial,
			"created_at":        item.CreatedAt.Format(time.RFC3339),
			"status":            normalizeReviewStatus(item.ReviewStatus),
			"is_owner":          isOwner,
			"can_edit":          canEditContent(isOwner, item.ReviewStatus),
			"can_delete":        canDeleteContent(isOwner, item.ReviewStatus),
			"can_mark_resolved": isOwner && normalizeReviewStatus(item.ReviewStatus) == statusPublished,
			"extra": map[string]any{
				"type":         item.Extra.Type,
				"category":     item.Extra.Category,
				"location":     item.Extra.Location,
				"event_time":   item.Extra.EventTime,
				"item_feature": item.Extra.ItemFeature,
				"contact":      visibleValue(canView, item.Extra.Contact),
			},
		})
	}
	return listEnvelope(paginateMaps(filtered, query.Pagination), len(filtered), query.Pagination), nil
}

func (s *Service) GetLostFoundDetail(ctx context.Context, principal auth.Principal, id string) (map[string]any, error) {
	item, err := s.repository.GetLostFound(ctx, id)
	if errors.Is(err, clrepo.ErrNotFound) {
		return nil, httpx.NotFound("失物招领不存在", nil)
	}
	if err != nil {
		return nil, httpx.Internal("读取失物招领详情失败", err)
	}
	if err := ensureContentVisible(principal, item.PublisherUserID, item.ReviewStatus, "失物招领不存在"); err != nil {
		return nil, err
	}
	canView := canViewContact(principal, item.PublisherUserID)
	role := simpleUserRole(item.PublisherUserID, principal)
	isOwner := role == "publisher"
	status := normalizeReviewStatus(item.ReviewStatus)
	return map[string]any{
		"id":                 item.ID,
		"title":              item.Title,
		"desc":               item.Desc,
		"publisher":          item.Publisher,
		"publisher_initial":  item.PublisherInitial,
		"created_at":         item.CreatedAt.Format(time.RFC3339),
		"status":             status,
		"user_role":          role,
		"is_owner":           isOwner,
		"can_view_contact":   canView,
		"can_edit":           canEditContent(isOwner, item.ReviewStatus),
		"can_delete":         canDeleteContent(isOwner, item.ReviewStatus),
		"can_mark_resolved":  isOwner && status == statusPublished,
		"extra": map[string]any{
			"type":         item.Extra.Type,
			"category":     item.Extra.Category,
			"location":     item.Extra.Location,
			"event_time":   item.Extra.EventTime,
			"item_feature": item.Extra.ItemFeature,
			"contact":      visibleValue(canView, item.Extra.Contact),
		},
	}, nil
}

func (s *Service) PublishLostFound(ctx context.Context, principal auth.Principal, request cltypes.LostFoundPublishRequest) (map[string]any, error) {
	if strings.TrimSpace(request.Title) == "" || strings.TrimSpace(request.Contact) == "" {
		return nil, httpx.BadRequest("标题和联系方式不能为空", nil)
	}
	id, err := s.repository.NextID(ctx, "lostfound")
	if err != nil {
		return nil, httpx.Internal("生成失物招领 ID 失败", err)
	}
	item := cltypes.LostFoundItem{
		ID:               id,
		Title:            strings.TrimSpace(request.Title),
		Desc:             strings.TrimSpace(request.Desc),
		ReviewStatus:     statusReviewing,
		PublisherUserID:  principal.UserID,
		Publisher:        displayName(principal),
		PublisherInitial: initialOf(displayName(principal)),
		CreatedAt:        time.Now().UTC(),
		Extra: cltypes.LostFoundExtra{
			Type:        strings.TrimSpace(request.Type),
			Category:    strings.TrimSpace(request.Category),
			Location:    strings.TrimSpace(request.Location),
			EventTime:   strings.TrimSpace(request.EventTime),
			ItemFeature: strings.TrimSpace(request.ItemFeature),
			Contact:     strings.TrimSpace(request.Contact),
		},
	}
	item, err = s.repository.SaveLostFound(ctx, item)
	if err != nil {
		return nil, httpx.Internal("保存失物招领失败", err)
	}
	s.recordAudit(ctx, principal, "campus_life.lost_found.publish", "lostFound", item.ID, "失物招领发布成功", map[string]any{
		"review_status": item.ReviewStatus,
		"type":          item.Extra.Type,
		"category":      item.Extra.Category,
	})
	return map[string]any{"id": item.ID}, nil
}

func (s *Service) DeleteLostFound(ctx context.Context, principal auth.Principal, id string) error {
	_, err := s.repository.UpdateLostFound(ctx, id, func(item *cltypes.LostFoundItem) error {
		if item.PublisherUserID != principal.UserID {
			return httpx.Forbidden("只有发布者可以下架", nil)
		}
		normalized := normalizeReviewStatus(item.ReviewStatus)
		if normalized != statusPublished && normalized != statusReviewing && normalized != statusRejected {
			return httpx.BadRequest("当前状态不允许下架", nil)
		}
		item.ReviewStatus = statusOffline
		return nil
	})
	if err != nil {
		if errors.Is(err, clrepo.ErrNotFound) {
			return httpx.NotFound("失物招领不存在", nil)
		}
		if isAppError(err) {
			return err
		}
		return httpx.Internal("下架失物招领失败", err)
	}
	s.recordAudit(ctx, principal, "campus_life.lost_found.delete", "lostFound", id, "失物招领下架成功", nil)
	return nil
}

func (s *Service) MarkLostFoundResolved(ctx context.Context, principal auth.Principal, id string) error {
	_, err := s.repository.UpdateLostFound(ctx, id, func(item *cltypes.LostFoundItem) error {
		if item.PublisherUserID != principal.UserID {
			return httpx.Forbidden("只有发布者可以标记已找到", nil)
		}
		if normalizeReviewStatus(item.ReviewStatus) != statusPublished {
			return httpx.BadRequest("当前状态不允许标记已找到", nil)
		}
		item.ReviewStatus = statusResolved
		return nil
	})
	if err != nil {
		if errors.Is(err, clrepo.ErrNotFound) {
			return httpx.NotFound("失物招领不存在", nil)
		}
		if isAppError(err) {
			return err
		}
		return httpx.Internal("标记失物招领已找到失败", err)
	}
	s.recordAudit(ctx, principal, "campus_life.lost_found.resolve", "lostFound", id, "失物招领标记已找到", nil)
	return nil
}

func (s *Service) ListCarpools(ctx context.Context, principal auth.Principal, query cltypes.CarpoolQuery) (map[string]any, error) {
	items, err := s.repository.ListCarpools(ctx)
	if err != nil {
		return nil, httpx.Internal("读取拼车列表失败", err)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.After(items[j].CreatedAt) })

	now := time.Now().In(chinaLocation)
	filtered := make([]map[string]any, 0)
	for _, item := range items {
		if normalizeReviewStatus(item.ReviewStatus) != statusPublished {
			continue
		}
		if query.Category != "" && query.Category != "all" && normalizedCarpoolCategory(item, now) != query.Category {
			continue
		}
		if !matchKeyword(query.Keyword, item.From, item.To, item.Note, item.Publisher, item.Type) {
			continue
		}
		canView := canViewContact(principal, item.PublisherUserID)
		isOwner := item.PublisherUserID == principal.UserID
		payload := buildCarpoolPayload(item, canView, now)
		payload["is_owner"] = isOwner
		payload["can_edit"] = canEditContent(isOwner, item.ReviewStatus) && item.TravelAt.After(now.UTC())
		payload["can_delete"] = canDeleteContent(isOwner, item.ReviewStatus)
		payload["can_join_carpool"] = !isOwner && normalizeReviewStatus(item.ReviewStatus) == statusPublished && principal.Authenticated
		filtered = append(filtered, payload)
	}

	return listEnvelope(paginateMaps(filtered, query.Pagination), len(filtered), query.Pagination), nil
}

func (s *Service) GetCarpoolDetail(ctx context.Context, principal auth.Principal, id string) (map[string]any, error) {
	item, err := s.repository.GetCarpool(ctx, id)
	if errors.Is(err, clrepo.ErrNotFound) {
		return nil, httpx.NotFound("拼车行程不存在", nil)
	}
	if err != nil {
		return nil, httpx.Internal("读取拼车详情失败", err)
	}
	if err := ensureContentVisible(principal, item.PublisherUserID, item.ReviewStatus, "拼车行程不存在"); err != nil {
		return nil, err
	}

	canView := canViewContact(principal, item.PublisherUserID)
	role := simpleUserRole(item.PublisherUserID, principal)
	isOwner := role == "publisher"
	now := time.Now().In(chinaLocation)
	status := normalizeReviewStatus(item.ReviewStatus)
	payload := buildCarpoolPayload(item, canView, now)
	payload["can_view_contact"] = canView
	payload["user_role"] = role
	payload["is_owner"] = isOwner
	payload["can_edit"] = canEditContent(isOwner, item.ReviewStatus) && item.TravelAt.After(now.UTC())
	payload["can_delete"] = canDeleteContent(isOwner, item.ReviewStatus)
	payload["can_join_carpool"] = !isOwner && status == statusPublished && principal.Authenticated
	return payload, nil
}

func (s *Service) PublishCarpool(ctx context.Context, principal auth.Principal, request cltypes.CarpoolPublishRequest) (map[string]any, error) {
	if strings.TrimSpace(request.From) == "" || strings.TrimSpace(request.To) == "" {
		return nil, httpx.BadRequest("出发地和目的地不能为空", nil)
	}
	if strings.TrimSpace(request.TravelDate) == "" || strings.TrimSpace(request.TravelTime) == "" {
		return nil, httpx.BadRequest("travel_date 和 travel_time 不能为空", nil)
	}
	if strings.TrimSpace(request.Contact) == "" {
		return nil, httpx.BadRequest("联系方式不能为空", nil)
	}

	travelAt, err := parseCarpoolTravelAt(request.TravelDate, request.TravelTime)
	if err != nil {
		return nil, httpx.BadRequest("travel_date/travel_time 格式错误", nil)
	}
	now := time.Now().In(chinaLocation)
	category := firstNonEmpty(strings.TrimSpace(request.Category), normalizedCarpoolCategory(cltypes.CarpoolItem{
		TravelAt: travelAt,
	}, now))
	if !isSupportedCarpoolCategory(category) {
		category = normalizedCarpoolCategory(cltypes.CarpoolItem{TravelAt: travelAt}, now)
	}
	id, err := s.repository.NextID(ctx, "carpool")
	if err != nil {
		return nil, httpx.Internal("生成拼车 ID 失败", err)
	}

	item := cltypes.CarpoolItem{
		ID:               id,
		Category:         category,
		From:             strings.TrimSpace(request.From),
		To:               strings.TrimSpace(request.To),
		TravelAt:         travelAt,
		Type:             firstNonEmpty(strings.TrimSpace(request.Type), defaultCarpoolType(category)),
		SeatsText:        strings.TrimSpace(request.SeatsText),
		Price:            strings.TrimSpace(request.Price),
		Note:             strings.TrimSpace(request.Note),
		Tags:             sanitizeTags(request.Tags),
		Contact:          strings.TrimSpace(request.Contact),
		ReviewStatus:     statusReviewing,
		PublisherUserID:  principal.UserID,
		Publisher:        displayName(principal),
		PublisherInitial: initialOf(displayName(principal)),
		CreatedAt:        time.Now().UTC(),
	}
	item, err = s.repository.SaveCarpool(ctx, item)
	if err != nil {
		return nil, httpx.Internal("保存拼车信息失败", err)
	}

	s.recordAudit(ctx, principal, "campus_life.carpool.publish", "carpool", item.ID, "拼车信息发布成功", map[string]any{
		"review_status": item.ReviewStatus,
		"category":      item.Category,
	})

	return map[string]any{"id": item.ID}, nil
}

func (s *Service) DeleteCarpool(ctx context.Context, principal auth.Principal, id string) error {
	_, err := s.repository.UpdateCarpool(ctx, id, func(item *cltypes.CarpoolItem) error {
		if item.PublisherUserID != principal.UserID {
			return httpx.Forbidden("只有发布者可以取消发布", nil)
		}
		normalized := normalizeReviewStatus(item.ReviewStatus)
		if normalized != statusPublished && normalized != statusReviewing && normalized != statusRejected {
			return httpx.BadRequest("当前状态不允许取消发布", nil)
		}
		item.ReviewStatus = statusOffline
		return nil
	})
	if err != nil {
		if errors.Is(err, clrepo.ErrNotFound) {
			return httpx.NotFound("拼车行程不存在", nil)
		}
		if isAppError(err) {
			return err
		}
		return httpx.Internal("取消拼车发布失败", err)
	}
	s.recordAudit(ctx, principal, "campus_life.carpool.delete", "carpool", id, "拼车取消发布成功", nil)
	return nil
}

func (s *Service) ListMeetups(ctx context.Context, principal auth.Principal, query cltypes.MeetupQuery) (map[string]any, error) {
	items, err := s.repository.ListMeetups(ctx)
	if err != nil {
		return nil, httpx.Internal("读取组局列表失败", err)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.After(items[j].CreatedAt) })

	filtered := make([]map[string]any, 0)
	now := time.Now().In(chinaLocation)
	for _, item := range items {
		if !matchMeetupUserRole(item, principal, query.UserRole) {
			continue
		}
		if !shouldExposeContent(principal, item.PublisherUserID, item.ReviewStatus, query.UserRole) {
			continue
		}
		if !shouldExposeMeetupState(principal, item, query.UserRole) {
			continue
		}
		if query.Category != "" && query.Category != "all" && query.Category != item.Category {
			continue
		}
		if !matchKeyword(query.Keyword, item.Title, item.Desc, item.Location, item.Publisher) {
			continue
		}
		payload := buildMeetupPayload(item, principal, now)
		payload["is_owner"] = meetupUserRole(item, principal) == "publisher"
		payload["can_edit"] = canEditContent(meetupUserRole(item, principal) == "publisher", item.ReviewStatus) && mergeMeetupStatus(item.Status, item.ReviewStatus) != statusCancelled
		payload["can_delete"] = meetupUserRole(item, principal) == "publisher" && mergeMeetupStatus(item.Status, item.ReviewStatus) != statusCancelled
		filtered = append(filtered, payload)
	}

	return listEnvelope(paginateMaps(filtered, query.Pagination), len(filtered), query.Pagination), nil
}

func (s *Service) GetMeetupDetail(ctx context.Context, principal auth.Principal, id string) (map[string]any, error) {
	item, err := s.repository.GetMeetup(ctx, id)
	if errors.Is(err, clrepo.ErrNotFound) {
		return nil, httpx.NotFound("组局不存在", nil)
	}
	if err != nil {
		return nil, httpx.Internal("读取组局详情失败", err)
	}
	if err := ensureContentVisible(principal, item.PublisherUserID, item.ReviewStatus, "组局不存在"); err != nil {
		return nil, err
	}
	if !shouldExposeMeetupState(principal, item, "") {
		return nil, httpx.NotFound("组局不存在", nil)
	}

	payload := buildMeetupPayload(item, principal, time.Now().In(chinaLocation))
	payload["can_view_contact"] = canViewContact(principal, item.PublisherUserID)
	payload["can_edit"] = canEditContent(meetupUserRole(item, principal) == "publisher", item.ReviewStatus) && mergeMeetupStatus(item.Status, item.ReviewStatus) != statusCancelled
	payload["can_delete"] = meetupUserRole(item, principal) == "publisher" && mergeMeetupStatus(item.Status, item.ReviewStatus) != statusCancelled
	return payload, nil
}

func (s *Service) PublishMeetup(ctx context.Context, principal auth.Principal, request cltypes.MeetupPublishRequest) (map[string]any, error) {
	if strings.TrimSpace(request.Title) == "" || strings.TrimSpace(request.Location) == "" || strings.TrimSpace(request.Contact) == "" {
		return nil, httpx.BadRequest("标题、地点和联系方式不能为空", nil)
	}
	if request.MaxParticipants <= 1 {
		return nil, httpx.BadRequest("max_participants 至少为 2", nil)
	}

	startAt, err := time.Parse(time.RFC3339, strings.TrimSpace(request.StartAt))
	if err != nil {
		return nil, httpx.BadRequest("start_at 必须为 RFC3339 时间", nil)
	}
	deadlineAt := startAt
	if strings.TrimSpace(request.DeadlineAt) != "" {
		deadlineAt, err = time.Parse(time.RFC3339, strings.TrimSpace(request.DeadlineAt))
		if err != nil {
			return nil, httpx.BadRequest("deadline_at 必须为 RFC3339 时间", nil)
		}
	}
	if deadlineAt.After(startAt) {
		return nil, httpx.BadRequest("deadline_at 不能晚于 start_at", nil)
	}

	id, err := s.repository.NextID(ctx, "meetup")
	if err != nil {
		return nil, httpx.Internal("生成组局 ID 失败", err)
	}

	item := cltypes.MeetupItem{
		ID:               id,
		Category:         strings.TrimSpace(request.Category),
		Title:            strings.TrimSpace(request.Title),
		Desc:             firstNonEmpty(strings.TrimSpace(request.Desc), strings.TrimSpace(request.Title)),
		Location:         strings.TrimSpace(request.Location),
		StartAt:          startAt,
		DeadlineAt:       deadlineAt,
		MaxParticipants:  request.MaxParticipants,
		FeeText:          strings.TrimSpace(request.FeeText),
		Tags:             sanitizeTags(request.Tags),
		Contact:          strings.TrimSpace(request.Contact),
		Status:           "open",
		ReviewStatus:     statusReviewing,
		PublisherUserID:  principal.UserID,
		Publisher:        displayName(principal),
		PublisherInitial: initialOf(displayName(principal)),
		CreatedAt:        time.Now().UTC(),
	}
	item = refreshMeetupStatus(item)
	item, err = s.repository.SaveMeetup(ctx, item)
	if err != nil {
		return nil, httpx.Internal("保存组局失败", err)
	}

	s.recordAudit(ctx, principal, "campus_life.meetup.publish", "meetup", item.ID, "组局发布成功", map[string]any{
		"review_status":    item.ReviewStatus,
		"category":         item.Category,
		"max_participants": item.MaxParticipants,
	})

	return map[string]any{"id": item.ID}, nil
}

func (s *Service) JoinMeetup(ctx context.Context, principal auth.Principal, meetupID string) error {
	_, err := s.repository.UpdateMeetup(ctx, meetupID, func(item *cltypes.MeetupItem) error {
		if item.PublisherUserID == principal.UserID {
			return httpx.BadRequest("不能报名自己发起的组局", nil)
		}
		status := mergeMeetupStatus(item.Status, item.ReviewStatus)
		if status != statusOpen {
			switch status {
			case statusReviewing:
				return httpx.BadRequest("该组局仍在审核中，暂不可报名", nil)
			case statusRejected:
				return httpx.BadRequest("该组局审核未通过，无法报名", nil)
			case statusOffline:
				return httpx.BadRequest("该组局已下线，无法报名", nil)
			case statusCancelled:
				return httpx.BadRequest("该组局已取消", nil)
			case statusFull:
				return httpx.BadRequest("该组局人数已满", nil)
			}
		}
		now := time.Now().UTC()
		if !item.DeadlineAt.IsZero() && item.DeadlineAt.Before(now) {
			return httpx.BadRequest("该组局报名已截止", nil)
		}
		if !item.StartAt.IsZero() && item.StartAt.Before(now) {
			return httpx.BadRequest("该组局已开始，无法再报名", nil)
		}
		if slices.Contains(item.ParticipantUserIDs, principal.UserID) {
			return httpx.BadRequest("你已经报名过该组局", nil)
		}
		item.ParticipantUserIDs = append(item.ParticipantUserIDs, principal.UserID)
		*item = refreshMeetupStatus(*item)
		return nil
	})
	if err != nil {
		if err == clrepo.ErrNotFound {
			return httpx.NotFound("组局不存在", nil)
		}
		if isAppError(err) {
			return err
		}
		return httpx.Internal("报名组局失败", err)
	}

	return nil
}

func (s *Service) CancelMeetupJoin(ctx context.Context, principal auth.Principal, meetupID string) error {
	_, err := s.repository.UpdateMeetup(ctx, meetupID, func(item *cltypes.MeetupItem) error {
		index := slices.Index(item.ParticipantUserIDs, principal.UserID)
		if index < 0 {
			return httpx.BadRequest("你尚未报名该组局", nil)
		}
		item.ParticipantUserIDs = append(item.ParticipantUserIDs[:index], item.ParticipantUserIDs[index+1:]...)
		*item = refreshMeetupStatus(*item)
		return nil
	})
	if err != nil {
		if err == clrepo.ErrNotFound {
			return httpx.NotFound("组局不存在", nil)
		}
		if isAppError(err) {
			return err
		}
		return httpx.Internal("取消报名组局失败", err)
	}

	return nil
}

func (s *Service) CancelMeetupPublish(ctx context.Context, principal auth.Principal, meetupID string) error {
	_, err := s.repository.UpdateMeetup(ctx, meetupID, func(item *cltypes.MeetupItem) error {
		if item.PublisherUserID != principal.UserID {
			return httpx.Forbidden("只有发起人可以取消组局", nil)
		}
		item.Status = "cancelled"
		return nil
	})
	if err != nil {
		if err == clrepo.ErrNotFound {
			return httpx.NotFound("组局不存在", nil)
		}
		if isAppError(err) {
			return err
		}
		return httpx.Internal("取消组局失败", err)
	}

	return nil
}

func (s *Service) ListReviewQueue(ctx context.Context, query cltypes.ReviewQuery) (map[string]any, error) {
	type row struct {
		createdAt time.Time
		payload   map[string]any
	}

	rows := make([]row, 0)
	appendRow := func(contentType string, createdAt time.Time, reviewStatus string, title, desc, publisher, id string, extra map[string]any) {
		if !matchReviewQuery(query, contentType, reviewStatus, title, desc, publisher, id) {
			return
		}
		rows = append(rows, row{
			createdAt: createdAt,
			payload: map[string]any{
				"content_type":  contentType,
				"content_id":    id,
				"title":         title,
				"desc":          desc,
				"publisher":     publisher,
				"created_at":    createdAt.Format(time.RFC3339),
				"review_status": normalizeReviewStatus(reviewStatus),
				"extra":         extra,
			},
		})
	}

	markets, err := s.repository.ListMarkets(ctx)
	if err != nil {
		return nil, httpx.Internal("读取二手审核列表失败", err)
	}
	for _, item := range markets {
		appendRow("market", item.CreatedAt, item.ReviewStatus, item.Title, item.Desc, item.Publisher, item.ID, map[string]any{
			"category":       item.Extra.Category,
			"price":          item.Extra.Price,
			"original_price": item.Extra.OriginalPrice,
			"condition":      item.Extra.Condition,
			"trade_mode":     item.Extra.TradeMode,
		})
	}

	errands, err := s.repository.ListErrands(ctx)
	if err != nil {
		return nil, httpx.Internal("读取跑腿审核列表失败", err)
	}
	for _, item := range errands {
		appendRow("errand", item.CreatedAt, item.ReviewStatus, item.Title, item.Desc, item.Publisher, item.ID, map[string]any{
			"category":    item.Category,
			"status":      item.Status,
			"route_start": item.RouteStart,
			"route_end":   item.RouteEnd,
			"deadline":    item.Deadline.Format(time.RFC3339),
			"reward":      item.Reward,
		})
	}

	resources, err := s.repository.ListResources(ctx)
	if err != nil {
		return nil, httpx.Internal("读取资料审核列表失败", err)
	}
	for _, item := range resources {
		appendRow("resource", item.CreatedAt, item.ReviewStatus, item.Title, item.Desc, item.Publisher, item.ID, map[string]any{
			"category":    item.Extra.Category,
			"course_name": item.Extra.CourseName,
			"file_type":   item.Extra.FileType,
			"file_size":   item.Extra.FileSize,
		})
	}

	lostFounds, err := s.repository.ListLostFound(ctx)
	if err != nil {
		return nil, httpx.Internal("读取失物招领审核列表失败", err)
	}
	for _, item := range lostFounds {
		appendRow("lostFound", item.CreatedAt, item.ReviewStatus, item.Title, item.Desc, item.Publisher, item.ID, map[string]any{
			"category":   item.Extra.Category,
			"type":       item.Extra.Type,
			"location":   item.Extra.Location,
			"event_time": item.Extra.EventTime,
		})
	}

	carpools, err := s.repository.ListCarpools(ctx)
	if err != nil {
		return nil, httpx.Internal("读取拼车审核列表失败", err)
	}
	now := time.Now().In(chinaLocation)
	for _, item := range carpools {
		appendRow("carpool", item.CreatedAt, item.ReviewStatus, carpoolTitle(item), item.Note, item.Publisher, item.ID, map[string]any{
			"category":   normalizedCarpoolCategory(item, now),
			"from":       item.From,
			"to":         item.To,
			"time":       formatCarpoolTravelText(item.TravelAt, now),
			"seats_text": item.SeatsText,
		})
	}
	meetups, err := s.repository.ListMeetups(ctx)
	if err != nil {
		return nil, httpx.Internal("读取组局审核列表失败", err)
	}
	for _, item := range meetups {
		appendRow("meetup", item.CreatedAt, item.ReviewStatus, item.Title, item.Desc, item.Publisher, item.ID, map[string]any{
			"category":         item.Category,
			"location":         item.Location,
			"start_at":         item.StartAt.In(chinaLocation).Format(time.RFC3339),
			"max_participants": item.MaxParticipants,
			"joined_count":     meetupJoinedCount(item),
			"status":           normalizeMeetupStatus(item.Status),
		})
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].createdAt.After(rows[j].createdAt)
	})

	items := paginateRows(rows, query.Pagination)
	list := make([]map[string]any, 0, len(items))
	for _, item := range items {
		list = append(list, item.payload)
	}

	return listEnvelope(list, len(rows), query.Pagination), nil
}

func (s *Service) UpdateReviewStatus(ctx context.Context, principal auth.Principal, request cltypes.ReviewUpdateRequest) error {
	contentType := strings.TrimSpace(request.ContentType)
	contentID := strings.TrimSpace(request.ContentID)
	reviewStatus := normalizeReviewStatus(request.ReviewStatus)

	if contentType == "" || contentID == "" {
		return httpx.BadRequest("content_type 和 content_id 不能为空", nil)
	}
	if !isSupportedReviewStatus(request.ReviewStatus) {
		return httpx.BadRequest("review_status 仅支持 reviewing/published/rejected/offline", nil)
	}

	var err error
	switch contentType {
	case "market":
		_, err = s.repository.UpdateMarket(ctx, contentID, func(item *cltypes.MarketItem) error {
			item.ReviewStatus = reviewStatus
			return nil
		})
	case "errand":
		_, err = s.repository.UpdateErrand(ctx, contentID, func(item *cltypes.ErrandItem) error {
			item.ReviewStatus = reviewStatus
			return nil
		})
	case "resource":
		_, err = s.repository.UpdateResource(ctx, contentID, func(item *cltypes.ResourceItem) error {
			item.ReviewStatus = reviewStatus
			return nil
		})
	case "lostFound":
		_, err = s.repository.UpdateLostFound(ctx, contentID, func(item *cltypes.LostFoundItem) error {
			item.ReviewStatus = reviewStatus
			return nil
		})
	case "carpool":
		_, err = s.repository.UpdateCarpool(ctx, contentID, func(item *cltypes.CarpoolItem) error {
			item.ReviewStatus = reviewStatus
			return nil
		})
	case "meetup":
		_, err = s.repository.UpdateMeetup(ctx, contentID, func(item *cltypes.MeetupItem) error {
			item.ReviewStatus = reviewStatus
			return nil
		})
	default:
		return httpx.BadRequest("content_type 仅支持 market/errand/resource/lostFound/carpool/meetup", nil)
	}

	if errors.Is(err, clrepo.ErrNotFound) {
		return httpx.NotFound("待审核内容不存在", nil)
	}
	if err != nil {
		if isAppError(err) {
			return err
		}
		return httpx.Internal("更新审核状态失败", err)
	}

	s.recordAudit(ctx, principal, "campus_life.review.update", contentType, contentID, "校园生活审核状态更新成功", map[string]any{
		"review_status": reviewStatus,
	})

	return nil
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

func errandUserRole(item cltypes.ErrandItem, principal auth.Principal) string {
	if !principal.Authenticated {
		return "viewer"
	}
	if item.PublisherUserID == principal.UserID {
		return "publisher"
	}
	if item.AcceptorUserID == principal.UserID {
		return "acceptor"
	}
	return "viewer"
}

func meetupUserRole(item cltypes.MeetupItem, principal auth.Principal) string {
	if !principal.Authenticated {
		return "viewer"
	}
	if item.PublisherUserID == principal.UserID {
		return "publisher"
	}
	if slices.Contains(item.ParticipantUserIDs, principal.UserID) {
		return "participant"
	}
	return "viewer"
}

func canEditContent(isOwner bool, reviewStatus string) bool {
	if !isOwner {
		return false
	}
	normalized := normalizeReviewStatus(reviewStatus)
	return normalized == statusPublished || normalized == statusReviewing
}

func canDeleteContent(isOwner bool, reviewStatus string, extraConditions ...bool) bool {
	if !isOwner {
		return false
	}
	normalized := normalizeReviewStatus(reviewStatus)
	if normalized != statusPublished && normalized != statusReviewing && normalized != statusRejected && normalized != statusResolved {
		return false
	}
	for _, cond := range extraConditions {
		if !cond {
			return false
		}
	}
	return true
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

func matchMeetupUserRole(item cltypes.MeetupItem, principal auth.Principal, userRole string) bool {
	if userRole == "" {
		return true
	}
	switch userRole {
	case "publisher":
		return principal.Authenticated && item.PublisherUserID == principal.UserID
	case "participant":
		return principal.Authenticated && slices.Contains(item.ParticipantUserIDs, principal.UserID)
	default:
		return true
	}
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

func buildCarpoolPayload(item cltypes.CarpoolItem, canView bool, now time.Time) map[string]any {
	category := normalizedCarpoolCategory(item, now)
	timeText := formatCarpoolTravelText(item.TravelAt, now)
	status := normalizeReviewStatus(item.ReviewStatus)
	return map[string]any{
		"id":                item.ID,
		"category":          category,
		"from":              item.From,
		"to":                item.To,
		"time":              timeText,
		"type":              item.Type,
		"seats_text":        item.SeatsText,
		"price":             item.Price,
		"note":              item.Note,
		"tags":              append([]string(nil), item.Tags...),
		"contact":           visibleValue(canView, item.Contact),
		"status":            status,
		"publisher":         item.Publisher,
		"publisher_initial": item.PublisherInitial,
		"created_at":        item.CreatedAt.Format(time.RFC3339),
		"extra": map[string]any{
			"category":   category,
			"from":       item.From,
			"to":         item.To,
			"time":       timeText,
			"type":       item.Type,
			"seats_text": item.SeatsText,
			"price":      item.Price,
			"note":       item.Note,
			"tags":       append([]string(nil), item.Tags...),
			"contact":    visibleValue(canView, item.Contact),
			"travel_at":  item.TravelAt.In(chinaLocation).Format(time.RFC3339),
			"status":     status,
		},
	}
}

func buildMeetupPayload(item cltypes.MeetupItem, principal auth.Principal, now time.Time) map[string]any {
	canView := canViewContact(principal, item.PublisherUserID)
	userRole := meetupUserRole(item, principal)
	status := mergeMeetupStatus(item.Status, item.ReviewStatus)
	joinedCount := meetupJoinedCount(item)
	remainingSeats := meetupRemainingSeats(item)
	canJoin := principal.Authenticated &&
		userRole == "viewer" &&
		status == statusOpen &&
		remainingSeats > 0 &&
		(item.DeadlineAt.IsZero() || item.DeadlineAt.After(now.UTC())) &&
		(item.StartAt.IsZero() || item.StartAt.After(now.UTC()))

	return map[string]any{
		"id":                 item.ID,
		"category":           item.Category,
		"title":              item.Title,
		"desc":               item.Desc,
		"location":           item.Location,
		"start_at":           item.StartAt.In(chinaLocation).Format(time.RFC3339),
		"deadline_at":        item.DeadlineAt.In(chinaLocation).Format(time.RFC3339),
		"max_participants":   item.MaxParticipants,
		"joined_count":       joinedCount,
		"remaining_seats":    remainingSeats,
		"fee_text":           item.FeeText,
		"tags":               append([]string(nil), item.Tags...),
		"contact":            visibleValue(canView, item.Contact),
		"status":             status,
		"publisher":          item.Publisher,
		"publisher_initial":  item.PublisherInitial,
		"created_at":         item.CreatedAt.Format(time.RFC3339),
		"user_role":          userRole,
		"joined":             userRole == "participant",
		"can_join":           canJoin,
		"can_cancel_join":    userRole == "participant" && status != statusCancelled,
		"extra": map[string]any{
			"category":         item.Category,
			"location":         item.Location,
			"start_at":         item.StartAt.In(chinaLocation).Format(time.RFC3339),
			"deadline_at":      item.DeadlineAt.In(chinaLocation).Format(time.RFC3339),
			"max_participants": item.MaxParticipants,
			"joined_count":     joinedCount,
			"remaining_seats":  remainingSeats,
			"fee_text":         item.FeeText,
			"tags":             append([]string(nil), item.Tags...),
			"contact":          visibleValue(canView, item.Contact),
			"status":           status,
		},
	}
}

func shouldExposeContent(principal auth.Principal, ownerUserID, reviewStatus, userRole string) bool {
	normalized := normalizeReviewStatus(reviewStatus)
	if normalized == statusPublished || normalized == statusResolved {
		return true
	}
	return canAccessPendingContent(principal, ownerUserID, userRole)
}

func ensureContentVisible(principal auth.Principal, ownerUserID, reviewStatus, notFoundMessage string) error {
	if shouldExposeContent(principal, ownerUserID, reviewStatus, "") {
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

func shouldExposeMeetupState(principal auth.Principal, item cltypes.MeetupItem, userRole string) bool {
	if normalizeMeetupStatus(item.Status) != "cancelled" {
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
	return slices.Contains(item.ParticipantUserIDs, principal.UserID)
}

func normalizeMeetupStatus(status string) string {
	switch strings.TrimSpace(strings.ToLower(status)) {
	case statusCancelled:
		return statusCancelled
	case statusFull:
		return statusFull
	default:
		return statusOpen
	}
}

func normalizeReviewStatus(status string) string {
	switch strings.TrimSpace(strings.ToLower(status)) {
	case statusReviewing:
		return statusReviewing
	case statusRejected:
		return statusRejected
	case statusOffline:
		return statusOffline
	case statusResolved:
		return statusResolved
	default:
		return statusPublished
	}
}

func isSupportedReviewStatus(status string) bool {
	switch strings.TrimSpace(strings.ToLower(status)) {
	case statusReviewing, statusPublished, statusRejected, statusOffline:
		return true
	default:
		return false
	}
}

func mergeMeetupStatus(businessStatus, reviewStatus string) string {
	normalizedReviewStatus := normalizeReviewStatus(reviewStatus)
	if normalizedReviewStatus != statusPublished {
		return normalizedReviewStatus
	}
	return normalizeMeetupStatus(businessStatus)
}

func mergeErrandStatus(businessStatus, reviewStatus string) string {
	normalizedReviewStatus := normalizeReviewStatus(reviewStatus)
	if normalizedReviewStatus != statusPublished {
		return normalizedReviewStatus
	}
	switch strings.TrimSpace(strings.ToLower(businessStatus)) {
	case statusCancelled:
		return statusCancelled
	case statusAccepted:
		return statusAccepted
	default:
		return statusPublished
	}
}

func meetupJoinedCount(item cltypes.MeetupItem) int {
	return len(item.ParticipantUserIDs) + 1
}

func meetupRemainingSeats(item cltypes.MeetupItem) int {
	remaining := item.MaxParticipants - meetupJoinedCount(item)
	if remaining < 0 {
		return 0
	}
	return remaining
}

func refreshMeetupStatus(item cltypes.MeetupItem) cltypes.MeetupItem {
	if normalizeMeetupStatus(item.Status) == "cancelled" {
		item.Status = "cancelled"
		return item
	}
	if meetupRemainingSeats(item) == 0 {
		item.Status = "full"
		return item
	}
	item.Status = "open"
	return item
}

func canModerateCampusLife(principal auth.Principal) bool {
	return principal.HasPermission(campusLifeModeratePermission)
}

func matchReviewQuery(query cltypes.ReviewQuery, contentType, reviewStatus string, values ...string) bool {
	if query.ContentType != "" && query.ContentType != contentType {
		return false
	}
	if query.ReviewStatus != "" && normalizeReviewStatus(query.ReviewStatus) != normalizeReviewStatus(reviewStatus) {
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

func normalizedCarpoolCategory(item cltypes.CarpoolItem, now time.Time) string {
	if item.TravelAt.IsZero() {
		if isSupportedCarpoolCategory(item.Category) {
			return item.Category
		}
		return "today"
	}

	travelDate := startOfDay(item.TravelAt.In(chinaLocation))
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

func carpoolTitle(item cltypes.CarpoolItem) string {
	if item.From != "" || item.To != "" {
		return strings.TrimSpace(item.From + " -> " + item.To)
	}
	return defaultCarpoolType(normalizedCarpoolCategory(item, time.Now().In(chinaLocation)))
}

func carpoolFeedDesc(item cltypes.CarpoolItem) string {
	parts := make([]string, 0, 3)
	if text := formatCarpoolTravelText(item.TravelAt, time.Now().In(chinaLocation)); text != "" {
		parts = append(parts, text)
	}
	if item.SeatsText != "" {
		parts = append(parts, item.SeatsText)
	}
	if item.Price != "" {
		parts = append(parts, item.Price)
	}
	if len(parts) > 0 {
		return strings.Join(parts, " · ")
	}
	return item.Note
}

func meetupFeedDesc(item cltypes.MeetupItem) string {
	parts := make([]string, 0, 4)
	if item.Location != "" {
		parts = append(parts, item.Location)
	}
	if !item.StartAt.IsZero() {
		parts = append(parts, item.StartAt.In(chinaLocation).Format("1月2日 15:04"))
	}
	if remaining := meetupRemainingSeats(item); remaining > 0 {
		parts = append(parts, "剩余 "+strconv.Itoa(remaining)+" 位")
	} else {
		parts = append(parts, "人数已满")
	}
	if item.FeeText != "" {
		parts = append(parts, item.FeeText)
	}
	return strings.Join(parts, " · ")
}

func (s *Service) recordAudit(
	ctx context.Context,
	principal auth.Principal,
	action string,
	resourceType string,
	resourceID string,
	message string,
	details map[string]any,
) {
	audit.RecordBestEffort(ctx, s.recorder, audit.Entry{
		ActorID:      principal.UserID,
		ActorName:    displayName(principal),
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Message:      message,
		Details:      details,
	})
}
