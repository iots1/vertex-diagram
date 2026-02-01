package repository

import (
	"context"
	"time"

	"github.com/iots1/vertex-diagram/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoTableRepository struct {
	Conn *mongo.Collection
}

// NewMongoTableRepository creates a new table repository
func NewMongoTableRepository(Conn *mongo.Collection) domain.TableRepository {
	return &mongoTableRepository{Conn}
}

func (m *mongoTableRepository) Store(ctx context.Context, t *domain.Table) error {
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now()
	}
	t.UpdatedAt = time.Now()

	res, err := m.Conn.InsertOne(ctx, t)
	if err == nil {
		if oid, ok := res.InsertedID.(string); ok {
			t.ID = oid
		}
	}
	return err
}

func (m *mongoTableRepository) StoreMultiple(ctx context.Context, tables []domain.Table) error {
	if len(tables) == 0 {
		return nil
	}

	now := time.Now()
	docs := make([]interface{}, len(tables))
	for i, table := range tables {
		table.CreatedAt = now
		table.UpdatedAt = now
		docs[i] = table
	}

	_, err := m.Conn.InsertMany(ctx, docs)
	return err
}

func (m *mongoTableRepository) GetByDiagramID(ctx context.Context, diagramID string) ([]domain.Table, error) {
	cursor, err := m.Conn.Find(ctx, bson.M{"diagram_id": diagramID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	tables := make([]domain.Table, 0)
	if err = cursor.All(ctx, &tables); err != nil {
		return nil, err
	}
	return tables, nil
}

func (m *mongoTableRepository) UpdateByDiagramID(ctx context.Context, diagramID string, tables []domain.Table) error {
	if len(tables) == 0 {
		return nil
	}

	return m.StoreMultiple(ctx, tables)
}

func (m *mongoTableRepository) DeleteByDiagramID(ctx context.Context, diagramID string) error {
	_, err := m.Conn.DeleteMany(ctx, bson.M{"diagram_id": diagramID})
	return err
}
