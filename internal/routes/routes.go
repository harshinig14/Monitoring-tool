package routes

import (
	"github.com/gin-gonic/gin"
	"MONITORING-TOOL/internal/handlers"
)

func RegisterRoutes(router *gin.Engine, regHandler *handlers.RegistrationHandler, metHandler *handlers.MetricsHandler, hbHandler *handlers.HeartbeatHandler, devHandler *handlers.DeviceHandler) {
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
	}
}
