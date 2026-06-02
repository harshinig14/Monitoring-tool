package service

import (
	"MONITORING-TOOL/internal/repository"
)

type HeartbeatService struct {
	heartbeatRepo *repository.HeartbeatRepository
}

func NewHeartbeatService(repo *repository.HeartbeatRepository) *HeartbeatService {
	return &HeartbeatService{heartbeatRepo: repo}
}

func (s *HeartbeatService) ProcessHeartbeat(userID int) error {
	return s.heartbeatRepo.UpdateHeartbeat(userID)
}
