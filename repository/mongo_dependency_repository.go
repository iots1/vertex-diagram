package repository

import (
	"context"
	"time"

	"github.com/iots1/vertex-diagram/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoDependencyRepository struct {
	Conn *mongo.Collection
}

// NewMongoDependencyRepository creates a new dependency repository
func NewMongoDependencyRepository(Conn *mongo.Collection) domain.DependencyRepository {
	return &mongoDependencyRepository{Conn}
}

func (m *mongoDependencyRepository) Store(ctx context.Context, d *domain.Dependency) error {
	if d.CreatedAt.IsZero() {
		d.CreatedAt = time.Now()
	}
	d.UpdatedAt = time.Now()

	res, err := m.Conn.InsertOne(ctx, d)
	if err == nil {
		if oid, ok := res.InsertedID.(string); ok {
			d.ID = oid
		}
	}
	return err
}

func (m *mongoDependencyRepository) StoreMultiple(ctx context.Context, dependencies []domain.Dependency) error {
	if len(dependencies) == 0 {
		return nil
	}

	now := time.Now()
	docs := make([]interface{}, len(dependencies))
	for i, dep := range dependencies {
		dep.CreatedAt = now
		dep.UpdatedAt = now
		docs[i] = dep
	}

	_, err := m.Conn.InsertMany(ctx, docs)
	return err
}

func (m *mongoDependencyRepository) GetByDiagramID(ctx context.Context, diagramID string) ([]domain.Dependency, error) {
	cursor, err := m.Conn.Find(ctx, bson.M{"diagram_id": diagramID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	dependencies := make([]domain.Dependency, 0)
	if err = cursor.All(ctx, &dependencies); err != nil {
		return nil, err
	}
	return dependencies, nil
}

func (m *mongoDependencyRepository) UpdateByDiagramID(ctx context.Context, diagramID string, dependencies []domain.Dependency) error {
	if len(dependencies) == 0 {
		return nil
	}

	return m.StoreMultiple(ctx, dependencies)
}

func (m *mongoDependencyRepository) DeleteByDiagramID(ctx context.Context, diagramID string) error {
	_, err := m.Conn.DeleteMany(ctx, bson.M{"diagram_id": diagramID})
	return err
}
