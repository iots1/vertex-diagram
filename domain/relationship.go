package domain

import (
	"context"
	"time"
)

// Relationship represents a relationship between tables in a diagram
type Relationship struct {
	ID             string    `bson:"_id,omitempty" json:"id"`
	DiagramID      string    `bson:"diagram_id" json:"diagram_id"` // FK to diagrams
	RelationshipID string    `bson:"relationship_id" json:"relationship_id"` // ID from diagram
	Name           string    `bson:"name" json:"name"`
	SourceTableID  string    `bson:"source_table_id" json:"source_table_id"`
	TargetTableID  string    `bson:"target_table_id" json:"target_table_id"`
	SourceFieldID  string    `bson:"source_field_id" json:"source_field_id"`
	TargetFieldID  string    `bson:"target_field_id" json:"target_field_id"`
	Type           string    `bson:"type" json:"type"` // relationship type
	CreatedAt      time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time `bson:"updated_at" json:"updated_at"`
}

// RelationshipRepository defines methods for relationship data access
type RelationshipRepository interface {
	Store(ctx context.Context, r *Relationship) error
	StoreMultiple(ctx context.Context, relationships []Relationship) error
	UpdateByDiagramID(ctx context.Context, diagramID string, relationships []Relationship) error
	GetByDiagramID(ctx context.Context, diagramID string) ([]Relationship, error)
	DeleteByDiagramID(ctx context.Context, diagramID string) error
}
