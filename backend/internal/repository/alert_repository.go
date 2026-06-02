package repository

import (
	"database/sql"
)

type AlertThresholds struct {
	CPUThreshold    float64 `json:"cpu_threshold"`
	MemoryThreshold float64 `json:"memory_threshold"`
	DiskThreshold   float64 `json:"disk_threshold"`
	NetworkThreshold float64 `json:"network_threshold"`
}

type AlertRepository struct {
	DB *sql.DB
}

func NewAlertRepository(db *sql.DB) *AlertRepository {
	return &AlertRepository{DB: db}
}

func (r *AlertRepository) GetThresholds() (*AlertThresholds, error) {
	var t AlertThresholds
	query := `SELECT cpu_threshold, memory_threshold, disk_threshold, network_threshold FROM alert_thresholds ORDER BY id DESC LIMIT 1`
	err := r.DB.QueryRow(query).Scan(&t.CPUThreshold, &t.MemoryThreshold, &t.DiskThreshold, &t.NetworkThreshold)
	if err != nil {
		return &AlertThresholds{CPUThreshold: 80, MemoryThreshold: 80, DiskThreshold: 80, NetworkThreshold: 80}, nil
	}
	return &t, nil
}

func (r *AlertRepository) SaveThresholds(t *AlertThresholds) error {
	query := `INSERT INTO alert_thresholds (cpu_threshold, memory_threshold, disk_threshold, network_threshold) VALUES ($1, $2, $3, $4)`
	_, err := r.DB.Exec(query, t.CPUThreshold, t.MemoryThreshold, t.DiskThreshold, t.NetworkThreshold)
	return err
}
