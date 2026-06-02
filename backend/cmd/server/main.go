package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"MONITORING-TOOL/internal/config"
	"MONITORING-TOOL/internal/database"
	"MONITORING-TOOL/internal/handlers"
	"MONITORING-TOOL/internal/repository"
	"MONITORING-TOOL/internal/routes"
	"MONITORING-TOOL/internal/service"
)

func main() {
	// Load Config
	cfg := config.LoadConfig()

	// Initialize Database
	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize Repo, Service, Handler
	userRepo := repository.NewUserRepository(db)
	regService := service.NewRegistrationService(userRepo)
	regHandler := handlers.NewRegistrationHandler(regService)

	metricsRepo := repository.NewMetricsRepository(db)
	metService := service.NewMetricsService(metricsRepo)
	metHandler := handlers.NewMetricsHandler(metService)

	heartbeatRepo := repository.NewHeartbeatRepository(db)
	hbService := service.NewHeartbeatService(heartbeatRepo)
	hbHandler := handlers.NewHeartbeatHandler(hbService)

	deviceRepo := repository.NewDeviceRepository(db)
	dashboardRepo := repository.NewMetricsDashboardRepository(db)
	devService := service.NewDeviceService(deviceRepo, dashboardRepo)
	devHandler := handlers.NewDeviceHandler(devService)

	configRepo := repository.NewConfigRepository(db)
	alertRepo := repository.NewAlertRepository(db)
	emailConfigRepo := repository.NewEmailConfigRepository(db)
	cfgService := service.NewConfigService(configRepo, alertRepo, emailConfigRepo)
	cfgHandler := handlers.NewConfigHandler(cfgService)

	// Register Routes
	router := gin.Default()
	routes.RegisterRoutes(router, regHandler, metHandler, hbHandler, devHandler, cfgHandler)

	// Start Gin Server
	log.Printf("Server Started On :%s\n", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
