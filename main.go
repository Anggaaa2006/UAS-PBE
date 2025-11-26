package main

import (
	"context"
	"log"

	"UAS PBE/config"
	pg "UAS PBE/internal/repository/postgres"
	mg "UAS PBE/internal/repository/mongo"
	"UAS PBE/internal/router"

	"UAS PBE/internal/service"
	"UAS PBE/internal/controller"
)

func main() {
	cfg := config.Load()

	// PostgreSQL
	pgDB, err := pg.NewPostgresConn(cfg.PostgresDSN)
	if err != nil {
		log.Fatal("Postgres error:", err)
	}

	// MongoDB
	mongoClient, err := mg.NewMongoClient(context.Background(), cfg.MongoURI)
	if err != nil {
		log.Fatal("Mongo error:", err)
	}
	mongoDB := mongoClient.Database(cfg.MongoDBName)

	// Repository instances
	userRepo := pg.NewUserRepo(pgDB)
	roleRepo := pg.NewRoleRepo(pgDB)
	refRepo := pg.NewAchievementRefRepo(pgDB)
	studentRepo := pg.NewStudentRepo(pgDB)
	lectRepo := pg.NewLecturerRepo(pgDB)
	achievementRepo := mg.NewAchievementRepo(mongoDB, "achievements")

	// Service layer
	authService := service.NewAuthService(userRepo, roleRepo, cfg)
	achievementService := service.NewAchievementService(achievementRepo, refRepo)
	userService := service.NewUserService(userRepo)
	studentService := service.NewStudentService(studentRepo)
	reportService := service.NewReportService(achievementRepo, refRepo)

	// Controllers
	authController := controller.NewAuthController(authService)
	achievementController := controller.NewAchievementController(achievementService)
	userController := controller.NewUserController(userService)
	studentController := controller.NewStudentController(studentService)
	reportController := controller.NewReportController(reportService)

	r := router.NewRouter(
		authController,
		achievementController,
		userController,
		studentController,
		reportController,
	)

	log.Println("Server running at :8080")
	r.Run(":8080")
}
