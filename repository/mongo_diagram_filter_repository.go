package repository

import (
	"context"
	"time"

	"github.com/iots1/vertex-diagram/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoDiagramFilterRepository struct {
	Conn *mongo.Collection
}

// NewMongoDiagramFilterRepository creates a new diagram filter repository
func NewMongoDiagramFilterRepository(Conn *mongo.Collection) domain.DiagramFilterRepository {
	return &mongoDiagramFilterRepository{Conn}
}

func (m *mongoDiagramFilterRepository) Store(ctx context.Context, df *domain.DiagramFilter) error {
	if df.CreatedAt.IsZero() {
		df.CreatedAt = time.Now()
	}
	df.UpdatedAt = time.Now()

	// First, try to find existing filter for this diagram
	var existing domain.DiagramFilter
	err := m.Conn.FindOne(ctx, bson.M{"diagram_id": df.DiagramID}).Decode(&existing)

	if err != nil && err != mongo.ErrNoDocuments {
		// Real error occurred
		return err
	}

	if err == mongo.ErrNoDocuments {
		// No existing filter, create new one
		res, err := m.Conn.InsertOne(ctx, df)
		if err == nil && res.InsertedID != nil {
			if oid, ok := res.InsertedID.(string); ok {
				df.ID = oid
			}
		}
		return err
	}

	// Existing filter found, update it
	df.ID = existing.ID // Preserve the original ID
	_, err = m.Conn.ReplaceOne(ctx, bson.M{"_id": existing.ID}, df)
	return err
}

func (m *mongoDiagramFilterRepository) GetByDiagramID(ctx context.Context, diagramID string) (*domain.DiagramFilter, error) {
	var filter domain.DiagramFilter
	err := m.Conn.FindOne(ctx, bson.M{"diagram_id": diagramID}).Decode(&filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &filter, nil
}

func (m *mongoDiagramFilterRepository) DeleteByDiagramID(ctx context.Context, diagramID string) error {
	_, err := m.Conn.DeleteOne(ctx, bson.M{"diagram_id": diagramID})
	return err
}
