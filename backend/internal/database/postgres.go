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
	
	// Run schema migrations
	if err := RunMigrations(db); err != nil {
		log.Printf("Failed to run database migrations: %v\n", err)
	}

	return db, nil
}

func RunMigrations(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS configurations (
			id SERIAL PRIMARY KEY,
			polling_frequency INTEGER,
			updated_at TIMESTAMP DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS alert_thresholds (
			id SERIAL PRIMARY KEY,
			cpu_threshold FLOAT,
			memory_threshold FLOAT,
			disk_threshold FLOAT,
			network_threshold FLOAT
		);`,
		`CREATE TABLE IF NOT EXISTS email_configuration (
			id SERIAL PRIMARY KEY,
			smtp_server TEXT,
			smtp_port INTEGER,
			username TEXT,
			password TEXT,
			from_email TEXT,
			primary_recipient TEXT,
			alternate_recipient TEXT,
			subject_template TEXT,
			body_template TEXT
		);`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return fmt.Errorf("query failed: %w (query: %s)", err, q)
		}
	}
	log.Println("Database migrations run successfully")
	return nil
}

