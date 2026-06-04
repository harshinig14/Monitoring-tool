package repository

import (
	"database/sql"
	"time"
)

type User struct {
	UserID      int
	Username    string
	MachineName string
	OSType      string
	LastSeen    time.Time
	CreatedAt   time.Time
}

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) GetUserByMachineName(machineName string) (*User, error) {
	var user User
	query := `SELECT user_id, username, machine_name, os_type, last_seen, created_at FROM users WHERE machine_name = $1`
	err := r.DB.QueryRow(query, machineName).Scan(
		&user.UserID,
		&user.Username,
		&user.MachineName,
		&user.OSType,
		&user.LastSeen,
		&user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No user found
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) CreateUser(username string, machineName string, osType string) (int, error) {
	var userID int
	query := `INSERT INTO users (username, machine_name, os_type, last_seen) VALUES ($1, $2, $3, NOW()) RETURNING user_id`
	err := r.DB.QueryRow(query, username, machineName, osType).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (r *UserRepository) GetUserByID(userID int) (*User, error) {
	var user User
	var lastSeen sql.NullTime
	query := `SELECT user_id, username, machine_name, os_type, last_seen, created_at FROM users WHERE user_id = $1`
	err := r.DB.QueryRow(query, userID).Scan(
		&user.UserID,
		&user.Username,
		&user.MachineName,
		&user.OSType,
		&lastSeen,
		&user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if lastSeen.Valid {
		user.LastSeen = lastSeen.Time
	}
	return &user, nil
}
