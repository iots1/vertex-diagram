package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/iots1/vertex-diagram/delivery/http"
	"github.com/iots1/vertex-diagram/infrastructure/config"
	"github.com/iots1/vertex-diagram/infrastructure/database"
	"github.com/iots1/vertex-diagram/repository"
	"github.com/iots1/vertex-diagram/usecase"

	"time"
)

func main() {
	// 1. Load Config (.env)
	cfg := config.LoadConfig()

	// 2. Connect Database (Singleton)
	dbClient, err := database.GetMongoClient(cfg.MongoURI)
	if err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}
	// อย่าลืมปิด DB เมื่อปิดโปรแกรม
	defer database.CloseMongoDB()

	// เลือก Database และ Collection
	db := dbClient.Database(cfg.DBName)
	col := db.Collection("diagrams")

	// 2. Initialize Fiber Web Server with larger body limit for big SQL diagrams
	app := fiber.New(fiber.Config{
		BodyLimit: 50 * 1024 * 1024, // 50MB
	})
	
	// Config CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // More flexible for development
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// 4. Clean Architecture Wiring
	// Repo -> Usecase -> Handler
	repo := repository.NewMongoRepository(col)
	uc := usecase.NewDiagramUsecase(repo, 5*time.Second)
	http.NewDiagramHandler(app, uc)

	// Global Config
	configCol := db.Collection("config")
	configRepo := repository.NewMongoConfigRepository(configCol)
	configUc := usecase.NewConfigUsecase(configRepo, 5*time.Second)
	http.NewConfigHandler(app, configUc)

	log.Printf("🚀 Vertex Backend running on :%s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}
}