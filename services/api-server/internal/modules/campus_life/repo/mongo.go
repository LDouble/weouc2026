package repo

import (
	"context"
	"fmt"

	cltypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/campus_life/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	markets   *mongo.Collection
	errands   *mongo.Collection
	resources *mongo.Collection
	lostFound *mongo.Collection
	carpools  *mongo.Collection
	meetups   *mongo.Collection
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
		markets:   database.Collection("campus_life_markets"),
		errands:   database.Collection("campus_life_errands"),
		resources: database.Collection("campus_life_resources"),
		lostFound: database.Collection("campus_life_lost_found"),
		carpools:  database.Collection("campus_life_carpools"),
		meetups:   database.Collection("campus_life_meetups"),
		sequences: database.Collection("campus_life_sequences"),
	}
}

func (r *MongoRepository) ListMarkets(ctx context.Context) ([]cltypes.MarketItem, error) {
	return listDocuments[cltypes.MarketItem](ctx, r.markets, "list campus_life markets failed")
}

func (r *MongoRepository) GetMarket(ctx context.Context, id string) (cltypes.MarketItem, error) {
	return getDocument[cltypes.MarketItem](ctx, r.markets, id, "get campus_life market failed")
}

func (r *MongoRepository) SaveMarket(ctx context.Context, item cltypes.MarketItem) (cltypes.MarketItem, error) {
	return saveDocument(ctx, r.markets, item.ID, item, "save campus_life market failed")
}

func (r *MongoRepository) UpdateMarket(ctx context.Context, id string, mutate func(*cltypes.MarketItem) error) (cltypes.MarketItem, error) {
	return updateDocument(ctx, r.markets, id, mutate, "update campus_life market failed")
}

func (r *MongoRepository) ListErrands(ctx context.Context) ([]cltypes.ErrandItem, error) {
	return listDocuments[cltypes.ErrandItem](ctx, r.errands, "list campus_life errands failed")
}

func (r *MongoRepository) GetErrand(ctx context.Context, id string) (cltypes.ErrandItem, error) {
	return getDocument[cltypes.ErrandItem](ctx, r.errands, id, "get campus_life errand failed")
}

func (r *MongoRepository) SaveErrand(ctx context.Context, item cltypes.ErrandItem) (cltypes.ErrandItem, error) {
	return saveDocument(ctx, r.errands, item.ID, item, "save campus_life errand failed")
}

func (r *MongoRepository) UpdateErrand(ctx context.Context, id string, mutate func(*cltypes.ErrandItem) error) (cltypes.ErrandItem, error) {
	return updateDocument(ctx, r.errands, id, mutate, "update campus_life errand failed")
}

func (r *MongoRepository) ListResources(ctx context.Context) ([]cltypes.ResourceItem, error) {
	return listDocuments[cltypes.ResourceItem](ctx, r.resources, "list campus_life resources failed")
}

func (r *MongoRepository) GetResource(ctx context.Context, id string) (cltypes.ResourceItem, error) {
	return getDocument[cltypes.ResourceItem](ctx, r.resources, id, "get campus_life resource failed")
}

func (r *MongoRepository) SaveResource(ctx context.Context, item cltypes.ResourceItem) (cltypes.ResourceItem, error) {
	return saveDocument(ctx, r.resources, item.ID, item, "save campus_life resource failed")
}

func (r *MongoRepository) UpdateResource(ctx context.Context, id string, mutate func(*cltypes.ResourceItem) error) (cltypes.ResourceItem, error) {
	return updateDocument(ctx, r.resources, id, mutate, "update campus_life resource failed")
}

func (r *MongoRepository) ListLostFound(ctx context.Context) ([]cltypes.LostFoundItem, error) {
	return listDocuments[cltypes.LostFoundItem](ctx, r.lostFound, "list campus_life lost_found failed")
}

func (r *MongoRepository) GetLostFound(ctx context.Context, id string) (cltypes.LostFoundItem, error) {
	return getDocument[cltypes.LostFoundItem](ctx, r.lostFound, id, "get campus_life lost_found failed")
}

func (r *MongoRepository) SaveLostFound(ctx context.Context, item cltypes.LostFoundItem) (cltypes.LostFoundItem, error) {
	return saveDocument(ctx, r.lostFound, item.ID, item, "save campus_life lost_found failed")
}

func (r *MongoRepository) UpdateLostFound(ctx context.Context, id string, mutate func(*cltypes.LostFoundItem) error) (cltypes.LostFoundItem, error) {
	return updateDocument(ctx, r.lostFound, id, mutate, "update campus_life lost_found failed")
}

func (r *MongoRepository) ListCarpools(ctx context.Context) ([]cltypes.CarpoolItem, error) {
	return listDocuments[cltypes.CarpoolItem](ctx, r.carpools, "list campus_life carpools failed")
}

func (r *MongoRepository) GetCarpool(ctx context.Context, id string) (cltypes.CarpoolItem, error) {
	return getDocument[cltypes.CarpoolItem](ctx, r.carpools, id, "get campus_life carpool failed")
}

func (r *MongoRepository) SaveCarpool(ctx context.Context, item cltypes.CarpoolItem) (cltypes.CarpoolItem, error) {
	return saveDocument(ctx, r.carpools, item.ID, item, "save campus_life carpool failed")
}

func (r *MongoRepository) UpdateCarpool(ctx context.Context, id string, mutate func(*cltypes.CarpoolItem) error) (cltypes.CarpoolItem, error) {
	return updateDocument(ctx, r.carpools, id, mutate, "update campus_life carpool failed")
}

func (r *MongoRepository) ListMeetups(ctx context.Context) ([]cltypes.MeetupItem, error) {
	return listDocuments[cltypes.MeetupItem](ctx, r.meetups, "list campus_life meetups failed")
}

func (r *MongoRepository) GetMeetup(ctx context.Context, id string) (cltypes.MeetupItem, error) {
	return getDocument[cltypes.MeetupItem](ctx, r.meetups, id, "get campus_life meetup failed")
}

func (r *MongoRepository) SaveMeetup(ctx context.Context, item cltypes.MeetupItem) (cltypes.MeetupItem, error) {
	return saveDocument(ctx, r.meetups, item.ID, item, "save campus_life meetup failed")
}

func (r *MongoRepository) UpdateMeetup(ctx context.Context, id string, mutate func(*cltypes.MeetupItem) error) (cltypes.MeetupItem, error) {
	return updateDocument(ctx, r.meetups, id, mutate, "update campus_life meetup failed")
}

func (r *MongoRepository) NextID(ctx context.Context, prefix string) (string, error) {
	if r.sequences == nil {
		return "", fmt.Errorf("mongo campus_life sequence collection is nil")
	}

	var sequence sequenceDocument
	if err := r.sequences.FindOneAndUpdate(
		ctx,
		bson.M{"_id": prefix},
		bson.M{"$inc": bson.M{"value": 1}},
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	).Decode(&sequence); err != nil {
		return "", fmt.Errorf("increment campus_life sequence failed: %w", err)
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

func updateDocument[T any](
	ctx context.Context,
	collection *mongo.Collection,
	id string,
	mutate func(*T) error,
	message string,
) (T, error) {
	current, err := getDocument[T](ctx, collection, id, message)
	if err != nil {
		var zero T
		return zero, err
	}
	if err := mutate(&current); err != nil {
		var zero T
		return zero, err
	}

	return saveDocument(ctx, collection, id, current, message)
}
