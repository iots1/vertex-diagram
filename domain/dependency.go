package domain

import (
	"context"
	"time"
)

// Dependency represents a table dependency in a diagram
type Dependency struct {
	ID                  string    `bson:"_id,omitempty" json:"id"`
	DiagramID           string    `bson:"diagram_id" json:"diagram_id"` // FK to diagrams
	DependencyID        string    `bson:"dependency_id" json:"dependency_id"` // ID from diagram
	Schema              string    `bson:"schema" json:"schema"`
	TableID             string    `bson:"table_id" json:"table_id"`
	DependentSchema     string    `bson:"dependent_schema" json:"dependent_schema"`
	DependentTableID    string    `bson:"dependent_table_id" json:"dependent_table_id"`
	CreatedAt           time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt           time.Time `bson:"updated_at" json:"updated_at"`
}

// DependencyRepository defines methods for dependency data access
type DependencyRepository interface {
	Store(ctx context.Context, d *Dependency) error
	StoreMultiple(ctx context.Context, dependencies []Dependency) error
	UpdateByDiagramID(ctx context.Context, diagramID string, dependencies []Dependency) error
	GetByDiagramID(ctx context.Context, diagramID string) ([]Dependency, error)
	DeleteByDiagramID(ctx context.Context, diagramID string) error
}
