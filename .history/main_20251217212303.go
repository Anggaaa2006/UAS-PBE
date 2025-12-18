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

