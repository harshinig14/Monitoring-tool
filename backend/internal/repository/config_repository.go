package repository

import (
	"database/sql"
)

type ConfigRepository struct {
	DB *sql.DB
}

func NewConfigRepository(db *sql.DB) *ConfigRepository {
	return &ConfigRepository{DB: db}
}

func (r *ConfigRepository) GetConfiguration() (int, error) {
	var frequency int
	query := `SELECT polling_frequency FROM configurations ORDER BY id DESC LIMIT 1`
	err := r.DB.QueryRow(query).Scan(&frequency)
	if err != nil {
		return 60, nil // Default fallback is 60 seconds (also catches missing tables)
	}
	return frequency, nil
}

func (r *ConfigRepository) SaveConfiguration(frequency int) error {
	query := `INSERT INTO configurations (polling_frequency) VALUES ($1)`
	_, err := r.DB.Exec(query, frequency)
	return err
}
