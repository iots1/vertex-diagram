package domain

import (
	"context"
	"time"
)

// Table represents a database table in a diagram
type Table struct {
	ID        string      `bson:"_id,omitempty" json:"id"`
	DiagramID string      `bson:"diagram_id" json:"diagram_id"` // FK to diagrams
	TableID   string      `bson:"table_id" json:"table_id"`     // ID from diagram
	Name      string      `bson:"name" json:"name"`
	Schema    string      `bson:"schema" json:"schema"`
	Fields    interface{} `bson:"fields" json:"fields"`         // JSON array
	Indexes   interface{} `bson:"indexes" json:"indexes"`       // JSON array
	Color     string      `bson:"color" json:"color"`
	X         int         `bson:"x" json:"x"`
	Y         int         `bson:"y" json:"y"`
	IsView    bool        `bson:"isView" json:"isView"`
	Order     int         `bson:"order" json:"order"`
	CreatedAt time.Time   `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time   `bson:"updated_at" json:"updated_at"`
}

// TableRepository defines methods for table data access
type TableRepository interface {
	Store(ctx context.Context, t *Table) error
	StoreMultiple(ctx context.Context, tables []Table) error
	UpdateByDiagramID(ctx context.Context, diagramID string, tables []Table) error
	GetByDiagramID(ctx context.Context, diagramID string) ([]Table, error)
	DeleteByDiagramID(ctx context.Context, diagramID string) error
}
