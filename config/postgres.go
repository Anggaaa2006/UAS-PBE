package config

import (
	"database/sql"

	_ "github.com/lib/pq"
)

/*
	InitPostgres
	Menghubungkan aplikasi ke PostgreSQL menggunakan DSN dari .env
*/
func InitPostgres(cfg Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.PostgresURL)
	if err != nil {
		return nil, err
	}

	// Coba ping dulu
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
