package domain

import (
	"context"
	"time"
)

// Relationship represents a relationship between tables in a diagram
type Relationship struct {
	ID                  string    `bson:"_id,omitempty" json:"mongoId"`
	DiagramID           string    `bson:"diagram_id" json:"diagram_id"` // FK to diagrams
	RelationshipID      string    `bson:"relationship_id" json:"id"`    // ID from diagram - returned as "id" for frontend
	Name                string    `bson:"name" json:"name"`
	SourceTableID       string    `bson:"source_table_id" json:"sourceTableId"` // Convert to camelCase in JSON
	TargetTableID       string    `bson:"target_table_id" json:"targetTableId"` // Convert to camelCase in JSON
	SourceFieldID       string    `bson:"source_field_id" json:"sourceFieldId"` // Convert to camelCase in JSON
	TargetFieldID       string    `bson:"target_field_id" json:"targetFieldId"` // Convert to camelCase in JSON
	Type                string    `bson:"type" json:"type"`
	SourceCardinality   string    `bson:"source_cardinality" json:"sourceCardinality"` // Convert to camelCase in JSON
	TargetCardinality   string    `bson:"target_cardinality" json:"targetCardinality"` // Convert to camelCase in JSON
	CreatedAt           time.Time `bson:"created_at" json:"createdAt"`           // Convert to camelCase in JSON
	UpdatedAt           time.Time `bson:"updated_at" json:"updatedAt"`           // Convert to camelCase in JSON
}

// RelationshipRepository defines methods for relationship data access
type RelationshipRepository interface {
	Store(ctx context.Context, r *Relationship) error
	StoreMultiple(ctx context.Context, relationships []Relationship) error
	UpdateByDiagramID(ctx context.Context, diagramID string, relationships []Relationship) error
	GetByDiagramID(ctx context.Context, diagramID string) ([]Relationship, error)
	DeleteByDiagramID(ctx context.Context, diagramID string) error
}
