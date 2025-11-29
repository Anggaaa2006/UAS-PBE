package postgres

import (
	"database/sql"
	_ "github.com/lib/pq" // driver PostgreSQL
)

/*
	NewPostgresConn
	Membuat koneksi ke database PostgreSQL.
	Parameter:
	- dsn = connection string PostgreSQL

	Return:
	- *sql.DB = object koneksi database
	- error   = error jika gagal connect
*/
func NewPostgresConn(dsn string) (*sql.DB, error) {

	// Membuka koneksi ke database PostgreSQL
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Test koneksi
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
