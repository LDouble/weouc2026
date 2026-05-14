package repo

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type legacyEnvelope struct {
	ID      string `bson:"_id"`
	Payload bson.M `bson:"payload"`
}

type legacyMapping struct {
	Collection  string
	ContentType string
}

var legacyMappings = []legacyMapping{
	{Collection: "campus_life_markets", ContentType: "market"},
	{Collection: "campus_life_errands", ContentType: "errand"},
	{Collection: "campus_life_resources", ContentType: "resource"},
	{Collection: "campus_life_lost_found", ContentType: "lost_found"},
	{Collection: "campus_life_carpools", ContentType: "carpool"},
	{Collection: "campus_life_meetups", ContentType: "meetup"},
}

func MigrateFromLegacyCollections(ctx context.Context, database *mongo.Database) error {
	targetColl := database.Collection("community_content")
	var totalMigrated int
	var totalFailed int

	for _, mapping := range legacyMappings {
		sourceColl := database.Collection(mapping.Collection)
		cursor, err := sourceColl.Find(ctx, bson.M{})
		if err != nil {
			log.Printf("migrate: failed to open cursor for %s: %v", mapping.Collection, err)
			totalFailed++
			continue
		}

		migrated := 0
		failed := 0

		for cursor.Next(ctx) {
			var envelope legacyEnvelope
			if err := cursor.Decode(&envelope); err != nil {
				log.Printf("migrate: failed to decode document from %s: %v", mapping.Collection, err)
				failed++
				continue
			}

			doc, err := buildNewDocument(envelope.Payload, mapping.ContentType)
			if err != nil {
				log.Printf("migrate: failed to build document from %s (id=%s): %v", mapping.Collection, envelope.ID, err)
				failed++
				continue
			}

			if _, err := targetColl.InsertOne(ctx, doc); err != nil {
				log.Printf("migrate: failed to insert document from %s (id=%s): %v", mapping.Collection, envelope.ID, err)
				failed++
				continue
			}
			migrated++
		}

		if err := cursor.Err(); err != nil {
			log.Printf("migrate: cursor error for %s: %v", mapping.Collection, err)
		}
		cursor.Close(ctx)

		log.Printf("migrate: %s → %s: migrated=%d, failed=%d", mapping.Collection, mapping.ContentType, migrated, failed)
		totalMigrated += migrated
		totalFailed += failed
	}

	log.Printf("migrate: total migrated=%d, total failed=%d", totalMigrated, totalFailed)

	if totalMigrated == 0 && totalFailed > 0 {
		return fmt.Errorf("migration completed with 0 successes and %d failures", totalFailed)
	}
	return nil
}

func buildNewDocument(payload bson.M, contentType string) (bson.M, error) {
	doc := make(bson.M)

	delete(payload, "_id")

	doc["content_type"] = contentType

	reviewStatus := ""
	if rs, ok := payload["review_status"]; ok {
		reviewStatus = fmt.Sprintf("%v", rs)
	}

	originalStatus := ""
	if s, ok := payload["status"]; ok {
		originalStatus = fmt.Sprintf("%v", s)
	}

	switch contentType {
	case "errand", "meetup":
		if reviewStatus == "published" {
			doc["status"] = originalStatus
		} else {
			doc["status"] = reviewStatus
		}
	default:
		doc["status"] = reviewStatus
	}

	if contentType == "lost_found" && doc["status"] == "resolved" {
		doc["status"] = "resolved"
	}

	delete(payload, "review_status")
	delete(payload, "status")

	if puid, ok := payload["publisher_user_id"]; ok {
		doc["publisher_user_id"] = puid
		puidStr := fmt.Sprintf("%v", puid)
		doc["created_by"] = puidStr
		doc["updated_by"] = puidStr
	}

	delete(payload, "publisher")
	delete(payload, "publisher_initial")

	extractContact(payload, doc)
	extractImages(payload, doc)
	extractTags(payload, doc)

	if ca, ok := payload["created_at"]; ok {
		doc["created_at"] = ca
		doc["updated_at"] = ca
	} else {
		delete(payload, "created_at")
		doc["created_at"] = nil
		doc["updated_at"] = nil
	}
	delete(payload, "created_at")

	doc["deleted_at"] = nil

	typePayload := buildTypePayload(payload, contentType)
	doc["type_payload"] = typePayload

	if contentType == "market" {
		if lb, ok := payload["liked_by_user_ids"]; ok {
			doc["liked_by_user_ids"] = lb
		}
		if l, ok := payload["likes"]; ok {
			doc["likes"] = l
		}
	}
	delete(payload, "liked_by_user_ids")
	delete(payload, "likes")

	delete(payload, "extra")
	delete(payload, "image")

	for k, v := range payload {
		doc[k] = v
	}

	return doc, nil
}

