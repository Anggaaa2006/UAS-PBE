package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"uas_pbe/controller"
	"uas_pbe/repository/mongo"
	"uas_pbe/repository/postgres"
	"uas_pbe/router"
	"uas_pbe/service"
)

func main() {

	/*
		==========================================
		1. Load ENV
		==========================================
	*/
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file tidak ditemukan (ini tidak masalah di production)")
	}

	postgresDSN := os.Getenv("POSTGRES_DSN")
	mongoURI := os.Getenv("MONGO_URI")
	mongoDBName := os.Getenv("MONGO_DB")

	if postgresDSN == "" || mongoURI == "" || mongoDBName == "" {
		log.Fatal("ENV belum lengkap. Pastikan POSTGRES_DSN, MONGO_URI, dan MONGO_DB diisi.")
	}

	/*
		==========================================
		2. Connect ke PostgreSQL
		==========================================
	*/
	pgDB, err := postgres.NewPostgresConn(postgresDSN)
	if err != nil {
		log.Fatal("Gagal connect ke PostgreSQL:", err)
	}
	fmt.Println("PostgreSQL connected ✔")

	/*
		==========================================
		3. Connect ke MongoDB
		==========================================
	*/
	ctx := context.Background()
	mongoClient, err := mongo.NewMongoConn(ctx, mongoURI)
	if err != nil {
		log.Fatal("Gagal connect ke MongoDB:", err)
	}
	mongoDB := mongoClient.Database(mongoDBName)
	fmt.Println("MongoDB connected ✔")

	/*
		==========================================
		4. Inisialisasi Repository
		==========================================
	*/
	achievementRefRepo := postgres.NewAchievementReferenceRepo(pgDB)
	achievementDetailRepo := mongo.NewAchievementDetailRepo(mongoDB)

	/*
		==========================================
		5. Inisialisasi Service
		==========================================
	*/
	achievementService := service.NewAchievementService(
		achievementRefRepo,
		achievementDetailRepo,
	)

	/*
		==========================================
		6. Inisialisasi Controller
		==========================================
	*/
	achievementController := controller.NewAchievementController(achievementService)

	/*
		==========================================
		7. Setup Router
		==========================================
	*/
	r := gin.Default()

	// Tambahkan route prestasi
	router.RegisterAchievementRoutes(r, achievementController)

	// Endpoint test
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	/*
		==========================================
		8. Run Server
		==========================================
	*/
	fmt.Println("Server running on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Server gagal berjalan:", err)
	}
}
