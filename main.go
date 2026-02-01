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

	// 3. Setup Fiber
	app := fiber.New()
	
	// Config CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173", // Frontend URL
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// 4. Clean Architecture Wiring
	// Repo -> Usecase -> Handler
	repo := repository.NewMongoRepository(col)
	uc := usecase.NewDiagramUsecase(repo, 5*time.Second) // เพิ่ม Timeout ให้เหมาะสม
	
	// Register Routes
	http.NewDiagramHandler(app, uc)

	log.Printf("🚀 Vertex Backend running on :%s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}
}