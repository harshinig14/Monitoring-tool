package service

import (
	"MONITORING-TOOL/internal/models"
	"MONITORING-TOOL/internal/repository"
)

type MetricsService struct {
	metricsRepo *repository.MetricsRepository
}

func NewMetricsService(repo *repository.MetricsRepository) *MetricsService {
	return &MetricsService{metricsRepo: repo}
}

func (s *MetricsService) SaveMetrics(req models.MetricsRequest) error {
	// In the future: Add logic to validate UserID existence or thresholds
	return s.metricsRepo.InsertMetrics(req)
}
