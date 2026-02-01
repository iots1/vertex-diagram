package domain

import (
	"context"
)

type Config struct {
	ID               string `bson:"_id" json:"id"`
	DefaultDiagramID string `bson:"default_diagram_id" json:"defaultDiagramId"`
}

type ConfigRepository interface {
	Get(ctx context.Context) (*Config, error)
	Upsert(ctx context.Context, c *Config) error
}

type ConfigUsecase interface {
	Get(ctx context.Context) (*Config, error)
	Save(ctx context.Context, c *Config) error
}
