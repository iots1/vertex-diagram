package repository

import (
	"context"
	"time"

	"github.com/iots1/vertex-diagram/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoRepository struct {
	Conn *mongo.Collection
}

// Constructor
func NewMongoRepository(Conn *mongo.Collection) domain.DiagramRepository {
	return &mongoRepository{Conn}
}

func (m *mongoRepository) Fetch(ctx context.Context) ([]domain.Diagram, error) {
	// ดึงข้อมูลไม่เอา Content (เพื่อความเร็ว)
	opts := options.Find()
	cursor, err := m.Conn.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	diagrams := make([]domain.Diagram, 0)
	if err = cursor.All(ctx, &diagrams); err != nil {
		return nil, err
	}
	return diagrams, nil
}

func (m *mongoRepository) GetByID(ctx context.Context, id string) (*domain.Diagram, error) {
	var d domain.Diagram
	// 1. Try string ID
	err := m.Conn.FindOne(ctx, bson.M{"_id": id}).Decode(&d)
	if err == nil {
		return &d, nil
	}

	// 2. Try ObjectID
	if oid, oerr := primitive.ObjectIDFromHex(id); oerr == nil {
		err = m.Conn.FindOne(ctx, bson.M{"_id": oid}).Decode(&d)
		if err == nil {
			return &d, nil
		}
	}

	return nil, err
}

func (m *mongoRepository) Store(ctx context.Context, d *domain.Diagram) error {
	if d.CreatedAt.IsZero() {
		d.CreatedAt = time.Now()
	}
	d.UpdatedAt = time.Now()

	if d.ID == "" {
		res, err := m.Conn.InsertOne(ctx, d)
		if err == nil {
			if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
				d.ID = oid.Hex()
			} else if sid, ok := res.InsertedID.(string); ok {
				d.ID = sid
			}
		}
		return err
	}

	// For Upsert, we must handle both ID types in filter
	filter := bson.M{"_id": d.ID}
	if oid, oerr := primitive.ObjectIDFromHex(d.ID); oerr == nil {
		// Use ObjectID if possible for better compatibility with server-side generated IDs
		filter = bson.M{"_id": oid}
	}

	// Use UpdateOne with $set to avoid replacing _id field
	update := bson.M{
		"$set": bson.M{
			"name":       d.Name,
			"content":    d.Content,
			"created_at": d.CreatedAt,
			"updated_at": d.UpdatedAt,
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := m.Conn.UpdateOne(ctx, filter, update, opts)
	return err
}

func (m *mongoRepository) Update(ctx context.Context, d *domain.Diagram) error {
	d.UpdatedAt = time.Now()

	filter := bson.M{"_id": d.ID}
	if oid, oerr := primitive.ObjectIDFromHex(d.ID); oerr == nil {
		filter = bson.M{"_id": oid}
	}

	// Use UpdateOne with $set to avoid replacing _id field
	update := bson.M{
		"$set": bson.M{
			"name":       d.Name,
			"content":    d.Content,
			"created_at": d.CreatedAt,
			"updated_at": d.UpdatedAt,
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := m.Conn.UpdateOne(ctx, filter, update, opts)
	return err
}

func (m *mongoRepository) Delete(ctx context.Context, id string) error {
	// Try string ID
	res, err := m.Conn.DeleteOne(ctx, bson.M{"_id": id})
	if err == nil && res.DeletedCount > 0 {
		return nil
	}

	// Try ObjectID
	if oid, oerr := primitive.ObjectIDFromHex(id); oerr == nil {
		_, err = m.Conn.DeleteOne(ctx, bson.M{"_id": oid})
		return err
	}
	
	return err
}