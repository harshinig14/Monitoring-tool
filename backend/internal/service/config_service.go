package service

import (
	"MONITORING-TOOL/internal/repository"
)

type ConfigService struct {
	configRepo *repository.ConfigRepository
}

func NewConfigService(cr *repository.ConfigRepository) *ConfigService {
	return &ConfigService{
		configRepo: cr,
	}
}

// 1. Polling Configuration
func (s *ConfigService) GetPollingFrequency() (int, error) {
	return s.configRepo.GetConfiguration()
}

func (s *ConfigService) SavePollingFrequency(freq int) error {
	return s.configRepo.SaveConfiguration(freq)
}
