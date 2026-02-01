package domain

import (
	"context"
	"time"
)

// CustomType represents a custom database type definition in a diagram
type CustomType struct {
	ID        string      `bson:"_id,omitempty" json:"id"`
	DiagramID string      `bson:"diagram_id" json:"diagram_id"` // FK to diagrams
	Schema    string      `bson:"schema" json:"schema"`
	Type      string      `bson:"type" json:"type"`   // Type name
	Kind      string      `bson:"kind" json:"kind"`   // e.g., 'enum', 'composite', etc.
	Values    interface{} `bson:"values" json:"values"`   // JSON array for enum values
	Fields    interface{} `bson:"fields" json:"fields"`   // JSON array for composite fields
	CreatedAt time.Time   `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time   `bson:"updated_at" json:"updated_at"`
}

// CustomTypeRepository defines methods for custom type data access
type CustomTypeRepository interface {
	Store(ctx context.Context, ct *CustomType) error
	StoreMultiple(ctx context.Context, customTypes []CustomType) error
	UpdateByDiagramID(ctx context.Context, diagramID string, customTypes []CustomType) error
	GetByDiagramID(ctx context.Context, diagramID string) ([]CustomType, error)
	DeleteByDiagramID(ctx context.Context, diagramID string) error
}
