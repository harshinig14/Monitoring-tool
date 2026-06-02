package routes

import (
	"github.com/gin-gonic/gin"
	"MONITORING-TOOL/internal/handlers"
)

func RegisterRoutes(router *gin.Engine, regHandler *handlers.RegistrationHandler, metHandler *handlers.MetricsHandler, hbHandler *handlers.HeartbeatHandler, devHandler *handlers.DeviceHandler, cfgHandler *handlers.ConfigHandler) {
	router.GET("/health", handlers.HealthHandler)
	
	api := router.Group("/api/v1")
	{
		api.POST("/register", regHandler.Register)
		api.POST("/metrics", metHandler.UploadMetrics)
		api.POST("/heartbeat", hbHandler.ProcessHeartbeat)

		// Dashboard APIs
		api.GET("/devices", devHandler.GetDevices)
		api.GET("/devices/:id", devHandler.GetDeviceByID)
		api.GET("/metrics/realtime/:userId", devHandler.GetRealtimeMetrics)
		api.GET("/metrics/hourly/:userId", devHandler.GetHourlyMetrics)
		api.GET("/metrics/daily/:userId", devHandler.GetDailyMetrics)

		// Configuration APIs
		api.GET("/configuration", cfgHandler.GetConfiguration)
		api.PUT("/configuration", cfgHandler.UpdateConfiguration)

		// Alerts Settings APIs
		api.GET("/alerts/thresholds", cfgHandler.GetThresholds)
		api.PUT("/alerts/thresholds", cfgHandler.UpdateThresholds)

		// Message Configuration SMTP APIs
		api.GET("/alerts/email-config", cfgHandler.GetEmailConfig)
		api.PUT("/alerts/email-config", cfgHandler.UpdateEmailConfig)

		// Download Agent APIs
		api.GET("/download-agent/windows", handlers.DownloadAgentWindows)
		api.GET("/download-agent/linux", handlers.DownloadAgentLinux)
		api.GET("/download-agent/mac", handlers.DownloadAgentMac)

		// Reports APIs
		api.GET("/reports/daily", handlers.GetDailyReports)
		api.GET("/reports/weekly", handlers.GetWeeklyReports)
		api.GET("/reports/monthly", handlers.GetMonthlyReports)
	}
}
