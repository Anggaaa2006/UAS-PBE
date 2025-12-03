package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

/*
	Config struct
	Menampung semua konfigurasi aplikasi dari .env
*/
type Config struct {
	PostgresURL string
	MongoURL    string
	MongoDBName string
}

/*
	LoadConfig
	Membaca file .env dan memasukkan ke struct Config
*/
func LoadConfig() Config {
	// Load .env jika ada
	_ = godotenv.Load()

	cfg := Config{
		PostgresURL: os.Getenv("POSTGRES_URL"),
		MongoURL:    os.Getenv("MONGO_URL"),
		MongoDBName: os.Getenv("MONGO_DB"),
	}

	// Validasi sederhana
	if cfg.PostgresURL == "" {
		log.Println("WARNING: POSTGRES_URL kosong")
	}
	if cfg.MongoURL == "" {
		log.Println("WARNING: MONGO_URL kosong")
	}

	return cfg
}
