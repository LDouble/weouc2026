package audit

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStore struct {
	logs      *mongo.Collection
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

func NewMongoStore(database *mongo.Database) *MongoStore {
	if database == nil {
		return &MongoStore{}
	}

	return &MongoStore{
		logs:      database.Collection("audit_logs"),
		sequences: database.Collection("audit_sequences"),
	}
}

func (s *MongoStore) Record(ctx context.Context, entry Entry) error {
	if s.logs == nil || s.sequences == nil {
		return fmt.Errorf("mongo audit store collection is nil")
	}

	id, err := s.nextAuditID(ctx, entry.ID)
	if err != nil {
		return err
	}

	entry.ID = id
	entry.Result = firstNonEmpty(entry.Result, "success")
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now().UTC()
	}
	if entry.Details != nil {
		entry.Details = cloneDetails(entry.Details)
	}

	document := mongoEnvelope[Entry]{
		ID:      entry.ID,
		Payload: cloneEntry(entry),
	}
	if _, err := s.logs.InsertOne(ctx, document); err != nil {
		return fmt.Errorf("insert audit log failed: %w", err)
	}

	return nil
}

func (s *MongoStore) List(ctx context.Context) ([]Entry, error) {
	if s.logs == nil {
		return nil, fmt.Errorf("mongo audit log collection is nil")
	}

	cursor, err := s.logs.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("list audit logs failed: %w", err)
	}
	defer cursor.Close(ctx)

	entries := make([]Entry, 0)
	for cursor.Next(ctx) {
		var document mongoEnvelope[Entry]
		if err := cursor.Decode(&document); err != nil {
			return nil, fmt.Errorf("decode audit log failed: %w", err)
		}
		entries = append(entries, cloneEntry(document.Payload))
	}
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("iterate audit logs failed: %w", err)
	}

	return entries, nil
}

func (s *MongoStore) nextAuditID(ctx context.Context, currentID string) (string, error) {
	if currentID != "" {
		return currentID, nil
	}

	var sequence sequenceDocument
	if err := s.sequences.FindOneAndUpdate(
		ctx,
		bson.M{"_id": "audit"},
		bson.M{"$inc": bson.M{"value": 1}},
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	).Decode(&sequence); err != nil {
		return "", fmt.Errorf("increment audit id sequence failed: %w", err)
	}

	return fmt.Sprintf("audit-%03d", sequence.Value), nil
}
