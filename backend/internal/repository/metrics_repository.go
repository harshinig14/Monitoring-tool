package repository

import (
	"database/sql"
	"MONITORING-TOOL/internal/models"
)

type MetricsRepository struct {
	DB *sql.DB
}

func NewMetricsRepository(db *sql.DB) *MetricsRepository {
	return &MetricsRepository{DB: db}
}

func (r *MetricsRepository) InsertMetrics(req models.MetricsRequest) error {
	query := `
		INSERT INTO trace_results (user_id, cpu_usage, memory_usage, disk_usage, network_usage)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.DB.Exec(query, req.UserID, req.CPUUsage, req.MemoryUsage, req.DiskUsage, req.NetworkUsage)
	return err
}
