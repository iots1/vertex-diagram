package domain

import (
	"context"
	"time"
)

// Entity หลักของระบบ
type Diagram struct {
	ID        string                 `bson:"_id,omitempty" json:"id"`
	Name      string                 `bson:"name" json:"name"`
	Content   map[string]interface{} `bson:"content" json:"content"` // JSON ก้อนใหญ่ของ ChartDB
	UpdatedAt time.Time              `bson:"updated_at" json:"updated_at"`
	CreatedAt time.Time              `bson:"created_at" json:"created_at"`
}

// Repository Interface: สัญญาว่าต้องทำอะไรกับ DB ได้บ้าง
type DiagramRepository interface {
	Fetch(ctx context.Context) ([]Diagram, error)
	GetByID(ctx context.Context, id string) (*Diagram, error)
	Store(ctx context.Context, d *Diagram) error
	Update(ctx context.Context, d *Diagram) error
	Delete(ctx context.Context, id string) error
}

// Usecase Interface: สัญญาว่า Business Logic มีอะไรบ้าง
type DiagramUsecase interface {
	GetAll(ctx context.Context) ([]Diagram, error)
	GetOne(ctx context.Context, id string) (*Diagram, error)
	Save(ctx context.Context, d *Diagram) (*Diagram, error)
	Delete(ctx context.Context, id string) error
}