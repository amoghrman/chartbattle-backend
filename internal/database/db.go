package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	connStr := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Database not reachable:", err)
	}

	DB = db
	log.Println("Connected to database successfully")
}
