package repo

import (
	"context"
	"fmt"
	"sync"
	"time"

	portaltypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/portal/types"
)

type InMemoryRepository struct {
	mu      sync.RWMutex
	nextID  int
	banners map[string]portaltypes.BannerItem
	notices map[string]portaltypes.NoticeItem
}

func NewInMemoryRepository() *InMemoryRepository {
	repository := &InMemoryRepository{
		nextID:  300,
		banners: make(map[string]portaltypes.BannerItem),
		notices: make(map[string]portaltypes.NoticeItem),
	}
	repository.seed()
	return repository
}

func (r *InMemoryRepository) ListBanners(context.Context) ([]portaltypes.BannerItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]portaltypes.BannerItem, 0, len(r.banners))
	for _, item := range r.banners {
		items = append(items, item)
	}
	return items, nil
}

func (r *InMemoryRepository) GetBanner(_ context.Context, id string) (portaltypes.BannerItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, exists := r.banners[id]
	if !exists {
		return portaltypes.BannerItem{}, ErrNotFound
	}
	return item, nil
}

func (r *InMemoryRepository) SaveBanner(_ context.Context, item portaltypes.BannerItem) (portaltypes.BannerItem, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.banners[item.ID] = item
	return item, nil
}

func (r *InMemoryRepository) DeleteBanner(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.banners[id]; !exists {
		return ErrNotFound
	}
	delete(r.banners, id)
	return nil
}

func (r *InMemoryRepository) ListNotices(context.Context) ([]portaltypes.NoticeItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]portaltypes.NoticeItem, 0, len(r.notices))
	for _, item := range r.notices {
		items = append(items, cloneNotice(item))
	}
	return items, nil
}

func (r *InMemoryRepository) GetNotice(_ context.Context, id string) (portaltypes.NoticeItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, exists := r.notices[id]
	if !exists {
		return portaltypes.NoticeItem{}, ErrNotFound
	}
	return cloneNotice(item), nil
}

func (r *InMemoryRepository) SaveNotice(_ context.Context, item portaltypes.NoticeItem) (portaltypes.NoticeItem, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item = cloneNotice(item)
	r.notices[item.ID] = item
	return cloneNotice(item), nil
}

func (r *InMemoryRepository) DeleteNotice(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.notices[id]; !exists {
		return ErrNotFound
	}
	delete(r.notices, id)
	return nil
}

func (r *InMemoryRepository) NextID(_ context.Context, prefix string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.nextID++
	return fmt.Sprintf("%s-%03d", prefix, r.nextID), nil
}

func (r *InMemoryRepository) seed() {
	now := time.Date(2026, 5, 11, 9, 0, 0, 0, time.UTC)
	r.banners["banner-101"] = portaltypes.BannerItem{
		ID:          "banner-101",
		Title:       "新学期校园生活入口升级",
		Description: "跑腿、组局、二手、资料与失物招领统一接入新首页。",
		ImageURL:    "https://example.com/portal/banner-101.png",
		ActionURL:   "/pages/home/index",
		Sort:        1,
		CreatedAt:   now,
	}
	r.banners["banner-102"] = portaltypes.BannerItem{
		ID:          "banner-102",
		Title:       "教务绑定后可查看联系方式",
		Description: "涉及联系方式的内容均以后端绑定状态裁剪结果为准。",
		ImageURL:    "https://example.com/portal/banner-102.png",
		ActionURL:   "/pages/profile/academic/index",
		Sort:        2,
		CreatedAt:   now.Add(10 * time.Minute),
	}

	r.notices["notice-101"] = portaltypes.NoticeItem{
		ID:              "notice-101",
		Title:           "校园综合应用内测启动",
		Summary:         "微信小程序已开放跑腿、组局、二手、资料与失物招领基础能力。",
		Content:         "本周开放第一轮内测，欢迎同学体验并通过站内消息反馈问题。",
		Audience:        "all",
		Tags:            []string{"内测", "公告"},
		Pinned:          true,
		PublisherUserID: "admin-001",
		Publisher:       "校园运营中心",
		PublishedAt:     now,
		CreatedAt:       now,
	}
	r.notices["notice-102"] = portaltypes.NoticeItem{
		ID:              "notice-102",
		Title:           "二手与跑腿发布规范",
		Summary:         "新发布内容默认进入审核，违规内容会被驳回或下线。",
		Content:         "请勿发布违法违规、交易风险高或联系方式异常的信息，审核员会依据规则处理。",
		Audience:        "student",
		Tags:            []string{"审核", "发布规范"},
		Pinned:          false,
		PublisherUserID: "admin-001",
		Publisher:       "校园运营中心",
		PublishedAt:     now.Add(30 * time.Minute),
		CreatedAt:       now.Add(30 * time.Minute),
	}
}

func cloneNotice(item portaltypes.NoticeItem) portaltypes.NoticeItem {
	item.Tags = append([]string(nil), item.Tags...)
	return item
}
