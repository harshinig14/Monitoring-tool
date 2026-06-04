package routes

import (
	"log"

	"github.com/gin-gonic/gin"
	"MONITORING-TOOL/internal/handlers"
)

func RegisterRoutes(
	router *gin.Engine, 
	regHandler *handlers.RegistrationHandler, 
	metHandler *handlers.MetricsHandler, 
	hbHandler *handlers.HeartbeatHandler, 
	devHandler *handlers.DeviceHandler, 
	cfgHandler *handlers.ConfigHandler, 
	thresholdHandler *handlers.ThresholdHandler, 
	alertHandler *handlers.AlertHandler,
	emailHandler *handlers.EmailHandler,
	reportHandler *handlers.ReportHandler,
	wsHandler *handlers.WSHandler,
	lifecycleHandler *handlers.DeviceLifecycleHandler,
) {
	router.GET("/health", handlers.HealthHandler)
	router.GET("/ws", wsHandler.HandleWS)
	
	api := router.Group("/api/v1")
	{
		api.POST("/debug-log", func(c *gin.Context) {
			var payload struct {
				Message string      `json:"message"`
				Data    interface{} `json:"data"`
			}
			if err := c.BindJSON(&payload); err == nil {
				log.Printf("[FRONTEND DEBUG] %s: %+v\n", payload.Message, payload.Data)
				c.JSON(200, gin.H{"status": "ok"})
			} else {
				c.JSON(400, gin.H{"error": err.Error()})
			}
		})

		api.POST("/register", regHandler.Register)
		api.POST("/metrics", metHandler.UploadMetrics)
		api.POST("/heartbeat", hbHandler.ProcessHeartbeat)
		api.GET("/ws", wsHandler.HandleWS)

		// Dashboard APIs
		api.GET("/devices", devHandler.GetDevices)
		api.GET("/devices/:id", devHandler.GetDeviceByID)
		api.GET("/metrics/realtime/:userId", devHandler.GetRealtimeMetrics)
		api.GET("/metrics/hourly/:userId", devHandler.GetHourlyMetrics)
		api.GET("/metrics/daily/:userId", devHandler.GetDailyMetrics)

		// Device Lifecycle APIs (Phase 13)
		api.POST("/devices/request", lifecycleHandler.RequestAddDevice)
		api.POST("/devices/:id/deactivate", lifecycleHandler.DeactivateDevice)
		api.POST("/devices/:id/activate", lifecycleHandler.ActivateDevice)
		api.DELETE("/devices/:id", lifecycleHandler.RemoveDevice)
		api.GET("/devices/:id/history", lifecycleHandler.GetDeviceHistory)

		// Configuration APIs
		api.GET("/configuration", cfgHandler.GetConfiguration)
		api.PUT("/configuration", cfgHandler.UpdateConfiguration)

		// Threshold APIs (New & Backward compatible)
		api.GET("/thresholds", thresholdHandler.GetThresholds)
		api.PUT("/thresholds", thresholdHandler.UpdateThresholds)
		api.GET("/alerts/thresholds", thresholdHandler.GetThresholds)
		api.PUT("/alerts/thresholds", thresholdHandler.UpdateThresholds)

		// Alerts APIs
		api.GET("/alerts", alertHandler.GetAlerts)
		api.GET("/alerts/recent", alertHandler.GetRecentAlerts)

		// Message Configuration SMTP APIs (Both versions for backward compatibility)
		api.GET("/alerts/email-config", emailHandler.GetEmailConfig)
		api.PUT("/alerts/email-config", emailHandler.UpdateEmailConfig)
		api.GET("/email-config", emailHandler.GetEmailConfig)
		api.PUT("/email-config", emailHandler.UpdateEmailConfig)
		api.POST("/email/test", emailHandler.SendTestEmail)
		api.GET("/email/logs", emailHandler.GetEmailLogs)

		// Download Agent APIs
		api.GET("/download-agent/windows", handlers.DownloadAgentWindows)
		api.GET("/download-agent/linux", handlers.DownloadAgentLinux)
		api.GET("/download-agent/mac", handlers.DownloadAgentMac)

		// Reports APIs
		api.GET("/reports/daily", reportHandler.GetDailyReport)
		api.GET("/reports/weekly", reportHandler.GetWeeklyReport)
		api.GET("/reports/monthly", reportHandler.GetMonthlyReport)
		api.GET("/reports/history", reportHandler.GetDailyHistory)
		api.GET("/reports/export/csv", reportHandler.ExportCSV)
	}
}
