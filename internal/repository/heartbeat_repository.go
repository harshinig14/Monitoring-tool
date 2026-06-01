package repository

import (
	"database/sql"
)

type HeartbeatRepository struct {
	DB *sql.DB
}

func NewHeartbeatRepository(db *sql.DB) *HeartbeatRepository {
	return &HeartbeatRepository{DB: db}
}

func (r *HeartbeatRepository) UpdateHeartbeat(userID int) error {
	query := `UPDATE users SET last_seen = NOW() WHERE user_id = $1`
	_, err := r.DB.Exec(query, userID)
	return err
}
