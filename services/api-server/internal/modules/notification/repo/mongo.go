package repo

import (
	"context"
	"fmt"

	notificationtypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/notification/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	messages  *mongo.Collection
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
		messages:  database.Collection("notification_messages"),
		sequences: database.Collection("notification_sequences"),
	}
}

func (r *MongoRepository) ListMessages(ctx context.Context) ([]notificationtypes.MessageItem, error) {
	return listDocuments[notificationtypes.MessageItem](ctx, r.messages, "list notification messages failed")
}

func (r *MongoRepository) GetMessage(ctx context.Context, id string) (notificationtypes.MessageItem, error) {
	return getDocument[notificationtypes.MessageItem](ctx, r.messages, id, "get notification message failed")
}

func (r *MongoRepository) SaveMessage(
	ctx context.Context,
	item notificationtypes.MessageItem,
) (notificationtypes.MessageItem, error) {
	return saveDocument(ctx, r.messages, item.ID, item, "save notification message failed")
}

func (r *MongoRepository) UpdateMessage(
	ctx context.Context,
	id string,
	mutate func(*notificationtypes.MessageItem) error,
) (notificationtypes.MessageItem, error) {
	current, err := r.GetMessage(ctx, id)
	if err != nil {
		return notificationtypes.MessageItem{}, err
	}
	if err := mutate(&current); err != nil {
		return notificationtypes.MessageItem{}, err
	}

	return r.SaveMessage(ctx, current)
}

func (r *MongoRepository) NextID(ctx context.Context, prefix string) (string, error) {
	if r.sequences == nil {
		return "", fmt.Errorf("mongo notification sequence collection is nil")
	}

	var sequence sequenceDocument
	if err := r.sequences.FindOneAndUpdate(
		ctx,
		bson.M{"_id": prefix},
		bson.M{"$inc": bson.M{"value": 1}},
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	).Decode(&sequence); err != nil {
		return "", fmt.Errorf("increment notification sequence failed: %w", err)
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
