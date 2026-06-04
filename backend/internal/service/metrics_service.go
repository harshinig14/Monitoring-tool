package service

import (
	"MONITORING-TOOL/internal/alerts"
	"MONITORING-TOOL/internal/models"
	"MONITORING-TOOL/internal/repository"
	"MONITORING-TOOL/internal/websocket"
)

type MetricsService struct {
	metricsRepo   *repository.MetricsRepository
	thresholdRepo *repository.ThresholdRepository
	alertService  *AlertService
	userRepo      *repository.UserRepository
	lifecycleRepo *repository.DeviceLifecycleRepository
	wsHub         *websocket.Hub
}

func NewMetricsService(
	metricsRepo *repository.MetricsRepository,
	thresholdRepo *repository.ThresholdRepository,
	alertService *AlertService,
	userRepo *repository.UserRepository,
	lifecycleRepo *repository.DeviceLifecycleRepository,
	wsHub *websocket.Hub,
) *MetricsService {
	return &MetricsService{
		metricsRepo:   metricsRepo,
		thresholdRepo: thresholdRepo,
		alertService:  alertService,
		userRepo:      userRepo,
		lifecycleRepo: lifecycleRepo,
		wsHub:         wsHub,
	}
}

func (s *MetricsService) SaveMetrics(req models.MetricsRequest) error {
	// 1. Check deactivated/removed status
	status, err := s.lifecycleRepo.GetDeviceStatus(req.UserID)
	if err == nil && (status == "DEACTIVATED" || status == "REMOVED") {
		// Ignore metrics telemetry if deactivated or removed
		return nil
	}

	// 2. Save raw metrics to database
	err = s.metricsRepo.InsertMetrics(req)
	if err != nil {
		return err
	}

	// 3. Resolve machine name from UserID
	user, err := s.userRepo.GetUserByID(req.UserID)
	if err != nil || user == nil {
		return nil // Ignore if user details cannot be resolved
	}

	// 4. Load alert thresholds from DB
	t, err := s.thresholdRepo.GetThresholds()
	if err != nil {
		return nil
	}

	// 5. Evaluate metrics
	candidateAlerts := alerts.EvaluateMetrics(req, user.MachineName, t)

	// 6. Delegate alert creation to AlertService (handles spam cooldown & emails)
	for _, a := range candidateAlerts {
		_ = s.alertService.CreateAlert(&a)
	}

	// 7. Broadcast metrics to WebSocket clients
	websocket.BroadcastMetric(s.wsHub, req.UserID, req.CPUUsage, req.MemoryUsage, req.DiskUsage, req.NetworkUsage)

	return nil
}
