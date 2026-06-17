package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jereminathanael/healthqueue/internal/config"
	_ "github.com/lib/pq"
)

func Connection(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DBConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w ", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return db, nil
}