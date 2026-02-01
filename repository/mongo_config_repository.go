package repository

import (
	"context"

	"github.com/iots1/vertex-diagram/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoConfigRepository struct {
	Conn *mongo.Collection
}

func NewMongoConfigRepository(Conn *mongo.Collection) domain.ConfigRepository {
	return &mongoConfigRepository{Conn}
}

func (m *mongoConfigRepository) Get(ctx context.Context) (*domain.Config, error) {
	var c domain.Config
	err := m.Conn.FindOne(ctx, bson.M{"_id": "global"}).Decode(&c)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &domain.Config{ID: "global"}, nil
		}
		return nil, err
	}
	return &c, nil
}

func (m *mongoConfigRepository) Upsert(ctx context.Context, c *domain.Config) error {
	c.ID = "global"
	filter := bson.M{"_id": "global"}
	update := bson.M{"$set": c}
	opts := options.Update().SetUpsert(true)
	_, err := m.Conn.UpdateOne(ctx, filter, update, opts)
	return err
}
