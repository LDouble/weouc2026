package repo

import (
	"context"
	"fmt"
	"time"

	cltypes "github.com/liangluo/weouc2026/services/api-server/internal/modules/campus_life/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(database *mongo.Database) *MongoRepository {
	if database == nil {
		return &MongoRepository{}
	}
	coll := database.Collection("community_content")
	return &MongoRepository{collection: coll}
}

func (r *MongoRepository) EnsureIndexes(ctx context.Context) error {
	if r.collection == nil {
		return nil
	}
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "content_type", Value: 1}, {Key: "status", Value: 1}, {Key: "created_at", Value: -1}}},
		{Keys: bson.D{{Key: "status", Value: 1}, {Key: "created_at", Value: -1}}},
		{Keys: bson.D{{Key: "publisher_user_id", Value: 1}, {Key: "created_at", Value: -1}}},
		{Keys: bson.D{{Key: "publisher_user_id", Value: 1}, {Key: "status", Value: 1}, {Key: "created_at", Value: -1}}},
		{Keys: bson.D{{Key: "deleted_at", Value: 1}}, Options: options.Index().SetSparse(true)},
	}
	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}

func (r *MongoRepository) Save(ctx context.Context, item cltypes.CommunityContent) (cltypes.CommunityContent, error) {
	if r.collection == nil {
		return item, fmt.Errorf("mongo collection is nil")
	}
	now := time.Now().UTC()
	item.CreatedAt = now
	item.UpdatedAt = now
	if item.ID.IsZero() {
		doc, err := bson.Marshal(item)
		if err != nil {
			return item, fmt.Errorf("marshal community_content failed: %w", err)
		}
		var m bson.M
		if err := bson.Unmarshal(doc, &m); err != nil {
			return item, fmt.Errorf("unmarshal community_content failed: %w", err)
		}
		delete(m, "_id")
		res, err := r.collection.InsertOne(ctx, m)
		if err != nil {
			return item, fmt.Errorf("insert community_content failed: %w", err)
		}
		item.ID = res.InsertedID.(primitive.ObjectID)
		return item, nil
	}
	_, err := r.collection.InsertOne(ctx, item)
	if err != nil {
		return item, fmt.Errorf("insert community_content failed: %w", err)
	}
	return item, nil
}

func (r *MongoRepository) GetByID(ctx context.Context, id string) (cltypes.CommunityContent, error) {
	var zero cltypes.CommunityContent
	if r.collection == nil {
		return zero, fmt.Errorf("mongo collection is nil")
	}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return zero, ErrNotFound
	}
	filter := bson.M{"_id": objID, "deleted_at": nil}
	var result cltypes.CommunityContent
	if err := r.collection.FindOne(ctx, filter).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return zero, ErrNotFound
		}
		return zero, fmt.Errorf("get community_content by id failed: %w", err)
	}
	return result, nil
}

func (r *MongoRepository) Update(ctx context.Context, id string, mutate func(*cltypes.CommunityContent) error) (cltypes.CommunityContent, error) {
	current, err := r.GetByID(ctx, id)
	if err != nil {
		return current, err
	}
	if err := mutate(&current); err != nil {
		return current, err
	}
	current.UpdatedAt = time.Now().UTC()
	_, err = r.collection.ReplaceOne(ctx, bson.M{"_id": current.ID}, current)
	if err != nil {
		return current, fmt.Errorf("update community_content failed: %w", err)
	}
	return current, nil
}

