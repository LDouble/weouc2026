package repo

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	cltypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/campus_life/types"
)

var ErrNotFound = errors.New("item not found")

type Repository interface {
	ListMarkets(ctx context.Context) []cltypes.MarketItem
	GetMarket(ctx context.Context, id string) (cltypes.MarketItem, bool)
	SaveMarket(ctx context.Context, item cltypes.MarketItem) cltypes.MarketItem
	UpdateMarket(ctx context.Context, id string, mutate func(*cltypes.MarketItem) error) (cltypes.MarketItem, error)

	ListErrands(ctx context.Context) []cltypes.ErrandItem
	GetErrand(ctx context.Context, id string) (cltypes.ErrandItem, bool)
	SaveErrand(ctx context.Context, item cltypes.ErrandItem) cltypes.ErrandItem
	UpdateErrand(ctx context.Context, id string, mutate func(*cltypes.ErrandItem) error) (cltypes.ErrandItem, error)

	ListResources(ctx context.Context) []cltypes.ResourceItem
	GetResource(ctx context.Context, id string) (cltypes.ResourceItem, bool)
	SaveResource(ctx context.Context, item cltypes.ResourceItem) cltypes.ResourceItem

	ListLostFound(ctx context.Context) []cltypes.LostFoundItem
	GetLostFound(ctx context.Context, id string) (cltypes.LostFoundItem, bool)
	SaveLostFound(ctx context.Context, item cltypes.LostFoundItem) cltypes.LostFoundItem

	NextID(prefix string) string
}

type InMemoryRepository struct {
	mu         sync.RWMutex
	markets    map[string]cltypes.MarketItem
	errands    map[string]cltypes.ErrandItem
	resources  map[string]cltypes.ResourceItem
	lostFounds map[string]cltypes.LostFoundItem
	seq        int64
}

