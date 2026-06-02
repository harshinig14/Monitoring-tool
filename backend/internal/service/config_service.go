package service

import (
	"MONITORING-TOOL/internal/repository"
)

type ConfigService struct {
	configRepo      *repository.ConfigRepository
	alertRepo       *repository.AlertRepository
	emailConfigRepo *repository.EmailConfigRepository
}

func NewConfigService(cr *repository.ConfigRepository, ar *repository.AlertRepository, er *repository.EmailConfigRepository) *ConfigService {
	return &ConfigService{
		configRepo:      cr,
		alertRepo:       ar,
		emailConfigRepo: er,
	}
}

// 1. Polling Configuration
func (s *ConfigService) GetPollingFrequency() (int, error) {
	return s.configRepo.GetConfiguration()
}

func (s *ConfigService) SavePollingFrequency(freq int) error {
	return s.configRepo.SaveConfiguration(freq)
}

// 2. Alert Thresholds
func (s *ConfigService) GetAlertThresholds() (*repository.AlertThresholds, error) {
	return s.alertRepo.GetThresholds()
}

func (s *ConfigService) SaveAlertThresholds(t *repository.AlertThresholds) error {
	return s.alertRepo.SaveThresholds(t)
}

// 3. Email Configuration
func (s *ConfigService) GetEmailConfig() (*repository.EmailConfig, error) {
	return s.emailConfigRepo.GetConfig()
}

func (s *ConfigService) SaveEmailConfig(c *repository.EmailConfig) error {
	return s.emailConfigRepo.SaveConfig(c)
}
