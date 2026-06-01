package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"MONITORING-TOOL/internal/config"
)

func NewPostgresConnection(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("Connected To PostgreSQL")
	return db, nil
}
