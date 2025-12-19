package main
// @title UAS PBE API
// @version 1.0
// @description Sistem Manajemen Prestasi Mahasiswa (Test User)
// @termsOfService http://example.com/terms/

// @contact.name Test User
// @contact.email testuser@example.com

// @host localhost:8080
// @BasePath /

import (
	"log"

	"github.com/gin-gonic/gin"

	"uas_pbe/config"
	"uas_pbe/controller"
	"uas_pbe/router"

	"uas_pbe/repository/mongo"
	"uas_pbe/repository/postgres"
	"uas_pbe/service"
)

func main() {

	// ===============================
	// 1. Load ENV + connect database
	// ===============================
	cfg := config.LoadConfig()

	pgDB, err := config.InitPostgres(cfg)
	if err != nil {
		log.Fatal("failed connect postgres:", err)
	}

	mongoDB, err := config.InitMongo(cfg)
	if err != nil {
		log.Fatal("failed connect mongo:", err)
	}

	// ===============================
	// 2. Init Repository
	// ===============================
	achievementPGRepo := postgres.NewAchievementReferenceRepo(pgDB)
	achievementDetailRepo := mongo.NewAchievementDetailRepo(mongoDB)
	userPGRepo := postgres.NewUserRepo(pgDB)

	// ===============================
	// 3. Init Service
	// ===============================
	authService := service.NewAuthService(userPGRepo)

	achievementService := service.NewAchievementService(
		achievementPGRepo,
		achievementDetailRepo,
	)

	statsService := service.NewStatsService(achievementPGRepo)

	adminAchievementService := service.NewAdminAchievementService(
		achievementPGRepo,
		achievementDetailRepo,
	)

	adminUserService := service.NewAdminUserService(userPGRepo)

	// ===============================
	// 4. Init Controller
	// ===============================
	authController := controller.NewAuthController(authService)
	achievementController := controller.NewAchievementController(achievementService)
	statsController := controller.NewStatsController(statsService)

	adminAchievementController :=
		controller.NewAdminAchievementController(adminAchievementService)

	adminUserController :=
		controller.NewAdminUserController(adminUserService)

	// ✅ DASHBOARD CONTROLLER (INI YANG BENAR)
	dashboardController :=
		controller.NewDashboardController(
			achievementService, // *AchievementService
			statsService,       // StatsService
		)

	// ===============================
	// 5. Init Gin
	// ===============================
	r := gin.Default()

	// ===============================
	// 6. Register Routes
	// ===============================
	router.RegisterRoutes(
		r,
		authController,
		achievementController,
		statsController,
		adminAchievementController,
		adminUserController,
		dashboardController, // ⬅️ TAMBAHAN PENTING
	)

	// ===============================
	// 7. Run Server
	// ===============================
	log.Println("Server running on port 8080")
	r.Run(":8080")
}
