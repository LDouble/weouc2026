package repo

import (
	"context"
	"fmt"
	"sync"
	"time"

	cltypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/campus_life/types"
)

type InMemoryRepository struct {
	mu         sync.RWMutex
	markets    map[string]cltypes.MarketItem
	errands    map[string]cltypes.ErrandItem
	resources  map[string]cltypes.ResourceItem
	lostFounds map[string]cltypes.LostFoundItem
	carpools   map[string]cltypes.CarpoolItem
	meetups    map[string]cltypes.MeetupItem
	seq        int64
}

func NewInMemoryRepository() *InMemoryRepository {
	now := time.Now().UTC()
	repository := &InMemoryRepository{
		markets:    make(map[string]cltypes.MarketItem),
		errands:    make(map[string]cltypes.ErrandItem),
		resources:  make(map[string]cltypes.ResourceItem),
		lostFounds: make(map[string]cltypes.LostFoundItem),
		carpools:   make(map[string]cltypes.CarpoolItem),
		meetups:    make(map[string]cltypes.MeetupItem),
		seq:        100,
	}

	repository.markets["market-101"] = cltypes.MarketItem{
		ID:               "market-101",
		Title:            "九成新 iPad Pro 11 寸",
		Desc:             "M2 芯片，日常记笔记和刷题使用，配原装保护壳。",
		ReviewStatus:     "published",
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
		ReviewStatus:     "published",
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
		ReviewStatus:     "published",
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
		ReviewStatus:     "published",
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
		ReviewStatus:     "published",
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
		ReviewStatus:     "published",
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
		ReviewStatus:     "published",
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

	repository.carpools["carpool-101"] = cltypes.CarpoolItem{
		ID:               "carpool-101",
		Category:         "today",
		From:             "海大南门",
		To:               "高铁北站",
		TravelAt:         now.Add(2 * time.Hour),
		Type:             "今日顺路",
		SeatsText:        "余座 2",
		Price:            "人均 18 元",
		Note:             "可放 2 个行李箱，校内拼满即走。",
		Tags:             []string{"今日可拼", "校内出发"},
		Contact:          "wx-carpool-18",
		ReviewStatus:     "published",
		PublisherUserID:  "seed-u5",
		Publisher:        "车主学长",
		PublisherInitial: "车",
		CreatedAt:        now.Add(-40 * time.Minute),
	}
	repository.carpools["carpool-102"] = cltypes.CarpoolItem{
		ID:               "carpool-102",
		Category:         "longterm",
		From:             "学校东门",
		To:               "软件园二期",
		TravelAt:         now.Add(10 * 24 * time.Hour),
		Type:             "长期通勤",
		SeatsText:        "固定 3/4",
		Price:            "月结 AA",
		Note:             "工作日早八晚六，适合长期通勤同学。",
		Tags:             []string{"长期路线", "工作日"},
		Contact:          "站内私信",
		ReviewStatus:     "published",
		PublisherUserID:  "seed-u6",
		Publisher:        "通勤搭子",
		PublisherInitial: "通",
		CreatedAt:        now.Add(-6 * time.Hour),
	}

	repository.meetups["meetup-101"] = cltypes.MeetupItem{
		ID:               "meetup-101",
		Category:         "sports",
		Title:            "今晚东操羽毛球双打局",
		Desc:             "缺 2 位搭子，自带球拍更方便，结束后可一起夜宵。",
		Location:         "东操羽毛球场 3 号场",
		StartAt:          now.Add(5 * time.Hour),
		DeadlineAt:       now.Add(3 * time.Hour),
		MaxParticipants:  4,
		FeeText:          "AA 场地费 15 元/人",
		Tags:             []string{"羽毛球", "今晚", "新手友好"},
		Contact:          "wx-meetup-101",
		Status:           "open",
		ReviewStatus:     "published",
		PublisherUserID:  "seed-u7",
		Publisher:        "运动搭子",
		PublisherInitial: "运",
		ParticipantUserIDs: []string{
			"seed-u8",
		},
		CreatedAt: now.Add(-70 * time.Minute),
	}
	repository.meetups["meetup-102"] = cltypes.MeetupItem{
		ID:               "meetup-102",
		Category:         "study",
		Title:            "高数期末晚自习组队",
		Desc:             "想找 3 位同学一起在图书馆刷题，互相讲题更高效。",
		Location:         "图书馆五楼北区",
		StartAt:          now.Add(28 * time.Hour),
		DeadlineAt:       now.Add(24 * time.Hour),
		MaxParticipants:  5,
		FeeText:          "免费",
		Tags:             []string{"期末复习", "图书馆", "刷题"},
		Contact:          "站内私信",
		Status:           "open",
		ReviewStatus:     "published",
		PublisherUserID:  "seed-u9",
		Publisher:        "数院同学",
		PublisherInitial: "数",
		CreatedAt:        now.Add(-9 * time.Hour),
	}

	return repository
}

func (r *InMemoryRepository) ListMarkets(context.Context) ([]cltypes.MarketItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return cloneMarketSlice(r.markets), nil
}

func (r *InMemoryRepository) GetMarket(_ context.Context, id string) (cltypes.MarketItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, exists := r.markets[id]
	if !exists {
		return cltypes.MarketItem{}, ErrNotFound
	}
	return cloneMarket(item), nil
}

func (r *InMemoryRepository) SaveMarket(_ context.Context, item cltypes.MarketItem) (cltypes.MarketItem, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	item = cloneMarket(item)
	r.markets[item.ID] = item
	return cloneMarket(item), nil
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

func (r *InMemoryRepository) ListErrands(context.Context) ([]cltypes.ErrandItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return cloneErrandSlice(r.errands), nil
}

func (r *InMemoryRepository) GetErrand(_ context.Context, id string) (cltypes.ErrandItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, exists := r.errands[id]
	if !exists {
		return cltypes.ErrandItem{}, ErrNotFound
	}
	return item, nil
}

func (r *InMemoryRepository) SaveErrand(_ context.Context, item cltypes.ErrandItem) (cltypes.ErrandItem, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	cloned := item
	cloned.Images = append([]string(nil), item.Images...)
	r.errands[item.ID] = cloned
	return cloned, nil
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

func (r *InMemoryRepository) ListResources(context.Context) ([]cltypes.ResourceItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return cloneResourceSlice(r.resources), nil
}

func (r *InMemoryRepository) GetResource(_ context.Context, id string) (cltypes.ResourceItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, exists := r.resources[id]
	if !exists {
		return cltypes.ResourceItem{}, ErrNotFound
	}
	return cloneResource(item), nil
}

func (r *InMemoryRepository) SaveResource(_ context.Context, item cltypes.ResourceItem) (cltypes.ResourceItem, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	item = cloneResource(item)
	r.resources[item.ID] = item
	return cloneResource(item), nil
}

func (r *InMemoryRepository) UpdateResource(_ context.Context, id string, mutate func(*cltypes.ResourceItem) error) (cltypes.ResourceItem, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	item, exists := r.resources[id]
	if !exists {
		return cltypes.ResourceItem{}, ErrNotFound
	}
	next := cloneResource(item)
	if err := mutate(&next); err != nil {
		return cltypes.ResourceItem{}, err
	}
	r.resources[id] = next
	return cloneResource(next), nil
}

func (r *InMemoryRepository) ListLostFound(context.Context) ([]cltypes.LostFoundItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return cloneLostFoundSlice(r.lostFounds), nil
}

func (r *InMemoryRepository) GetLostFound(_ context.Context, id string) (cltypes.LostFoundItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, exists := r.lostFounds[id]
	if !exists {
		return cltypes.LostFoundItem{}, ErrNotFound
	}
	return cloneLostFound(item), nil
}

func (r *InMemoryRepository) SaveLostFound(_ context.Context, item cltypes.LostFoundItem) (cltypes.LostFoundItem, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	item = cloneLostFound(item)
	r.lostFounds[item.ID] = item
	return cloneLostFound(item), nil
}

func (r *InMemoryRepository) UpdateLostFound(_ context.Context, id string, mutate func(*cltypes.LostFoundItem) error) (cltypes.LostFoundItem, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	item, exists := r.lostFounds[id]
	if !exists {
		return cltypes.LostFoundItem{}, ErrNotFound
	}
	next := cloneLostFound(item)
	if err := mutate(&next); err != nil {
		return cltypes.LostFoundItem{}, err
	}
	r.lostFounds[id] = next
	return cloneLostFound(next), nil
}

func (r *InMemoryRepository) ListCarpools(context.Context) ([]cltypes.CarpoolItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return cloneCarpoolSlice(r.carpools), nil
}

func (r *InMemoryRepository) GetCarpool(_ context.Context, id string) (cltypes.CarpoolItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, exists := r.carpools[id]
	if !exists {
		return cltypes.CarpoolItem{}, ErrNotFound
	}
	return cloneCarpool(item), nil
}

func (r *InMemoryRepository) SaveCarpool(_ context.Context, item cltypes.CarpoolItem) (cltypes.CarpoolItem, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	item = cloneCarpool(item)
	r.carpools[item.ID] = item
	return cloneCarpool(item), nil
}

func (r *InMemoryRepository) UpdateCarpool(_ context.Context, id string, mutate func(*cltypes.CarpoolItem) error) (cltypes.CarpoolItem, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	item, exists := r.carpools[id]
	if !exists {
		return cltypes.CarpoolItem{}, ErrNotFound
	}
	next := cloneCarpool(item)
	if err := mutate(&next); err != nil {
		return cltypes.CarpoolItem{}, err
	}
	r.carpools[id] = next
	return cloneCarpool(next), nil
}

func (r *InMemoryRepository) ListMeetups(context.Context) ([]cltypes.MeetupItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return cloneMeetupSlice(r.meetups), nil
}

func (r *InMemoryRepository) GetMeetup(_ context.Context, id string) (cltypes.MeetupItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, exists := r.meetups[id]
	if !exists {
		return cltypes.MeetupItem{}, ErrNotFound
	}
	return cloneMeetup(item), nil
}

func (r *InMemoryRepository) SaveMeetup(_ context.Context, item cltypes.MeetupItem) (cltypes.MeetupItem, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	item = cloneMeetup(item)
	r.meetups[item.ID] = item
	return cloneMeetup(item), nil
}

func (r *InMemoryRepository) UpdateMeetup(_ context.Context, id string, mutate func(*cltypes.MeetupItem) error) (cltypes.MeetupItem, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	item, exists := r.meetups[id]
	if !exists {
		return cltypes.MeetupItem{}, ErrNotFound
	}
	next := cloneMeetup(item)
	if err := mutate(&next); err != nil {
		return cltypes.MeetupItem{}, err
	}
	r.meetups[id] = next
	return cloneMeetup(next), nil
}

func (r *InMemoryRepository) NextID(_ context.Context, prefix string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seq++
	return fmt.Sprintf("%s-%d", prefix, r.seq), nil
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
		result = append(result, cloneLostFound(item))
	}
	return result
}

func cloneLostFound(item cltypes.LostFoundItem) cltypes.LostFoundItem {
	return item
}

func cloneCarpoolSlice(items map[string]cltypes.CarpoolItem) []cltypes.CarpoolItem {
	result := make([]cltypes.CarpoolItem, 0, len(items))
	for _, item := range items {
		result = append(result, cloneCarpool(item))
	}
	return result
}

func cloneCarpool(item cltypes.CarpoolItem) cltypes.CarpoolItem {
	cloned := item
	cloned.Tags = append([]string(nil), item.Tags...)
	return cloned
}

func cloneMeetupSlice(items map[string]cltypes.MeetupItem) []cltypes.MeetupItem {
	result := make([]cltypes.MeetupItem, 0, len(items))
	for _, item := range items {
		result = append(result, cloneMeetup(item))
	}
	return result
}

func cloneMeetup(item cltypes.MeetupItem) cltypes.MeetupItem {
	cloned := item
	cloned.Tags = append([]string(nil), item.Tags...)
	cloned.ParticipantUserIDs = append([]string(nil), item.ParticipantUserIDs...)
	return cloned
}
