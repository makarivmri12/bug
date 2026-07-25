package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// InitPostgres initializes PostgreSQL connection
func InitPostgres(dsn string, logger *zap.Logger) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error("failed to open database", zap.Error(err))
		return nil, err
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		logger.Error("failed to ping database", zap.Error(err))
		return nil, err
	}

	logger.Info("database connected successfully")
	return db, nil
}

// Close closes the database connection
func Close(db *sql.DB, logger *zap.Logger) error {
	if err := db.Close(); err != nil {
		logger.Error("failed to close database", zap.Error(err))
		return err
	}
	logger.Info("database connection closed")
	return nil
}
