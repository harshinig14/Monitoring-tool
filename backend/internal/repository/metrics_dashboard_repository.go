package repository

import (
	"database/sql"
	"time"

	"MONITORING-TOOL/internal/dto"
)

type MetricsDashboardRepository struct {
	DB *sql.DB
}

func NewMetricsDashboardRepository(db *sql.DB) *MetricsDashboardRepository {
	return &MetricsDashboardRepository{DB: db}
}

func (r *MetricsDashboardRepository) GetLatestMetrics(userID int) (*dto.RealtimeMetrics, error) {
	query := `
		SELECT cpu_usage, memory_usage, disk_usage, network_usage, trace_date
		FROM trace_results
		WHERE user_id = $1
		ORDER BY trace_date DESC
		LIMIT 1
	`

	var m dto.RealtimeMetrics
	var traceDate time.Time
	err := r.DB.QueryRow(query, userID).Scan(&m.CPUUsage, &m.MemoryUsage, &m.DiskUsage, &m.NetworkUsage, &traceDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	m.TraceDate = toLocalTime(traceDate).Format(time.RFC3339)
	return &m, nil
}

func (r *MetricsDashboardRepository) GetMetricsLastHour(userID int) ([]dto.RealtimeMetrics, error) {
	query := `
		SELECT cpu_usage, memory_usage, disk_usage, network_usage, trace_date
		FROM trace_results
		WHERE user_id = $1
		AND trace_date >= NOW() - INTERVAL '1 hour'
		ORDER BY trace_date
	`
	return r.queryMetricsList(query, userID)
}

func (r *MetricsDashboardRepository) GetMetricsLast24Hours(userID int) ([]dto.RealtimeMetrics, error) {
	query := `
		SELECT cpu_usage, memory_usage, disk_usage, network_usage, trace_date
		FROM trace_results
		WHERE user_id = $1
		AND trace_date >= NOW() - INTERVAL '24 hours'
		ORDER BY trace_date
	`
	return r.queryMetricsList(query, userID)
}

func (r *MetricsDashboardRepository) queryMetricsList(query string, userID int) ([]dto.RealtimeMetrics, error) {
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []dto.RealtimeMetrics
	for rows.Next() {
		var m dto.RealtimeMetrics
		var traceDate time.Time
		err := rows.Scan(&m.CPUUsage, &m.MemoryUsage, &m.DiskUsage, &m.NetworkUsage, &traceDate)
		if err != nil {
			return nil, err
		}
		m.TraceDate = toLocalTime(traceDate).Format(time.RFC3339)
		results = append(results, m)
	}
	return results, nil
}
