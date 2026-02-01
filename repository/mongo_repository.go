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
	opts := options.Find().SetProjection(bson.M{"content": 0})
	cursor, err := m.Conn.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	var diagrams []domain.Diagram
	if err = cursor.All(ctx, &diagrams); err != nil {
		return nil, err
	}
	return diagrams, nil
}

func (m *mongoRepository) GetByID(ctx context.Context, id string) (*domain.Diagram, error) {
	oid, _ := primitive.ObjectIDFromHex(id)
	var d domain.Diagram
	err := m.Conn.FindOne(ctx, bson.M{"_id": oid}).Decode(&d)
	return &d, err
}

func (m *mongoRepository) Store(ctx context.Context, d *domain.Diagram) error {
	d.CreatedAt = time.Now()
	d.UpdatedAt = time.Now()
	res, err := m.Conn.InsertOne(ctx, d)
	if err == nil {
		d.ID = res.InsertedID.(primitive.ObjectID)
	}
	return err
}

func (m *mongoRepository) Update(ctx context.Context, d *domain.Diagram) error {
	d.UpdatedAt = time.Now()
	filter := bson.M{"_id": d.ID}
	update := bson.M{"$set": d}
	_, err := m.Conn.UpdateOne(ctx, filter, update)
	return err
}