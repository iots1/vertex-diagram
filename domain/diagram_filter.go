package domain

import (
	"context"
	"time"
)

// DiagramFilter represents filter settings for a diagram view
type DiagramFilter struct {
	ID         string        `bson:"_id,omitempty" json:"id"`
	DiagramID  string        `bson:"diagram_id" json:"diagram_id"` // FK to diagrams (unique)
	TableIDs   []string      `bson:"table_ids" json:"table_ids"`   // Filtered table IDs
	SchemaIDs  []string      `bson:"schema_ids" json:"schema_ids"` // Filtered schema IDs
	CreatedAt  time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time     `bson:"updated_at" json:"updated_at"`
}

// DiagramFilterRepository defines methods for diagram filter data access
type DiagramFilterRepository interface {
	Store(ctx context.Context, df *DiagramFilter) error
	GetByDiagramID(ctx context.Context, diagramID string) (*DiagramFilter, error)
	DeleteByDiagramID(ctx context.Context, diagramID string) error
}
