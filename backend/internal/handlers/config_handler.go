package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"MONITORING-TOOL/internal/repository"
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

// GET /api/v1/alerts/thresholds
func (h *ConfigHandler) GetThresholds(c *gin.Context) {
	t, err := h.configService.GetAlertThresholds()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, t)
}

// PUT /api/v1/alerts/thresholds
func (h *ConfigHandler) UpdateThresholds(c *gin.Context) {
	var t repository.AlertThresholds
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.configService.SaveAlertThresholds(&t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// GET /api/v1/alerts/email-config
func (h *ConfigHandler) GetEmailConfig(c *gin.Context) {
	config, err := h.configService.GetEmailConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, config)
}

// PUT /api/v1/alerts/email-config
func (h *ConfigHandler) UpdateEmailConfig(c *gin.Context) {
	var config repository.EmailConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.configService.SaveEmailConfig(&config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
