package service

import (
	"database/sql"

	"MONITORING-TOOL/internal/models"
	"MONITORING-TOOL/internal/repository"
	"MONITORING-TOOL/internal/websocket"
)

type RegistrationService struct {
	userRepo      *repository.UserRepository
	lifecycleRepo *repository.DeviceLifecycleRepository
	wsHub         *websocket.Hub
}

func NewRegistrationService(
	userRepo *repository.UserRepository,
	lifecycleRepo *repository.DeviceLifecycleRepository,
	wsHub *websocket.Hub,
) *RegistrationService {
	return &RegistrationService{
		userRepo:      userRepo,
		lifecycleRepo: lifecycleRepo,
		wsHub:         wsHub,
	}
}

func (s *RegistrationService) RegisterDevice(req models.RegisterRequest) (int, error) {
	user, err := s.userRepo.GetUserByMachineName(req.MachineName)
	if err != nil {
		return 0, err
	}

	var userID int
	var prevStatus string = "PENDING_REGISTRATION"

	if user != nil {
		userID = user.UserID
		// Retrieve current database status
		status, err := s.lifecycleRepo.GetDeviceStatus(userID)
		if err == nil {
			prevStatus = status
		}
	} else {
		// Attempt to claim a placeholder if one was created during the 'Add Device' download click
		placeholderID, err := s.claimPlaceholder(req.OSType, req.Username, req.MachineName)
		if err == nil && placeholderID != 0 {
			userID = placeholderID
			prevStatus = "PENDING_REGISTRATION"
		} else {
			// Create a brand new device row
			userID, err = s.userRepo.CreateUser(req.Username, req.MachineName, req.OSType)
			if err != nil {
				return 0, err
			}
		}
	}

	// Transition status to ONLINE if it is currently offline or pending registration
	if prevStatus == "PENDING_REGISTRATION" || prevStatus == "OFFLINE" {
		err = s.lifecycleRepo.UpdateDeviceStatus(userID, "ONLINE")
		if err == nil {
			_ = s.lifecycleRepo.InsertStatusHistory(userID, prevStatus, "ONLINE")
			websocket.BroadcastDeviceStatus(s.wsHub, userID, req.MachineName, "ONLINE")
		}
	}

	return userID, nil
}

func (s *RegistrationService) claimPlaceholder(osType, username, machineName string) (int, error) {
	var userID int
	query := `
		SELECT user_id 
		FROM users 
		WHERE device_status = 'PENDING_REGISTRATION' 
		  AND os_type = $1 
		  AND machine_name LIKE 'Pending Device%' 
		ORDER BY created_at ASC 
		LIMIT 1
	`
	err := s.lifecycleRepo.DB.QueryRow(query, osType).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	// Claim and update the user row with actual agent values
	updateQuery := `
		UPDATE users 
		SET machine_name = $1, username = $2, last_seen = NOW()
		WHERE user_id = $3
	`
	_, err = s.lifecycleRepo.DB.Exec(updateQuery, machineName, username, userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
