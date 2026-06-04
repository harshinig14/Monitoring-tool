package service

import (
	"MONITORING-TOOL/internal/repository"
	"MONITORING-TOOL/internal/websocket"
)

type HeartbeatService struct {
	heartbeatRepo *repository.HeartbeatRepository
	userRepo      *repository.UserRepository
	lifecycleRepo *repository.DeviceLifecycleRepository
	wsHub         *websocket.Hub
}

func NewHeartbeatService(
	hbRepo *repository.HeartbeatRepository, 
	userRepo *repository.UserRepository,
	lifecycleRepo *repository.DeviceLifecycleRepository,
	wsHub *websocket.Hub,
) *HeartbeatService {
	return &HeartbeatService{
		heartbeatRepo: hbRepo,
		userRepo:      userRepo,
		lifecycleRepo: lifecycleRepo,
		wsHub:         wsHub,
	}
}

func (s *HeartbeatService) ProcessHeartbeat(userID int) error {
	// 1. Get current status to inspect deactivated/removed ignore rules
	status, err := s.lifecycleRepo.GetDeviceStatus(userID)
	if err != nil {
		return err
	}

	if status == "DEACTIVATED" || status == "REMOVED" {
		// Ignore telemetry data silently for deactivated or removed nodes
		return nil
	}

	// 2. Update database heartbeat last_seen
	err = s.heartbeatRepo.UpdateHeartbeat(userID)
	if err != nil {
		return err
	}

	// 3. Transition to ONLINE if previously offline or pending
	if status == "OFFLINE" || status == "PENDING_REGISTRATION" {
		err = s.lifecycleRepo.UpdateDeviceStatus(userID, "ONLINE")
		if err == nil {
			_ = s.lifecycleRepo.InsertStatusHistory(userID, status, "ONLINE")
			user, err := s.userRepo.GetUserByID(userID)
			if err == nil && user != nil {
				websocket.BroadcastDeviceStatus(s.wsHub, userID, user.MachineName, "ONLINE")
			}
		}
	}

	return nil
}
