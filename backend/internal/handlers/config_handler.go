package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"MONITORING-TOOL/internal/service"
)

type ConfigHandler struct {
	configService *service.ConfigService
}

func NewConfigHandler(svc *service.ConfigService) *ConfigHandler {
	return &ConfigHandler{configService: svc}
}

// GET /api/v1/configuration
func (h *ConfigHandler) GetConfiguration(c *gin.Context) {
	freq, err := h.configService.GetPollingFrequency()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"polling_frequency": freq})
}

// PUT /api/v1/configuration
func (h *ConfigHandler) UpdateConfiguration(c *gin.Context) {
	var req struct {
		PollingFrequency int `json:"polling_frequency"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.configService.SavePollingFrequency(req.PollingFrequency)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
