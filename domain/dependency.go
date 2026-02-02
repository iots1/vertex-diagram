package domain

import (
	"context"
	"time"
)

// Dependency represents a table dependency in a diagram
type Dependency struct {
	ID                  string    `bson:"_id,omitempty" json:"mongoId"`
	DiagramID           string    `bson:"diagram_id" json:"diagram_id"` // FK to diagrams
	DependencyID        string    `bson:"dependency_id" json:"id"`      // ID from diagram - returned as "id" for frontend
	Schema              string    `bson:"schema" json:"schema"`
	TableID             string    `bson:"table_id" json:"tableId"`      // Convert to camelCase in JSON
	DependentSchema     string    `bson:"dependent_schema" json:"dependentSchema"` // Convert to camelCase in JSON
	DependentTableID    string    `bson:"dependent_table_id" json:"dependentTableId"` // Convert to camelCase in JSON
	CreatedAt           time.Time `bson:"created_at" json:"createdAt"` // Convert to camelCase in JSON
	UpdatedAt           time.Time `bson:"updated_at" json:"updatedAt"` // Convert to camelCase in JSON
}

// DependencyRepository defines methods for dependency data access
type DependencyRepository interface {
	Store(ctx context.Context, d *Dependency) error
	StoreMultiple(ctx context.Context, dependencies []Dependency) error
	UpdateByDiagramID(ctx context.Context, diagramID string, dependencies []Dependency) error
	GetByDiagramID(ctx context.Context, diagramID string) ([]Dependency, error)
	DeleteByDiagramID(ctx context.Context, diagramID string) error
}
