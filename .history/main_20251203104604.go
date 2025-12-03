package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"uas_pbe/config"
	"uas_pbe/controller"
	"uas_pbe/router"

	"uas_pbe/repository/postgres"
	"uas_pbe/repository/mongo"
	"uas_pbe/service"
)

func main() {

	// ===============================
	// 1. Load ENV + connect database
	// ===============================
	cfg := config.LoadConfig()

	// Connect PostgreSQL
	pgDB, err := config.InitPostgres(cfg)
	if err != nil {
		log.Fatal("failed connect postgres:", err)
	}

	// Connect MongoDB
	mongoDB, err := config.InitMongo(cfg)
	if err != nil {
		log.Fatal("failed connect mongo:", err)
	}

	// =======================================
	// 2. Init Repository (Postgres + Mongo)
	// =======================================

	achievementPGRepo := postgres.NewAchievementReferenceRepo(pgDB)
	achievementDetailRepo := mongo.NewAchievementDetailRepo(mongoDB)

	userPGRepo := postgres.NewUserRepo(pgDB) // untuk login

	// ===============================
	// 3. Init Service
	// ===============================
	authService := service.NewAuthService(userPGRepo)
	achievementService := service.NewAchievementService(achievementPGRepo, achievementDetailRepo)

	// ===============================
	// 4. Init Controller
	// ===============================
	authController := controller.NewAuthController(authService)
	achievementController := controller.NewAchievementController(achievementService)

	// ===============================
	// 5. Init Gin
	// ===============================
	r := gin.Default()

	// ===============================
	// 6. Register All Routes
	// ===============================
	router.RegisterRoutes(
		r,
		authController,
		achievementController,
	)

	// ===============================
	// 7. Run Server
	// ===============================
	log.Println("Server running on port 8080")
	r.Run(":8080")
}
