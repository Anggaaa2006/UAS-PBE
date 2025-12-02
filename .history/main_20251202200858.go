package main

import (
	"log"
	"os"

	"uas_pbe/config"
	"uas_pbe/controller"
	"uas_pbe/repository/mongo"
	"uas_pbe/repository/postgres"
	"uas_pbe/router"
	"uas_pbe/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// ============================
	// 1. Load .env
	// ============================
	if err := godotenv.Load(); err != nil {
		log.Println("WARNING: .env tidak ditemukan, menggunakan default ENV")
	}

	// ============================
	// 2. Connect Database
	// ============================

	// PostgreSQL
	pgDB := config.ConnectPostgres()
	if pgDB == nil {
		log.Fatal("Gagal konek PostgreSQL")
	}

	// MongoDB
	mongoClient, mongoDB := config.ConnectMongo()
	if mongoClient == nil {
		log.Fatal("Gagal konek MongoDB")
	}

	// ============================
	// 3. Inisialisasi Repository
	// ============================

	// Repo Achievement di PostgreSQL
	achievementRepo := postgres.NewAchievementRepo(pgDB)

	// Repo AchievementDetail di Mongo
	achievementDetailRepo := mongo.NewAchievementDetailRepo(mongoDB)

	// Repo User (untuk login/register)
	userRepo := postgres.NewUserRepo(pgDB)

	// ============================
	// 4. Inisialisasi Service
	// ============================

	// Auth (Login + Register)
	authService := service.NewAuthService(userRepo)

	// Achievement Service (CRUD + submit + approve)
	achievementService := service.NewAchievementService(achievementRepo, achievementDetailRepo)

	// ============================
	// 5. Inisialisasi Controller
	// ============================

	authController := controller.NewAuthController(authService)
	achievementController := controller.NewAchievementController(achievementService)

	// ============================
	// 6. Inisialisasi Gin Router
	// ============================
	r := gin.Default()

	// ============================
	// 7. Daftarkan Routes
	// ============================

	// Auth Route
	router.RegisterAuthRoutes(r, authController)

	// Achievement Route
	router.RegisterAchievementRoutes(r, *achievementController)

	// ============================
	// 8. Run Server
	// ============================
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port " + port)
	r.Run(":" + port)
}
