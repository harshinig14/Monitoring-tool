package service

import (
	"MONITORING-TOOL/internal/dto"
	"MONITORING-TOOL/internal/repository"
)

type DeviceService struct {
	deviceRepo   *repository.DeviceRepository
	dashboardRepo *repository.MetricsDashboardRepository
}

func NewDeviceService(deviceRepo *repository.DeviceRepository, dashboardRepo *repository.MetricsDashboardRepository) *DeviceService {
	return &DeviceService{deviceRepo: deviceRepo, dashboardRepo: dashboardRepo}
}

func (s *DeviceService) GetDevices() ([]dto.DeviceResponse, error) {
	return s.deviceRepo.GetAllDevices()
}

func (s *DeviceService) GetDeviceByID(id int) (*dto.DeviceResponse, error) {
	return s.deviceRepo.GetDeviceByID(id)
}

func (s *DeviceService) GetLatestMetrics(userID int) (*dto.RealtimeMetrics, error) {
	return s.dashboardRepo.GetLatestMetrics(userID)
}

func (s *DeviceService) GetMetricsLastHour(userID int) ([]dto.RealtimeMetrics, error) {
	return s.dashboardRepo.GetMetricsLastHour(userID)
}

func (s *DeviceService) GetMetricsLast24Hours(userID int) ([]dto.RealtimeMetrics, error) {
	return s.dashboardRepo.GetMetricsLast24Hours(userID)
}
