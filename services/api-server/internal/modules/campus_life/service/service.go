package service

import (
	"context"
	"path"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	clrepo "github.com/liangluo/weouc2026/services/api-server/internal/modules/campus_life/repo"
	cltypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/campus_life/types"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/auth"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/httpx"
	"github.com/liangluo/weouc2026/services/api-server/internal/providers/storage_provider"
)

type Service struct {
	repository      clrepo.Repository
	storageProvider storage_provider.Provider
}

func New(repository clrepo.Repository, storageProvider storage_provider.Provider) *Service {
	return &Service{
		repository:      repository,
		storageProvider: storageProvider,
	}
}

func (s *Service) ListFeed(ctx context.Context, principal auth.Principal, query cltypes.FeedQuery) map[string]any {
	type row struct {
		id        string
		feedType  string
		createdAt time.Time
		payload   map[string]any
	}

	rows := make([]row, 0)
	for _, item := range s.repository.ListMarkets(ctx) {
		if !matchUserRole(item.PublisherUserID, "", principal, query.UserRole) {
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
				"image":           resolvedImage,
				"extra": map[string]any{
					"images":   resolvedImages,
					"likes":    item.Likes,
					"comments": 0,
				},
			},
		})
	}
	for _, item := range s.repository.ListErrands(ctx) {
		role := errandUserRole(item, principal)
		if !matchUserRole(item.PublisherUserID, item.AcceptorUserID, principal, query.UserRole) {
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
				"extra": map[string]any{
					"likes":     0,
					"comments":  0,
					"user_role": role,
				},
			},
		})
	}
	for _, item := range s.repository.ListResources(ctx) {
		if !matchUserRole(item.PublisherUserID, "", principal, query.UserRole) {
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
			},
		})
	}
	for _, item := range s.repository.ListLostFound(ctx) {
		if !matchUserRole(item.PublisherUserID, "", principal, query.UserRole) {
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

	return listEnvelope(list, len(rows), query.Pagination)
}

func (s *Service) ListMarket(ctx context.Context, principal auth.Principal, query cltypes.MarketQuery) map[string]any {
	items := s.repository.ListMarkets(ctx)
	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt.After(items[j].CreatedAt)
	})

	filtered := make([]map[string]any, 0)
	for _, item := range items {
		if query.Category != "" && query.Category != item.Extra.Category {
			continue
		}
		if !matchKeyword(query.Keyword, item.Title, item.Desc) {
			continue
		}
		canViewContact := canViewContact(principal, item.PublisherUserID)
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

	return listEnvelope(paginateMaps(filtered, query.Pagination), len(filtered), query.Pagination)
}

func (s *Service) GetMarketDetail(ctx context.Context, principal auth.Principal, id string) (map[string]any, error) {
	item, exists := s.repository.GetMarket(ctx, id)
	if !exists {
		return nil, httpx.NotFound("二手商品不存在", nil)
	}

	canView := canViewContact(principal, item.PublisherUserID)
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
		"can_view_contact":  canView,
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
	item := cltypes.MarketItem{
		ID:               s.repository.NextID("market"),
		Title:            strings.TrimSpace(request.Title),
		Desc:             strings.TrimSpace(request.Desc),
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
	item = s.repository.SaveMarket(ctx, item)

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
		return err
	}

	return nil
}

func (s *Service) ListErrands(ctx context.Context, principal auth.Principal, query cltypes.ErrandQuery) map[string]any {
	items := s.repository.ListErrands(ctx)
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.After(items[j].CreatedAt) })
	filtered := make([]map[string]any, 0)
	for _, item := range items {
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
		accepted := item.Status == "accepted"
		filtered = append(filtered, map[string]any{
			"id":                item.ID,
			"category":          item.Category,
			"title":             item.Title,
			"desc":              item.Desc,
			"route_start":       item.RouteStart,
			"route_end":         item.RouteEnd,
			"deadline":          item.Deadline.Format(time.RFC3339),
			"reward":            item.Reward,
			"status":            item.Status,
			"user_role":         role,
			"is_accepted":       accepted,
			"views":             0,
			"publisher":         item.Publisher,
			"publisher_initial": item.PublisherInitial,
			"created_at":        item.CreatedAt.Format(time.RFC3339),
		})
	}

	return listEnvelope(paginateMaps(filtered, query.Pagination), len(filtered), query.Pagination)
}

func (s *Service) GetErrandDetail(ctx context.Context, principal auth.Principal, id string) (map[string]any, error) {
	item, exists := s.repository.GetErrand(ctx, id)
	if !exists {
		return nil, httpx.NotFound("跑腿任务不存在", nil)
	}
	role := errandUserRole(item, principal)
	canView := canViewContact(principal, item.PublisherUserID)
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
			"status":            item.Status,
			"is_accepted":       item.Status == "accepted",
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
				"status":      item.Status,
				"urgent":      item.Urgent,
			},
		},
		"user_role":        role,
		"can_view_contact": canView,
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
	item := cltypes.ErrandItem{
		ID:               s.repository.NextID("errand"),
		Title:            strings.TrimSpace(request.Title),
		Desc:             strings.TrimSpace(request.Desc),
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
	item = s.repository.SaveErrand(ctx, item)
	return map[string]any{"id": item.ID}, nil
}

