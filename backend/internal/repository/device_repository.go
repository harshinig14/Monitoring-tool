package repository

import (
	"database/sql"
	"time"

	"MONITORING-TOOL/internal/dto"
)

type DeviceRepository struct {
	DB *sql.DB
}

func NewDeviceRepository(db *sql.DB) *DeviceRepository {
	return &DeviceRepository{DB: db}
}

func (r *DeviceRepository) GetAllDevices() ([]dto.DeviceResponse, error) {
	query := `
		SELECT
			user_id,
			machine_name,
			os_type,
			last_seen,
			CASE
				WHEN last_seen IS NULL THEN 'PENDING'
				WHEN NOW() - last_seen < INTERVAL '60 seconds' THEN 'ONLINE'
				ELSE 'OFFLINE'
			END AS status
		FROM users
		ORDER BY machine_name
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []dto.DeviceResponse
	for rows.Next() {
		var d dto.DeviceResponse
		var lastSeen sql.NullTime
		err := rows.Scan(&d.UserID, &d.MachineName, &d.OSType, &lastSeen, &d.Status)
		if err != nil {
			return nil, err
		}
		if lastSeen.Valid {
			d.LastSeen = lastSeen.Time.Format(time.RFC3339)
		} else {
			d.LastSeen = ""
		}
		devices = append(devices, d)
	}
	return devices, nil
}

func (r *DeviceRepository) GetDeviceByID(id int) (*dto.DeviceResponse, error) {
	query := `
		SELECT
			user_id,
			machine_name,
			os_type,
			last_seen,
			CASE
				WHEN last_seen IS NULL THEN 'PENDING'
				WHEN NOW() - last_seen < INTERVAL '60 seconds' THEN 'ONLINE'
				ELSE 'OFFLINE'
			END AS status
		FROM users
		WHERE user_id = $1
	`

	var d dto.DeviceResponse
	var lastSeen sql.NullTime
	err := r.DB.QueryRow(query, id).Scan(&d.UserID, &d.MachineName, &d.OSType, &lastSeen, &d.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if lastSeen.Valid {
		d.LastSeen = lastSeen.Time.Format(time.RFC3339)
	} else {
		d.LastSeen = ""
	}
	return &d, nil
}