func (r *MongoRepository) ListByType(ctx context.Context, contentType string, filter cltypes.ContentFilter) ([]cltypes.CommunityContent, int64, error) {
	if r.collection == nil {
		return nil, 0, fmt.Errorf("mongo collection is nil")
	}
	query := bson.M{"content_type": contentType, "deleted_at": nil}
	if len(filter.Statuses) > 0 {
		query["status"] = bson.M{"$in": filter.Statuses}
	}
	if filter.Category != "" {
		query["type_payload.category"] = filter.Category
	}
	if filter.SubType != "" {
		query["type_payload.type"] = filter.SubType
	}
	if filter.Keyword != "" {
		query["$or"] = []bson.M{
			{"title": bson.M{"$regex": filter.Keyword, "$options": "i"}},
			{"desc": bson.M{"$regex": filter.Keyword, "$options": "i"}},
		}
	}
	if !filter.IncludeAllStatus {
		visibleStatuses := filter.VisibleStatuses
		if len(visibleStatuses) == 0 {
			visibleStatuses = cltypes.VisibleStatuses
		}
		if filter.CurrentUserID != "" {
			orClauses := []bson.M{
				{"publisher_user_id": filter.CurrentUserID},
				{"status": bson.M{"$in": visibleStatuses}},
			}
			if filter.AcceptorUserID != "" {
				orClauses = append(orClauses, bson.M{"type_payload.acceptor_user_id": filter.AcceptorUserID})
			}
			if filter.ParticipantUserID != "" {
				orClauses = append(orClauses, bson.M{"type_payload.participant_user_ids": filter.ParticipantUserID})
			}
			if existingOr, ok := query["$or"]; ok {
				query["$and"] = []bson.M{
					{"$or": existingOr},
					{"$or": orClauses},
				}
				delete(query, "$or")
			} else {
				query["$or"] = orClauses
			}
		} else {
			query["status"] = bson.M{"$in": visibleStatuses}
			if filter.AcceptorUserID != "" {
				query["type_payload.acceptor_user_id"] = filter.AcceptorUserID
			}
			if filter.ParticipantUserID != "" {
				query["type_payload.participant_user_ids"] = filter.ParticipantUserID
			}
		}
	}
	if filter.PublisherUserID != "" {
		query["publisher_user_id"] = filter.PublisherUserID
	}
	total, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("count community_content by type failed: %w", err)
	}
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}})
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	opts.SetSkip(int64((page - 1) * pageSize))
	opts.SetLimit(int64(pageSize))

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("list community_content by type failed: %w", err)
	}
	defer cursor.Close(ctx)

	items := make([]cltypes.CommunityContent, 0)
	for cursor.Next(ctx) {
		var doc cltypes.CommunityContent
		if err := cursor.Decode(&doc); err != nil {
			return nil, 0, fmt.Errorf("decode community_content failed: %w", err)
		}
		items = append(items, doc)
	}
	if err := cursor.Err(); err != nil {
		return nil, 0, fmt.Errorf("cursor error: %w", err)
	}
	return items, total, nil
}

func (r *MongoRepository) ListForFeed(ctx context.Context, filter cltypes.FeedFilter) ([]cltypes.CommunityContent, int64, error) {
	if r.collection == nil {
		return nil, 0, fmt.Errorf("mongo collection is nil")
	}
	query := bson.M{"deleted_at": nil}
	if len(filter.FeedTypes) > 0 {
		query["content_type"] = bson.M{"$in": filter.FeedTypes}
	}
	if filter.Keyword != "" {
		query["$or"] = []bson.M{
			{"title": bson.M{"$regex": filter.Keyword, "$options": "i"}},
			{"desc": bson.M{"$regex": filter.Keyword, "$options": "i"}},
		}
	}
	if !filter.IncludeAllStatus {
		visibleStatuses := filter.VisibleStatuses
		if len(visibleStatuses) == 0 {
			visibleStatuses = cltypes.VisibleStatuses
		}
		if filter.CurrentUserID != "" {
			orClauses := []bson.M{
				{"publisher_user_id": filter.CurrentUserID},
				{"status": bson.M{"$in": visibleStatuses}},
			}
			if filter.AcceptorUserID != "" {
				orClauses = append(orClauses, bson.M{"type_payload.acceptor_user_id": filter.AcceptorUserID})
			}
			if filter.ParticipantUserID != "" {
				orClauses = append(orClauses, bson.M{"type_payload.participant_user_ids": filter.ParticipantUserID})
			}
			if existingOr, ok := query["$or"]; ok {
				query["$and"] = []bson.M{
					{"$or": existingOr},
					{"$or": orClauses},
				}
				delete(query, "$or")
			} else {
				query["$or"] = orClauses
			}
		} else {
			query["status"] = bson.M{"$in": visibleStatuses}
			if filter.AcceptorUserID != "" {
				query["type_payload.acceptor_user_id"] = filter.AcceptorUserID
			}
			if filter.ParticipantUserID != "" {
				query["type_payload.participant_user_ids"] = filter.ParticipantUserID
			}
		}
	}
	if filter.PublisherUserID != "" {
		query["publisher_user_id"] = filter.PublisherUserID
	}
	total, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("count community_content for feed failed: %w", err)
	}
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}})
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	opts.SetSkip(int64((page - 1) * pageSize))
	opts.SetLimit(int64(pageSize))

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("list community_content for feed failed: %w", err)
	}
	defer cursor.Close(ctx)

	items := make([]cltypes.CommunityContent, 0)
	for cursor.Next(ctx) {
		var doc cltypes.CommunityContent
		if err := cursor.Decode(&doc); err != nil {
			return nil, 0, fmt.Errorf("decode community_content failed: %w", err)
		}
		items = append(items, doc)
	}
	if err := cursor.Err(); err != nil {
		return nil, 0, fmt.Errorf("cursor error: %w", err)
	}
	return items, total, nil
}
