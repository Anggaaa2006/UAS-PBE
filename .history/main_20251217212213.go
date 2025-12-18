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

// ✅ TAMBAHAN
adminUserService := service.NewAdminUserService(userPGRepo)

// ===============================
// 4. Init Controller
// ===============================
authController := controller.NewAuthController(authService)
achievementController := controller.NewAchievementController(achievementService)
statsController := controller.NewStatsController(statsService)

adminAchievementController :=
    controller.NewAdminAchievementController(adminAchievementService)

// ✅ TAMBAHAN
adminUserController :=
    controller.NewAdminUserController(adminUserService)

// ===============================
// 6. Register All Routes
// ===============================
router.RegisterRoutes(
    r,
    authController,
    achievementController,
    statsController,
    adminAchievementController,
    adminUserController,
)
