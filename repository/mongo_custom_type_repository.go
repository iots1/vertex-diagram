package repository

import (
	"context"
	"time"

	"github.com/iots1/vertex-diagram/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoCustomTypeRepository struct {
	Conn *mongo.Collection
}

// NewMongoCustomTypeRepository creates a new custom type repository
func NewMongoCustomTypeRepository(Conn *mongo.Collection) domain.CustomTypeRepository {
	return &mongoCustomTypeRepository{Conn}
}

func (m *mongoCustomTypeRepository) Store(ctx context.Context, ct *domain.CustomType) error {
	if ct.CreatedAt.IsZero() {
		ct.CreatedAt = time.Now()
	}
	ct.UpdatedAt = time.Now()

	res, err := m.Conn.InsertOne(ctx, ct)
	if err == nil {
		if oid, ok := res.InsertedID.(string); ok {
			ct.ID = oid
		}
	}
	return err
}

func (m *mongoCustomTypeRepository) StoreMultiple(ctx context.Context, customTypes []domain.CustomType) error {
	if len(customTypes) == 0 {
		return nil
	}

	now := time.Now()
	docs := make([]interface{}, len(customTypes))
	for i, ct := range customTypes {
		ct.CreatedAt = now
		ct.UpdatedAt = now
		docs[i] = ct
	}

	_, err := m.Conn.InsertMany(ctx, docs)
	return err
}

func (m *mongoCustomTypeRepository) GetByDiagramID(ctx context.Context, diagramID string) ([]domain.CustomType, error) {
	cursor, err := m.Conn.Find(ctx, bson.M{"diagram_id": diagramID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	customTypes := make([]domain.CustomType, 0)
	if err = cursor.All(ctx, &customTypes); err != nil {
		return nil, err
	}
	return customTypes, nil
}

func (m *mongoCustomTypeRepository) UpdateByDiagramID(ctx context.Context, diagramID string, customTypes []domain.CustomType) error {
	if len(customTypes) == 0 {
		return nil
	}

	return m.StoreMultiple(ctx, customTypes)
}

func (m *mongoCustomTypeRepository) DeleteByDiagramID(ctx context.Context, diagramID string) error {
	_, err := m.Conn.DeleteMany(ctx, bson.M{"diagram_id": diagramID})
	return err
}
