package domain

import (
	"context"
	"time"
)

// Area represents a visual grouping area on the canvas for organizing tables
type Area struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	DiagramID string    `bson:"diagram_id" json:"diagram_id"` // FK to diagrams
	Name      string    `bson:"name" json:"name"`
	X         int       `bson:"x" json:"x"`           // Canvas position
	Y         int       `bson:"y" json:"y"`
	Width     int       `bson:"width" json:"width"`   // Area dimensions
	Height    int       `bson:"height" json:"height"`
	Color     string    `bson:"color" json:"color"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// AreaRepository defines methods for area data access
type AreaRepository interface {
	Store(ctx context.Context, a *Area) error
	StoreMultiple(ctx context.Context, areas []Area) error
	UpdateByDiagramID(ctx context.Context, diagramID string, areas []Area) error
	GetByDiagramID(ctx context.Context, diagramID string) ([]Area, error)
	DeleteByDiagramID(ctx context.Context, diagramID string) error
}
