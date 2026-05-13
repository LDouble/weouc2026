package audit

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

type Entry struct {
	ID           string
	ActorID      string
	ActorName    string
	Action       string
	ResourceType string
	ResourceID   string
	Result       string
	Message      string
	Details      map[string]any
	CreatedAt    time.Time
}

type ListQuery struct {
	ActorID      string
	Action       string
	ResourceType string
	ResourceID   string
}

type Recorder interface {
	Record(ctx context.Context, entry Entry) error
}

type Repository interface {
	List(ctx context.Context, query ListQuery) ([]Entry, error)
}

type Store interface {
	Recorder
	Repository
}

type InMemoryStore struct {
	mu      sync.RWMutex
	nextID  int
	entries []Entry
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		nextID:  900,
		entries: make([]Entry, 0),
	}
}

func (s *InMemoryStore) Record(_ context.Context, entry Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.nextID++
	entry.ID = firstNonEmpty(entry.ID, fmt.Sprintf("audit-%03d", s.nextID))
	entry.Result = firstNonEmpty(entry.Result, "success")
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now().UTC()
	}
	if entry.Details != nil {
		entry.Details = cloneDetails(entry.Details)
	}
	s.entries = append(s.entries, entry)
	return nil
}

func (s *InMemoryStore) List(_ context.Context, query ListQuery) ([]Entry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]Entry, 0, len(s.entries))
	for _, entry := range s.entries {
		if !matchEntryQuery(entry, query) {
			continue
		}
		result = append(result, cloneEntry(entry))
	}
	return result, nil
}

func RecordBestEffort(ctx context.Context, recorder Recorder, entry Entry) {
	if recorder == nil {
		return
	}
	_ = recorder.Record(ctx, entry)
}

func cloneEntry(entry Entry) Entry {
	entry.Details = cloneDetails(entry.Details)
	return entry
}

func cloneDetails(details map[string]any) map[string]any {
	if details == nil {
		return nil
	}
	cloned := make(map[string]any, len(details))
	for key, value := range details {
		cloned[key] = value
	}
	return cloned
}

func matchEntryQuery(entry Entry, query ListQuery) bool {
	if strings.TrimSpace(query.ActorID) != "" && entry.ActorID != strings.TrimSpace(query.ActorID) {
		return false
	}
	if strings.TrimSpace(query.Action) != "" && entry.Action != strings.TrimSpace(query.Action) {
		return false
	}
	if strings.TrimSpace(query.ResourceType) != "" && entry.ResourceType != strings.TrimSpace(query.ResourceType) {
		return false
	}
	if strings.TrimSpace(query.ResourceID) != "" && entry.ResourceID != strings.TrimSpace(query.ResourceID) {
		return false
	}
	return true
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