func (s *Service) AcceptErrand(ctx context.Context, principal auth.Principal, taskID string) error {
	_, err := s.repository.UpdateErrand(ctx, taskID, func(item *cltypes.ErrandItem) error {
		if item.PublisherUserID == principal.UserID {
			return httpx.BadRequest("不能接自己发布的任务", nil)
		}
		if item.Status == "accepted" {
			return httpx.BadRequest("该任务已被接单", nil)
		}
		if item.Status == "cancelled" {
			return httpx.BadRequest("该任务已取消", nil)
		}
		item.Status = "accepted"
		item.AcceptorUserID = principal.UserID
		return nil
	})
	if err != nil {
		if err == clrepo.ErrNotFound {
			return httpx.NotFound("跑腿任务不存在", nil)
		}
		return err
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
		return err
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
		return err
	}

	return nil
}

func (s *Service) ListResources(ctx context.Context, query cltypes.ResourceQuery) map[string]any {
	items := s.repository.ListResources(ctx)
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.After(items[j].CreatedAt) })
	filtered := make([]map[string]any, 0)
	for _, item := range items {
		if query.Category != "" && query.Category != item.Extra.Category {
			continue
		}
		if !matchKeyword(query.Keyword, item.Title, item.Desc) {
			continue
		}
		resolvedFiles := resolveResourceFiles(ctx, s.storageProvider, item.Extra.Files)
		filtered = append(filtered, map[string]any{
			"id":                item.ID,
			"title":             item.Title,
			"desc":              item.Desc,
			"publisher":         item.Publisher,
			"publisher_initial": item.PublisherInitial,
			"created_at":        item.CreatedAt.Format(time.RFC3339),
			"extra": map[string]any{
				"category":     item.Extra.Category,
				"course_name":  item.Extra.CourseName,
				"contact":      item.Extra.Contact,
				"files":        resolvedFiles,
				"file_size":    item.Extra.FileSize,
				"file_type":    item.Extra.FileType,
				"download_url": firstResourceURL(resolvedFiles, resolveManagedURL(ctx, s.storageProvider, item.Extra.DownloadURL)),
				"likes":        item.Extra.Likes,
				"views":        item.Extra.Views,
			},
		})
	}
	return listEnvelope(paginateMaps(filtered, query.Pagination), len(filtered), query.Pagination)
}

func (s *Service) GetResourceDetail(ctx context.Context, id string) (map[string]any, error) {
	item, exists := s.repository.GetResource(ctx, id)
	if !exists {
		return nil, httpx.NotFound("资料不存在", nil)
	}
	resolvedFiles := resolveResourceFiles(ctx, s.storageProvider, item.Extra.Files)
	return map[string]any{
		"id":                item.ID,
		"title":             item.Title,
		"desc":              item.Desc,
		"publisher":         item.Publisher,
		"publisher_initial": item.PublisherInitial,
		"created_at":        item.CreatedAt.Format(time.RFC3339),
		"extra": map[string]any{
			"category":     item.Extra.Category,
			"course_name":  item.Extra.CourseName,
			"contact":      item.Extra.Contact,
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
	item := cltypes.ResourceItem{
		ID:               s.repository.NextID("resource"),
		Title:            strings.TrimSpace(request.Title),
		Desc:             firstNonEmpty(strings.TrimSpace(request.Desc), strings.TrimSpace(request.Title)),
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
	item = s.repository.SaveResource(ctx, item)
	return map[string]any{"id": item.ID}, nil
}

func (s *Service) ListLostFound(ctx context.Context, principal auth.Principal, query cltypes.LostFoundQuery) map[string]any {
	items := s.repository.ListLostFound(ctx)
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.After(items[j].CreatedAt) })
	filtered := make([]map[string]any, 0)
	for _, item := range items {
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
		filtered = append(filtered, map[string]any{
			"id":                item.ID,
			"title":             item.Title,
			"desc":              item.Desc,
			"publisher":         item.Publisher,
			"publisher_initial": item.PublisherInitial,
			"created_at":        item.CreatedAt.Format(time.RFC3339),
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
	return listEnvelope(paginateMaps(filtered, query.Pagination), len(filtered), query.Pagination)
}

func (s *Service) GetLostFoundDetail(ctx context.Context, principal auth.Principal, id string) (map[string]any, error) {
	item, exists := s.repository.GetLostFound(ctx, id)
	if !exists {
		return nil, httpx.NotFound("失物招领不存在", nil)
	}
	canView := canViewContact(principal, item.PublisherUserID)
	return map[string]any{
		"id":                item.ID,
		"title":             item.Title,
		"desc":              item.Desc,
		"publisher":         item.Publisher,
		"publisher_initial": item.PublisherInitial,
		"created_at":        item.CreatedAt.Format(time.RFC3339),
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
	item := cltypes.LostFoundItem{
		ID:               s.repository.NextID("lostfound"),
		Title:            strings.TrimSpace(request.Title),
		Desc:             strings.TrimSpace(request.Desc),
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
	item = s.repository.SaveLostFound(ctx, item)
	return map[string]any{"id": item.ID}, nil
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
