package main

import (
	"Stock-Suggester-API/internal/database"
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func StartDB() {
	err := godotenv.Load("data.env")
	if err != nil {
		fmt.Println(err)
		return
	}
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("sqlite", dbURL)

	if err != nil {
		fmt.Println(err)
		return
	}

	cfg.db = database.New(db)
}

type Config struct {
	db *database.Queries
}

var cfg Config
