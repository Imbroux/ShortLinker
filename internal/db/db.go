package db

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
)

var DB *sql.DB

func InitDB(dataSourceName string) {
	var err error
	DB, err = sql.Open("pgx", dataSourceName)
	if err != nil {
		log.Fatalf("Error connecting to db: %v", err)
	}

	log.Println("Database connection established.")
}
