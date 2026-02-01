package repository

import (
	"context"
	"time"

	"github.com/iots1/vertex-diagram/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoNoteRepository struct {
	Conn *mongo.Collection
}

// NewMongoNoteRepository creates a new note repository
func NewMongoNoteRepository(Conn *mongo.Collection) domain.NoteRepository {
	return &mongoNoteRepository{Conn}
}

func (m *mongoNoteRepository) Store(ctx context.Context, n *domain.Note) error {
	if n.CreatedAt.IsZero() {
		n.CreatedAt = time.Now()
	}
	n.UpdatedAt = time.Now()

	res, err := m.Conn.InsertOne(ctx, n)
	if err == nil {
		if oid, ok := res.InsertedID.(string); ok {
			n.ID = oid
		}
	}
	return err
}

func (m *mongoNoteRepository) StoreMultiple(ctx context.Context, notes []domain.Note) error {
	if len(notes) == 0 {
		return nil
	}

	now := time.Now()
	docs := make([]interface{}, len(notes))
	for i, note := range notes {
		note.CreatedAt = now
		note.UpdatedAt = now
		docs[i] = note
	}

	_, err := m.Conn.InsertMany(ctx, docs)
	return err
}

func (m *mongoNoteRepository) GetByDiagramID(ctx context.Context, diagramID string) ([]domain.Note, error) {
	cursor, err := m.Conn.Find(ctx, bson.M{"diagram_id": diagramID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	notes := make([]domain.Note, 0)
	if err = cursor.All(ctx, &notes); err != nil {
		return nil, err
	}
	return notes, nil
}

func (m *mongoNoteRepository) UpdateByDiagramID(ctx context.Context, diagramID string, notes []domain.Note) error {
	if len(notes) == 0 {
		return nil
	}

	return m.StoreMultiple(ctx, notes)
}

func (m *mongoNoteRepository) DeleteByDiagramID(ctx context.Context, diagramID string) error {
	_, err := m.Conn.DeleteMany(ctx, bson.M{"diagram_id": diagramID})
	return err
}