func extractContact(payload bson.M, doc bson.M) {
	if c, ok := payload["contact"]; ok {
		doc["contact"] = c
		delete(payload, "contact")
		return
	}
	if extra, ok := payload["extra"].(bson.M); ok {
		if c, ok := extra["contact"]; ok {
			doc["contact"] = c
			delete(extra, "contact")
		}
	}
}

func extractImages(payload bson.M, doc bson.M) {
	if img, ok := payload["images"]; ok {
		doc["images"] = img
		delete(payload, "images")
		return
	}
	if extra, ok := payload["extra"].(bson.M); ok {
		if img, ok := extra["images"]; ok {
			doc["images"] = img
			delete(extra, "images")
		}
	}
}

func extractTags(payload bson.M, doc bson.M) {
	if t, ok := payload["tags"]; ok {
		doc["tags"] = t
		delete(payload, "tags")
		return
	}
	if extra, ok := payload["extra"].(bson.M); ok {
		if t, ok := extra["tags"]; ok {
			doc["tags"] = t
			delete(extra, "tags")
		}
	}
}

func buildTypePayload(payload bson.M, contentType string) bson.M {
	extra := bson.M{}
	if e, ok := payload["extra"].(bson.M); ok {
		for k, v := range e {
			extra[k] = v
		}
	}

	tp := bson.M{}

	switch contentType {
	case "market":
		extractTo(tp, extra, "category")
		extractTo(tp, extra, "price")
		extractTo(tp, extra, "original_price")
		extractTo(tp, extra, "condition")
		extractTo(tp, extra, "trade_mode")
	case "errand":
		extractTo(tp, payload, "category")
		extractTo(tp, payload, "route_start")
		extractTo(tp, payload, "route_end")
		extractTo(tp, payload, "deadline")
		extractTo(tp, payload, "reward")
		extractTo(tp, payload, "urgent")
		extractTo(tp, payload, "acceptor_user_id")
	case "resource":
		extractTo(tp, extra, "category")
		extractTo(tp, extra, "course_name")
		extractTo(tp, extra, "files")
		extractTo(tp, extra, "file_size")
		extractTo(tp, extra, "file_type")
		extractTo(tp, extra, "download_url")
		extractTo(tp, extra, "likes")
		extractTo(tp, extra, "views")
	case "lost_found":
		extractTo(tp, extra, "type")
		extractTo(tp, extra, "category")
		extractTo(tp, extra, "location")
		extractTo(tp, extra, "event_time")
		extractTo(tp, extra, "item_feature")
	case "carpool":
		extractTo(tp, extra, "category")
		extractTo(tp, extra, "from")
		extractTo(tp, extra, "to")
		extractTo(tp, extra, "travel_at")
		extractTo(tp, extra, "type")
		extractTo(tp, extra, "seats_text")
		extractTo(tp, extra, "price")
		extractTo(tp, extra, "note")
	case "meetup":
		extractTo(tp, extra, "category")
		extractTo(tp, extra, "location")
		extractTo(tp, extra, "start_at")
		extractTo(tp, extra, "deadline_at")
		extractTo(tp, extra, "max_participants")
		extractTo(tp, extra, "fee_text")
		extractTo(tp, extra, "participant_user_ids")
	}

	return tp
}

func extractTo(target bson.M, source bson.M, key string) {
	if v, ok := source[key]; ok {
		target[key] = v
		delete(source, key)
	}
}