func NewInMemoryRepository() *InMemoryRepository {
	now := time.Now().UTC()
	repository := &InMemoryRepository{
		markets:    make(map[string]cltypes.MarketItem),
		errands:    make(map[string]cltypes.ErrandItem),
		resources:  make(map[string]cltypes.ResourceItem),
		lostFounds: make(map[string]cltypes.LostFoundItem),
		seq:        100,
	}

	repository.markets["market-101"] = cltypes.MarketItem{
		ID:               "market-101",
		Title:            "九成新 iPad Pro 11 寸",
		Desc:             "M2 芯片，日常记笔记和刷题使用，配原装保护壳。",
		PublisherUserID:  "seed-u1",
		Publisher:        "海大同学",
		PublisherInitial: "海",
		Image:            "https://images.unsplash.com/photo-1544244015-0df4b3ffc6b0?auto=format&fit=crop&w=1200&q=80",
		CreatedAt:        now.Add(-2 * time.Hour),
		Likes:            18,
		LikedByUserIDs:   map[string]bool{},
		Extra: cltypes.MarketExtra{
			Category:      "digital",
			Price:         "4299",
			OriginalPrice: "6499",
			Condition:     "9成新",
			TradeMode:     "校内当面交易",
			Contact:       "wx-hd-ipad",
			Images: []string{
				"https://images.unsplash.com/photo-1544244015-0df4b3ffc6b0?auto=format&fit=crop&w=1200&q=80",
				"https://images.unsplash.com/photo-1511707171634-5f897ff02aa9?auto=format&fit=crop&w=1200&q=80",
			},
		},
	}
	repository.markets["market-102"] = cltypes.MarketItem{
		ID:               "market-102",
		Title:            "求购高数下复习资料",
		Desc:             "想收一本有完整笔记的高数下资料，价格可谈。",
		PublisherUserID:  "seed-u2",
		Publisher:        "理工同学",
		PublisherInitial: "理",
		CreatedAt:        now.Add(-6 * time.Hour),
		Likes:            4,
		LikedByUserIDs:   map[string]bool{},
		Extra: cltypes.MarketExtra{
			Category:  "wanted",
			Price:     "50",
			Condition: "不限",
			TradeMode: "图书馆见面",
			Contact:   "tel-13800000000",
		},
	}

	repository.errands["errand-101"] = cltypes.ErrandItem{
		ID:               "errand-101",
		Title:            "求带一份一食堂鸡腿饭",
		Desc:             "一食堂窗口 3，少辣，送到图书馆东门。",
		Category:         "food",
		RouteStart:       "一食堂",
		RouteEnd:         "图书馆东门",
		Deadline:         now.Add(90 * time.Minute),
		Reward:           "6",
		Contact:          "wx-food-2026",
		Status:           "published",
		PublisherUserID:  "seed-u1",
		Publisher:        "海大同学",
		PublisherInitial: "海",
		CreatedAt:        now.Add(-30 * time.Minute),
	}
	repository.errands["errand-102"] = cltypes.ErrandItem{
		ID:               "errand-102",
		Title:            "西门快递代取",
		Desc:             "菜鸟驿站大件，帮忙带到 8 号宿舍楼。",
		Category:         "parcel",
		RouteStart:       "西门菜鸟驿站",
		RouteEnd:         "8号宿舍楼",
		Deadline:         now.Add(4 * time.Hour),
		Reward:           "5",
		Contact:          "wx-parcel-88",
		Status:           "accepted",
		PublisherUserID:  "seed-u2",
		Publisher:        "理工同学",
		PublisherInitial: "理",
		AcceptorUserID:   "seed-u3",
		CreatedAt:        now.Add(-4 * time.Hour),
	}

	repository.resources["resource-101"] = cltypes.ResourceItem{
		ID:               "resource-101",
		Title:            "离散数学期末冲刺笔记",
		Desc:             "包含重点题型整理、证明题模板和 3 套自测题。",
		PublisherUserID:  "seed-u3",
		Publisher:        "计院学长",
		PublisherInitial: "计",
		CreatedAt:        now.Add(-24 * time.Hour),
		Extra: cltypes.ResourceExtra{
			Category:   "notes",
			CourseName: "离散数学",
			Files: []cltypes.ResourceFile{
				{
					Name:     "离散数学冲刺笔记.pdf",
					URL:      "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
					FileType: "application/pdf",
					FileSize: "2.4MB",
				},
			},
			FileSize: "2.4MB",
			FileType: "application/pdf",
			Likes:    22,
			Views:    153,
		},
	}

	repository.lostFounds["lostfound-101"] = cltypes.LostFoundItem{
		ID:               "lostfound-101",
		Title:            "蓝色校园卡套",
		Desc:             "卡套正面有白色鲸鱼贴纸，内含校园卡。",
		PublisherUserID:  "seed-u1",
		Publisher:        "海大同学",
		PublisherInitial: "海",
		CreatedAt:        now.Add(-3 * time.Hour),
		Extra: cltypes.LostFoundExtra{
			Type:      "lost",
			Category:  "card",
			Location:  "图书馆三楼自习区",
			EventTime: now.Add(-5 * time.Hour).Format(time.RFC3339),
			Contact:   "wx-card-2026",
		},
	}
	repository.lostFounds["lostfound-102"] = cltypes.LostFoundItem{
		ID:               "lostfound-102",
		Title:            "黑色保温杯待认领",
		Desc:             "在二教 302 教室最后一排发现，杯盖有轻微磨损。",
		PublisherUserID:  "seed-u4",
		Publisher:        "校园志愿者",
		PublisherInitial: "校",
		CreatedAt:        now.Add(-26 * time.Hour),
		Extra: cltypes.LostFoundExtra{
			Type:      "found",
			Category:  "daily",
			Location:  "二教 302",
			EventTime: now.Add(-27 * time.Hour).Format(time.RFC3339),
			Contact:   "站内私信",
		},
	}

	return repository
}

func (r *InMemoryRepository) ListMarkets(context.Context) []cltypes.MarketItem {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return cloneMarketSlice(r.markets)
}

func (r *InMemoryRepository) GetMarket(_ context.Context, id string) (cltypes.MarketItem, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, exists := r.markets[id]
	return cloneMarket(item), exists
}

func (r *InMemoryRepository) SaveMarket(_ context.Context, item cltypes.MarketItem) cltypes.MarketItem {
	r.mu.Lock()
	defer r.mu.Unlock()
	item = cloneMarket(item)
	r.markets[item.ID] = item
	return cloneMarket(item)
}

