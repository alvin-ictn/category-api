package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDB(connectionString string) (*sql.DB, error) {
	// Open database
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	// Test connection
	err = db.Ping()
	if err != nil {
		log.Println("Database connection failed (will retry):", err)
		return db, err // Return db object anyway so we can keep trying
	}

	// Set connection pool settings (optional but recommended)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("Database connected successfully")
	return db, nil
}
