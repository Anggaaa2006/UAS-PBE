package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"uas_pbe/config"
	"uas_pbe/router"
	"UAS PBE/repository"
	"UAS PBE/service"
	"UAS PBE/controller"

	"context"
)

func main() {
	// ======================================================
	// 1. LOAD CONFIG (.env)
	// ======================================================
	cfg := config.LoadConfig()

	// ======================================================
	// 2. CONNECT POSTGRESQL
	// ======================================================
	postgresDB, err := repository.NewPostgresConn(cfg.PostgresDSN)
	if err != nil {
		log.Fatal("Failed to connect PostgreSQL:", err)
	}

	// ======================================================
	// 3. CONNECT MONGODB
	// ======================================================
	mongoClient, err := repository.NewMongoConn(context.Background(), cfg.MongoURI)
	if err != nil {
		log.Fatal("Failed to connect MongoDB:", err)
	}
	mongoDB := mongoClient.Database(cfg.MongoDBName)

	// ======================================================
	// 4. INIT REPOSITORY
	// ======================================================
	userRepo := repository.NewUserRepo(postgresDB)
	roleRepo := repository.NewRoleRepo(postgresDB)
	studentRepo := repository.NewStudentRepo(postgresDB)
	achRefRepo := repository.NewAchievementReferenceRepo(postgresDB)

	achDetailRepo := repository.NewAchievementDetailRepo(mongoDB)

	// ======================================================
	// 5. INIT SERVICE
	// ======================================================
	authService := service.NewAuthService(userRepo, roleRepo)
	studentService := service.NewStudentService(studentRepo)
	achievementService := service.NewAchievementService(achRefRepo, achDetailRepo)

	// ======================================================
	// 6. INIT CONTROLLER
	// ======================================================
	authController := controller.NewAuthController(authService)
	studentController := controller.NewStudentController(studentService)
	achievementController := controller.NewAchievementController(achievementService)

	// ======================================================
	// 7. SETUP ROUTER
	// ======================================================
	r := gin.Default()
	router.SetupRoutes(r, authController, studentController, achievementController)

	// ======================================================
	// 8. RUN SERVER
	// ======================================================
	log.Println("Server running on", cfg.AppPort)
	r.Run(cfg.AppPort)
}
