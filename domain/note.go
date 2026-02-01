package domain

import (
	"context"
	"time"
)

// Note represents a text note annotation on the diagram canvas
type Note struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	DiagramID string    `bson:"diagram_id" json:"diagram_id"` // FK to diagrams
	Content   string    `bson:"content" json:"content"`
	X         int       `bson:"x" json:"x"`           // Canvas position
	Y         int       `bson:"y" json:"y"`
	Width     int       `bson:"width" json:"width"`   // Note dimensions
	Height    int       `bson:"height" json:"height"`
	Color     string    `bson:"color" json:"color"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// NoteRepository defines methods for note data access
type NoteRepository interface {
	Store(ctx context.Context, n *Note) error
	StoreMultiple(ctx context.Context, notes []Note) error
	UpdateByDiagramID(ctx context.Context, diagramID string, notes []Note) error
	GetByDiagramID(ctx context.Context, diagramID string) ([]Note, error)
	DeleteByDiagramID(ctx context.Context, diagramID string) error
}
