package repo

import (
	"context"
	"fmt"

	portaltypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/portal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	banners   *mongo.Collection
	notices   *mongo.Collection
	sequences *mongo.Collection
}

type mongoEnvelope[T any] struct {
	ID      string `bson:"_id"`
	Payload T      `bson:"payload"`
}

type sequenceDocument struct {
	ID    string `bson:"_id"`
	Value int64  `bson:"value"`
}

func NewMongoRepository(database *mongo.Database) *MongoRepository {
	if database == nil {
		return &MongoRepository{}
	}

	return &MongoRepository{
		banners:   database.Collection("portal_banners"),
		notices:   database.Collection("portal_notices"),
		sequences: database.Collection("portal_sequences"),
	}
}

func (r *MongoRepository) ListBanners(ctx context.Context) ([]portaltypes.BannerItem, error) {
	return listDocuments[portaltypes.BannerItem](ctx, r.banners, "list portal banners failed")
}

func (r *MongoRepository) GetBanner(ctx context.Context, id string) (portaltypes.BannerItem, error) {
	return getDocument[portaltypes.BannerItem](ctx, r.banners, id, "get portal banner failed")
}

func (r *MongoRepository) SaveBanner(ctx context.Context, item portaltypes.BannerItem) (portaltypes.BannerItem, error) {
	return saveDocument(ctx, r.banners, item.ID, item, "save portal banner failed")
}

func (r *MongoRepository) DeleteBanner(ctx context.Context, id string) error {
	return deleteDocument(ctx, r.banners, id, "delete portal banner failed")
}

func (r *MongoRepository) ListNotices(ctx context.Context) ([]portaltypes.NoticeItem, error) {
	return listDocuments[portaltypes.NoticeItem](ctx, r.notices, "list portal notices failed")
}

func (r *MongoRepository) GetNotice(ctx context.Context, id string) (portaltypes.NoticeItem, error) {
	return getDocument[portaltypes.NoticeItem](ctx, r.notices, id, "get portal notice failed")
}

func (r *MongoRepository) SaveNotice(ctx context.Context, item portaltypes.NoticeItem) (portaltypes.NoticeItem, error) {
	return saveDocument(ctx, r.notices, item.ID, item, "save portal notice failed")
}

func (r *MongoRepository) DeleteNotice(ctx context.Context, id string) error {
	return deleteDocument(ctx, r.notices, id, "delete portal notice failed")
}

func (r *MongoRepository) NextID(ctx context.Context, prefix string) (string, error) {
	if r.sequences == nil {
		return "", fmt.Errorf("mongo portal sequence collection is nil")
	}

	var sequence sequenceDocument
	if err := r.sequences.FindOneAndUpdate(
		ctx,
		bson.M{"_id": prefix},
		bson.M{"$inc": bson.M{"value": 1}},
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	).Decode(&sequence); err != nil {
		return "", fmt.Errorf("increment portal sequence failed: %w", err)
	}

	return fmt.Sprintf("%s-%03d", prefix, sequence.Value), nil
}

func listDocuments[T any](ctx context.Context, collection *mongo.Collection, message string) ([]T, error) {
	if collection == nil {
		return nil, fmt.Errorf("mongo collection is nil")
	}

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", message, err)
	}
	defer cursor.Close(ctx)

	items := make([]T, 0)
	for cursor.Next(ctx) {
		var document mongoEnvelope[T]
		if err := cursor.Decode(&document); err != nil {
			return nil, fmt.Errorf("%s: %w", message, err)
		}
		items = append(items, document.Payload)
	}
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", message, err)
	}

	return items, nil
}

func getDocument[T any](ctx context.Context, collection *mongo.Collection, id, message string) (T, error) {
	var zero T
	if collection == nil {
		return zero, fmt.Errorf("mongo collection is nil")
	}

	var document mongoEnvelope[T]
	if err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&document); err != nil {
		if err == mongo.ErrNoDocuments {
			return zero, ErrNotFound
		}
		return zero, fmt.Errorf("%s: %w", message, err)
	}

	return document.Payload, nil
}

func saveDocument[T any](ctx context.Context, collection *mongo.Collection, id string, payload T, message string) (T, error) {
	if collection == nil {
		return payload, fmt.Errorf("mongo collection is nil")
	}

	document := mongoEnvelope[T]{
		ID:      id,
		Payload: payload,
	}
	if _, err := collection.ReplaceOne(ctx, bson.M{"_id": id}, document, options.Replace().SetUpsert(true)); err != nil {
		return payload, fmt.Errorf("%s: %w", message, err)
	}

	return payload, nil
}

func deleteDocument(ctx context.Context, collection *mongo.Collection, id, message string) error {
	if collection == nil {
		return fmt.Errorf("mongo collection is nil")
	}

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("%s: %w", message, err)
	}
	if result.DeletedCount == 0 {
		return ErrNotFound
	}

	return nil
}
