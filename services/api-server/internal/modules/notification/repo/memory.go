package repo

import (
	"context"
	"fmt"
	"sync"
	"time"

	notificationtypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/notification/types"
)

type InMemoryRepository struct {
	mu       sync.RWMutex
	nextID   int
	messages map[string]notificationtypes.MessageItem
}

func NewInMemoryRepository() *InMemoryRepository {
	repository := &InMemoryRepository{
		nextID:   500,
		messages: make(map[string]notificationtypes.MessageItem),
	}
	repository.seed()
	return repository
}

func (r *InMemoryRepository) ListMessages(context.Context) ([]notificationtypes.MessageItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]notificationtypes.MessageItem, 0, len(r.messages))
	for _, item := range r.messages {
		items = append(items, cloneMessage(item))
	}
	return items, nil
}

func (r *InMemoryRepository) GetMessage(_ context.Context, id string) (notificationtypes.MessageItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, exists := r.messages[id]
	if !exists {
		return notificationtypes.MessageItem{}, ErrNotFound
	}
	return cloneMessage(item), nil
}

func (r *InMemoryRepository) SaveMessage(_ context.Context, item notificationtypes.MessageItem) (notificationtypes.MessageItem, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item = cloneMessage(item)
	r.messages[item.ID] = item
	return cloneMessage(item), nil
}

func (r *InMemoryRepository) UpdateMessage(
	_ context.Context,
	id string,
	mutate func(*notificationtypes.MessageItem) error,
) (notificationtypes.MessageItem, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, exists := r.messages[id]
	if !exists {
		return notificationtypes.MessageItem{}, ErrNotFound
	}
	next := cloneMessage(current)
	if err := mutate(&next); err != nil {
		return notificationtypes.MessageItem{}, err
	}
	r.messages[id] = cloneMessage(next)
	return cloneMessage(next), nil
}

func (r *InMemoryRepository) NextID(_ context.Context, prefix string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.nextID++
	return fmt.Sprintf("%s-%03d", prefix, r.nextID), nil
}

func (r *InMemoryRepository) seed() {
	now := time.Date(2026, 5, 11, 10, 0, 0, 0, time.UTC)
	r.messages["notification-101"] = notificationtypes.MessageItem{
		ID:              "notification-101",
		Title:           "欢迎使用校园综合应用",
		Content:         "你可以通过首页动态快速查看跑腿、组局、二手、资料和失物招领。",
		Category:        "system",
		TargetScope:     "all",
		ActionURL:       "/pages/home/index",
		Publisher:       "系统助手",
		PublisherUserID: "system",
		CreatedAt:       now,
		ReadByUserIDs:   map[string]time.Time{},
	}
	r.messages["notification-102"] = notificationtypes.MessageItem{
		ID:              "notification-102",
		Title:           "教务绑定后可查看联系方式",
		Content:         "若你需要联系发布者，请先完成教务绑定以解锁联系方式查看权限。",
		Category:        "reminder",
		TargetScope:     "users",
		TargetUserIDs:   []string{"u-1"},
		ActionURL:       "/pages/profile/academic/index",
		Publisher:       "校园运营中心",
		PublisherUserID: "admin-001",
		CreatedAt:       now.Add(15 * time.Minute),
		ReadByUserIDs:   map[string]time.Time{},
	}
}

func cloneMessage(item notificationtypes.MessageItem) notificationtypes.MessageItem {
	item.TargetUserIDs = append([]string(nil), item.TargetUserIDs...)
	if item.ReadByUserIDs != nil {
		cloned := make(map[string]time.Time, len(item.ReadByUserIDs))
		for key, value := range item.ReadByUserIDs {
			cloned[key] = value
		}
		item.ReadByUserIDs = cloned
	}
	return item
}
