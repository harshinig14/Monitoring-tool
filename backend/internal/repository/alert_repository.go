package repository

import (
	"database/sql"
	"MONITORING-TOOL/internal/models"
)

type AlertRepository struct {
	DB *sql.DB
}

func NewAlertRepository(db *sql.DB) *AlertRepository {
	return &AlertRepository{DB: db}
}

func (r *AlertRepository) CreateAlert(a *models.Alert) error {
	query := `
		INSERT INTO alerts (user_id, machine_name, metric_name, current_value, threshold_value, severity)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.DB.Exec(query, a.UserID, a.MachineName, a.MetricName, a.CurrentValue, a.ThresholdValue, a.Severity)
	return err
}

func (r *AlertRepository) GetAlerts() ([]models.Alert, error) {
	query := `
		SELECT alert_id, user_id, machine_name, metric_name, current_value, threshold_value, severity, created_at
		FROM alerts
		ORDER BY created_at DESC
	`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []models.Alert
	for rows.Next() {
		var a models.Alert
		err := rows.Scan(&a.AlertID, &a.UserID, &a.MachineName, &a.MetricName, &a.CurrentValue, &a.ThresholdValue, &a.Severity, &a.CreatedAt)
		if err != nil {
			return nil, err
		}
		alerts = append(alerts, a)
	}
	return alerts, nil
}

func (r *AlertRepository) GetRecentAlerts() ([]models.Alert, error) {
	query := `
		SELECT alert_id, user_id, machine_name, metric_name, current_value, threshold_value, severity, created_at
		FROM alerts
		ORDER BY created_at DESC
		LIMIT 50
	`
	rows, err := r.DB.Query(query)
	if err != nil {
		// Return empty slice instead of error if table is missing or doesn't exist
		return []models.Alert{}, nil
	}
	defer rows.Close()

	var alerts []models.Alert
	for rows.Next() {
		var a models.Alert
		err := rows.Scan(&a.AlertID, &a.UserID, &a.MachineName, &a.MetricName, &a.CurrentValue, &a.ThresholdValue, &a.Severity, &a.CreatedAt)
		if err != nil {
			return nil, err
		}
		alerts = append(alerts, a)
	}
	return alerts, nil
}

func (r *AlertRepository) IsSpamAlert(machineName string, metricName string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM alerts 
			WHERE machine_name = $1 
			AND metric_name = $2 
			AND created_at > NOW() - INTERVAL '5 minutes'
		)
	`
	var exists bool
	err := r.DB.QueryRow(query, machineName, metricName).Scan(&exists)
	if err != nil {
		// If query fails (e.g. missing table), assume not spam
		return false, nil
	}
	return exists, nil
}
