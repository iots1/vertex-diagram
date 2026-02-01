package repository

import (
	"context"
	"time"

	"github.com/iots1/vertex-diagram/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoAreaRepository struct {
	Conn *mongo.Collection
}

// NewMongoAreaRepository creates a new area repository
func NewMongoAreaRepository(Conn *mongo.Collection) domain.AreaRepository {
	return &mongoAreaRepository{Conn}
}

func (m *mongoAreaRepository) Store(ctx context.Context, a *domain.Area) error {
	if a.CreatedAt.IsZero() {
		a.CreatedAt = time.Now()
	}
	a.UpdatedAt = time.Now()

	res, err := m.Conn.InsertOne(ctx, a)
	if err == nil {
		if oid, ok := res.InsertedID.(string); ok {
			a.ID = oid
		}
	}
	return err
}

func (m *mongoAreaRepository) StoreMultiple(ctx context.Context, areas []domain.Area) error {
	if len(areas) == 0 {
		return nil
	}

	now := time.Now()
	docs := make([]interface{}, len(areas))
	for i, area := range areas {
		area.CreatedAt = now
		area.UpdatedAt = now
		docs[i] = area
	}

	_, err := m.Conn.InsertMany(ctx, docs)
	return err
}

func (m *mongoAreaRepository) GetByDiagramID(ctx context.Context, diagramID string) ([]domain.Area, error) {
	cursor, err := m.Conn.Find(ctx, bson.M{"diagram_id": diagramID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	areas := make([]domain.Area, 0)
	if err = cursor.All(ctx, &areas); err != nil {
		return nil, err
	}
	return areas, nil
}

func (m *mongoAreaRepository) UpdateByDiagramID(ctx context.Context, diagramID string, areas []domain.Area) error {
	if len(areas) == 0 {
		return nil
	}

	return m.StoreMultiple(ctx, areas)
}

func (m *mongoAreaRepository) DeleteByDiagramID(ctx context.Context, diagramID string) error {
	_, err := m.Conn.DeleteMany(ctx, bson.M{"diagram_id": diagramID})
	return err
}
