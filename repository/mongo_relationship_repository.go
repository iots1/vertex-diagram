package repository

import (
	"context"
	"time"

	"github.com/iots1/vertex-diagram/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRelationshipRepository struct {
	Conn *mongo.Collection
}

// NewMongoRelationshipRepository creates a new relationship repository
func NewMongoRelationshipRepository(Conn *mongo.Collection) domain.RelationshipRepository {
	return &mongoRelationshipRepository{Conn}
}

func (m *mongoRelationshipRepository) Store(ctx context.Context, r *domain.Relationship) error {
	if r.CreatedAt.IsZero() {
		r.CreatedAt = time.Now()
	}
	r.UpdatedAt = time.Now()

	res, err := m.Conn.InsertOne(ctx, r)
	if err == nil {
		if oid, ok := res.InsertedID.(string); ok {
			r.ID = oid
		}
	}
	return err
}

func (m *mongoRelationshipRepository) StoreMultiple(ctx context.Context, relationships []domain.Relationship) error {
	if len(relationships) == 0 {
		return nil
	}

	now := time.Now()
	docs := make([]interface{}, len(relationships))
	for i, rel := range relationships {
		rel.CreatedAt = now
		rel.UpdatedAt = now
		docs[i] = rel
	}

	_, err := m.Conn.InsertMany(ctx, docs)
	return err
}

func (m *mongoRelationshipRepository) GetByDiagramID(ctx context.Context, diagramID string) ([]domain.Relationship, error) {
	cursor, err := m.Conn.Find(ctx, bson.M{"diagram_id": diagramID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	relationships := make([]domain.Relationship, 0)
	if err = cursor.All(ctx, &relationships); err != nil {
		return nil, err
	}
	return relationships, nil
}

func (m *mongoRelationshipRepository) UpdateByDiagramID(ctx context.Context, diagramID string, relationships []domain.Relationship) error {
	if len(relationships) == 0 {
		return nil
	}

	return m.StoreMultiple(ctx, relationships)
}

func (m *mongoRelationshipRepository) DeleteByDiagramID(ctx context.Context, diagramID string) error {
	_, err := m.Conn.DeleteMany(ctx, bson.M{"diagram_id": diagramID})
	return err
}
