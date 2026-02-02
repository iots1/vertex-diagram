package domain

import (
	"context"
	"time"
)

// Table represents a database table in a diagram
type Table struct {
	ID        string                   `bson:"_id,omitempty" json:"mongoId"`
	DiagramID string                   `bson:"diagram_id" json:"diagram_id"` // FK to diagrams
	TableID   string                   `bson:"table_id" json:"id"`           // ID from diagram - returned as "id" for frontend
	Name      string                   `bson:"name" json:"name"`
	Schema    string                   `bson:"schema" json:"schema"`
	Fields    []map[string]interface{} `bson:"fields" json:"fields"`         // Array of field objects
	Indexes   []map[string]interface{} `bson:"indexes" json:"indexes"`       // Array of index objects
	Color     string                   `bson:"color" json:"color"`
	X         int                      `bson:"x" json:"x"`
	Y         int                      `bson:"y" json:"y"`
	IsView    bool                     `bson:"isView" json:"isView"`
	Order     int                      `bson:"order" json:"order"`
	CreatedAt time.Time                `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time                `bson:"updated_at" json:"updatedAt"`
}

// TableRepository defines methods for table data access
type TableRepository interface {
	Store(ctx context.Context, t *Table) error
	StoreMultiple(ctx context.Context, tables []Table) error
	UpdateByDiagramID(ctx context.Context, diagramID string, tables []Table) error
	GetByDiagramID(ctx context.Context, diagramID string) ([]Table, error)
	DeleteByDiagramID(ctx context.Context, diagramID string) error
}
