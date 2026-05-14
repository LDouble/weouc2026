package service

import (
	"context"
	"time"

	clrepo "github.com/liangluo/weouc2026/services/api-server/internal/modules/campus_life/repo"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/audit"
	"github.com/liangluo/weouc2026/services/api-server/internal/platform/auth"
	"github.com/liangluo/weouc2026/services/api-server/internal/providers/storage_provider"
	"go.mongodb.org/mongo-driver/bson"
)

type Service struct {
	repository      clrepo.Repository
	storageProvider storage_provider.Provider
	recorder        audit.Recorder
}

var chinaLocation = time.FixedZone("Asia/Shanghai", 8*3600)

const campusLifeModeratePermission = "campus_life:moderate"

func New(repository clrepo.Repository, storageProvider storage_provider.Provider, recorder audit.Recorder) *Service {
	return &Service{
		repository:      repository,
		storageProvider: storageProvider,
		recorder:        recorder,
	}
}

func marshalPayload(v any) bson.M {
	raw, err := bson.Marshal(v)
	if err != nil {
		return bson.M{}
	}
	var result bson.M
	if err := bson.Unmarshal(raw, &result); err != nil {
		return bson.M{}
	}
	return result
}

func unmarshalPayload[T any](payload bson.M) (T, error) {
	raw, err := bson.Marshal(payload)
	if err != nil {
		var zero T
		return zero, err
	}
	var result T
	if err := bson.Unmarshal(raw, &result); err != nil {
		var zero T
		return zero, err
	}
	return result, nil
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
