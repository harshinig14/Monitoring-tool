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
			cpu_threshold FLOAT NOT NULL DEFAULT 80,
			memory_threshold FLOAT NOT NULL DEFAULT 80,
			disk_threshold FLOAT NOT NULL DEFAULT 80,
			network_threshold FLOAT NOT NULL DEFAULT 80,
			updated_at TIMESTAMP DEFAULT NOW()
		);`,
		`INSERT INTO alert_thresholds (cpu_threshold, memory_threshold, disk_threshold, network_threshold)
		 SELECT 80, 80, 80, 80 WHERE NOT EXISTS (SELECT 1 FROM alert_thresholds);`,
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
			body_template TEXT,
			updated_at TIMESTAMP DEFAULT NOW()
		);`,
		`ALTER TABLE email_configuration ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT NOW();`,
		`INSERT INTO email_configuration (smtp_server, smtp_port, updated_at)
		 SELECT '', 587, NOW() WHERE NOT EXISTS (SELECT 1 FROM email_configuration);`,
		`CREATE TABLE IF NOT EXISTS alerts (
			alert_id BIGSERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			machine_name VARCHAR(255),
			metric_name VARCHAR(50),
			current_value FLOAT,
			threshold_value FLOAT,
			severity VARCHAR(20),
			created_at TIMESTAMP DEFAULT NOW(),
			FOREIGN KEY(user_id) REFERENCES users(user_id)
		);`,
		`CREATE TABLE IF NOT EXISTS email_logs (
			email_id BIGSERIAL PRIMARY KEY,
			recipient VARCHAR(255),
			subject TEXT,
			status VARCHAR(50),
			sent_at TIMESTAMP DEFAULT NOW()
		);`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS device_status VARCHAR(50) DEFAULT 'PENDING_REGISTRATION';`,
		`CREATE TABLE IF NOT EXISTS device_status_history (
			id BIGSERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			previous_status VARCHAR(50),
			current_status VARCHAR(50),
			changed_at TIMESTAMP DEFAULT NOW(),
			FOREIGN KEY(user_id) REFERENCES users(user_id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS device_requests (
			request_id BIGSERIAL PRIMARY KEY,
			os_type VARCHAR(50),
			requested_by VARCHAR(255),
			created_at TIMESTAMP DEFAULT NOW()
		);`,
		`UPDATE users SET device_status = 'ONLINE' WHERE device_status IS NULL AND last_seen IS NOT NULL AND NOW() - last_seen < INTERVAL '60 seconds';`,
		`UPDATE users SET device_status = 'OFFLINE' WHERE device_status IS NULL AND last_seen IS NOT NULL AND NOW() - last_seen >= INTERVAL '60 seconds';`,
		`UPDATE users SET device_status = 'PENDING_REGISTRATION' WHERE device_status IS NULL;`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			log.Printf("Migration query failed (skipping): %v\n", err)
		}
	}
	log.Println("Database migrations process finished")
	return nil
}

