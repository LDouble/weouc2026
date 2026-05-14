package audit

import (
	"context"
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

type Recorder interface {
	Record(ctx context.Context, entry Entry) error
}

type Repository interface {
	List(ctx context.Context) ([]Entry, error)
}

type Store interface {
	Recorder
	Repository
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

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
