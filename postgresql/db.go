package postgresql

import (
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/zap"
)

var DB *sql.DB

func InitDB(dataSourceName string, logger *zap.Logger) *sql.DB {
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		logger.Fatal("Error connecting to the database: ", zap.Error(err))
	}

	if err := db.Ping(); err != nil {
		logger.Fatal("Unable to reach the database: ", zap.Error(err))
	}

	logger.Info("Connected to the database successfully")

	DB = db
	return DB
}
