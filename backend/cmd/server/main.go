package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"MONITORING-TOOL/internal/config"
	"MONITORING-TOOL/internal/database"
	"MONITORING-TOOL/internal/handlers"
	"MONITORING-TOOL/internal/jobs"
	"MONITORING-TOOL/internal/middleware"
	"MONITORING-TOOL/internal/repository"
	"MONITORING-TOOL/internal/routes"
	"MONITORING-TOOL/internal/service"
	"MONITORING-TOOL/internal/websocket"
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

	// Initialize WebSocket Hub
	hub := websocket.NewHub()
	go hub.Run()

	// Initialize Repositories
	userRepo := repository.NewUserRepository(db)
	metricsRepo := repository.NewMetricsRepository(db)
	heartbeatRepo := repository.NewHeartbeatRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)
	dashboardRepo := repository.NewMetricsDashboardRepository(db)
	configRepo := repository.NewConfigRepository(db)
	alertRepo := repository.NewAlertRepository(db)
	thresholdRepo := repository.NewThresholdRepository(db)
	emailRepo := repository.NewEmailRepository(db)
	reportRepo := repository.NewReportRepository(db)
	lifecycleRepo := repository.NewDeviceLifecycleRepository(db)

	// Initialize Services
	lifecycleService := service.NewDeviceLifecycleService(lifecycleRepo, hub)
	regService := service.NewRegistrationService(userRepo, lifecycleRepo, hub)
	emailService := service.NewEmailService(emailRepo)
	alertService := service.NewAlertService(alertRepo, emailRepo, hub)
	metService := service.NewMetricsService(metricsRepo, thresholdRepo, alertService, userRepo, lifecycleRepo, hub)
	hbService := service.NewHeartbeatService(heartbeatRepo, userRepo, lifecycleRepo, hub)
	devService := service.NewDeviceService(deviceRepo, dashboardRepo)
	cfgService := service.NewConfigService(configRepo)
	thresholdService := service.NewThresholdService(thresholdRepo)
	reportService := service.NewReportService(reportRepo)

	// Initialize Handlers
	regHandler := handlers.NewRegistrationHandler(regService)
	metHandler := handlers.NewMetricsHandler(metService)
	hbHandler := handlers.NewHeartbeatHandler(hbService)
	devHandler := handlers.NewDeviceHandler(devService)
	cfgHandler := handlers.NewConfigHandler(cfgService)
	thresholdHandler := handlers.NewThresholdHandler(thresholdService)
	alertHandler := handlers.NewAlertHandler(alertService)
	emailHandler := handlers.NewEmailHandler(emailService)
	reportHandler := handlers.NewReportHandler(reportService)
	wsHandler := handlers.NewWSHandler(hub)
	lifecycleHandler := handlers.NewDeviceLifecycleHandler(lifecycleService)

	// Initialize and Start Background Jobs
	statusJob := jobs.NewDeviceStatusJob(db, hub)
	statusJob.Start()

	// Register Routes
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())
	routes.RegisterRoutes(router, regHandler, metHandler, hbHandler, devHandler, cfgHandler, thresholdHandler, alertHandler, emailHandler, reportHandler, wsHandler, lifecycleHandler)

	// Start Gin Server
	log.Printf("Server Started On :%s\n", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
