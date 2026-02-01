package database

import (
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	clientInstance *mongo.Client
	mongoOnce      sync.Once
	mongoError     error
)

// GetMongoClient คืนค่า Connection เดิมเสมอ (Singleton)
func GetMongoClient(uri string) (*mongo.Client, error) {
	// sync.Once จะทำงานแค่ครั้งแรกที่ถูกเรียก
	mongoOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		clientOptions := options.Client().ApplyURI(uri)
		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			mongoError = err
			return
		}

		// Ping เช็คว่าต่อได้จริงไหม
		err = client.Ping(ctx, nil)
		if err != nil {
			mongoError = err
			return
		}

		log.Println("✅ Connected to MongoDB (Singleton Instance)")
		clientInstance = client
	})

	return clientInstance, mongoError
}

// CreateIndexes creates indexes for all collections with diagram_id foreign key
func CreateIndexes(db *mongo.Database) error {
	collections := []string{"tables", "relationships", "dependencies", "areas", "custom_types", "notes", "diagram_filters"}

	for _, collectionName := range collections {
		_, err := db.Collection(collectionName).Indexes().CreateOne(
			context.Background(),
			mongo.IndexModel{
				Keys: bson.D{{Key: "diagram_id", Value: 1}},
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// CloseMongoDB ปิด Connection เมื่อจบโปรแกรม
func CloseMongoDB() {
	if clientInstance != nil {
		if err := clientInstance.Disconnect(context.TODO()); err != nil {
			log.Printf("❌ Error disconnecting MongoDB: %v\n", err)
		} else {
			log.Println("👋 MongoDB connection closed")
		}
	}
}