func (r *InMemoryRepository) UpdateMarket(_ context.Context, id string, mutate func(*cltypes.MarketItem) error) (cltypes.MarketItem, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	item, exists := r.markets[id]
	if !exists {
		return cltypes.MarketItem{}, ErrNotFound
	}
	next := cloneMarket(item)
	if err := mutate(&next); err != nil {
		return cltypes.MarketItem{}, err
	}
	r.markets[id] = next
	return cloneMarket(next), nil
}

func (r *InMemoryRepository) ListErrands(context.Context) []cltypes.ErrandItem {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return cloneErrandSlice(r.errands)
}

func (r *InMemoryRepository) GetErrand(_ context.Context, id string) (cltypes.ErrandItem, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, exists := r.errands[id]
	return item, exists
}

func (r *InMemoryRepository) SaveErrand(_ context.Context, item cltypes.ErrandItem) cltypes.ErrandItem {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.errands[item.ID] = item
	return item
}

func (r *InMemoryRepository) UpdateErrand(_ context.Context, id string, mutate func(*cltypes.ErrandItem) error) (cltypes.ErrandItem, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	item, exists := r.errands[id]
	if !exists {
		return cltypes.ErrandItem{}, ErrNotFound
	}
	if err := mutate(&item); err != nil {
		return cltypes.ErrandItem{}, err
	}
	r.errands[id] = item
	return item, nil
}

func (r *InMemoryRepository) ListResources(context.Context) []cltypes.ResourceItem {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return cloneResourceSlice(r.resources)
}

func (r *InMemoryRepository) GetResource(_ context.Context, id string) (cltypes.ResourceItem, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, exists := r.resources[id]
	return cloneResource(item), exists
}

func (r *InMemoryRepository) SaveResource(_ context.Context, item cltypes.ResourceItem) cltypes.ResourceItem {
	r.mu.Lock()
	defer r.mu.Unlock()
	item = cloneResource(item)
	r.resources[item.ID] = item
	return cloneResource(item)
}

func (r *InMemoryRepository) ListLostFound(context.Context) []cltypes.LostFoundItem {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return cloneLostFoundSlice(r.lostFounds)
}

func (r *InMemoryRepository) GetLostFound(_ context.Context, id string) (cltypes.LostFoundItem, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, exists := r.lostFounds[id]
	return item, exists
}

func (r *InMemoryRepository) SaveLostFound(_ context.Context, item cltypes.LostFoundItem) cltypes.LostFoundItem {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lostFounds[item.ID] = item
	return item
}

func (r *InMemoryRepository) NextID(prefix string) string {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seq++
	return fmt.Sprintf("%s-%d", prefix, r.seq)
}

func cloneMarket(item cltypes.MarketItem) cltypes.MarketItem {
	cloned := item
	if item.LikedByUserIDs != nil {
		cloned.LikedByUserIDs = make(map[string]bool, len(item.LikedByUserIDs))
		for key, value := range item.LikedByUserIDs {
			cloned.LikedByUserIDs[key] = value
		}
	}
	cloned.Extra.Images = append([]string(nil), item.Extra.Images...)
	return cloned
}

func cloneMarketSlice(items map[string]cltypes.MarketItem) []cltypes.MarketItem {
	result := make([]cltypes.MarketItem, 0, len(items))
	for _, item := range items {
		result = append(result, cloneMarket(item))
	}
	return result
}

func cloneErrandSlice(items map[string]cltypes.ErrandItem) []cltypes.ErrandItem {
	result := make([]cltypes.ErrandItem, 0, len(items))
	for _, item := range items {
		cloned := item
		cloned.Images = append([]string(nil), item.Images...)
		result = append(result, cloned)
	}
	return result
}

func cloneResource(item cltypes.ResourceItem) cltypes.ResourceItem {
	cloned := item
	cloned.Extra.Files = append([]cltypes.ResourceFile(nil), item.Extra.Files...)
	return cloned
}

func cloneResourceSlice(items map[string]cltypes.ResourceItem) []cltypes.ResourceItem {
	result := make([]cltypes.ResourceItem, 0, len(items))
	for _, item := range items {
		result = append(result, cloneResource(item))
	}
	return result
}

func cloneLostFoundSlice(items map[string]cltypes.LostFoundItem) []cltypes.LostFoundItem {
	result := make([]cltypes.LostFoundItem, 0, len(items))
	for _, item := range items {
		result = append(result, item)
	}
	return result
}
