package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"MONITORING-TOOL/internal/service"
)

type DeviceHandler struct {
	deviceService *service.DeviceService
}

func NewDeviceHandler(svc *service.DeviceService) *DeviceHandler {
	return &DeviceHandler{deviceService: svc}
}

// GET /api/v1/devices
func (h *DeviceHandler) GetDevices(c *gin.Context) {
	devices, err := h.deviceService.GetDevices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, devices)
}

// GET /api/v1/devices/:id
func (h *DeviceHandler) GetDeviceByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}

	device, err := h.deviceService.GetDeviceByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if device == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
		return
	}
	c.JSON(http.StatusOK, device)
}

// GET /api/v1/metrics/realtime/:userId
func (h *DeviceHandler) GetRealtimeMetrics(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	metrics, err := h.deviceService.GetLatestMetrics(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if metrics == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no metrics found"})
		return
	}
	c.JSON(http.StatusOK, metrics)
}

// GET /api/v1/metrics/hourly/:userId
func (h *DeviceHandler) GetHourlyMetrics(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	metrics, err := h.deviceService.GetMetricsLastHour(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, metrics)
}

// GET /api/v1/metrics/daily/:userId
func (h *DeviceHandler) GetDailyMetrics(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	metrics, err := h.deviceService.GetMetricsLast24Hours(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, metrics)
}
