package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
	dsn := os.Getenv("DB_DSN")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to open DB:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

	log.Println("PostgreSQL connected successfully!")
	return db
}
