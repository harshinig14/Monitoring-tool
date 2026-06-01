package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"MONITORING-TOOL/internal/models"
	"MONITORING-TOOL/internal/service"
)

type MetricsHandler struct {
	metricsService *service.MetricsService
}

func NewMetricsHandler(svc *service.MetricsService) *MetricsHandler {
	return &MetricsHandler{metricsService: svc}
}

func (h *MetricsHandler) UploadMetrics(c *gin.Context) {
	var req models.MetricsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.metricsService.SaveMetrics(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.MetricsResponse{Success: true})
}
