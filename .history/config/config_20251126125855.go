package config

import "os"

type Config struct {
	PostgresDSN string
	MongoURI    string
	MongoDBName string
	JWTSecret   string
}

func Load() *Config {
	return &Config{
		PostgresDSN: getOr("POSTGRES_DSN", "postgres://dev:dev@localhost:5432/uas?sslmode=disable"),
		MongoURI:    getOr("MONGO_URI", "mongodb://localhost:27017"),
		MongoDBName: getOr("MONGO_DB_NAME", "uas"),
		JWTSecret:   getOr("JWT_SECRET", "secret123"),
	}
}

func getOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
