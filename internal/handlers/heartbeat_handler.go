package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"MONITORING-TOOL/internal/models"
	"MONITORING-TOOL/internal/service"
)

type HeartbeatHandler struct {
	heartbeatService *service.HeartbeatService
}

func NewHeartbeatHandler(svc *service.HeartbeatService) *HeartbeatHandler {
	return &HeartbeatHandler{heartbeatService: svc}
}

func (h *HeartbeatHandler) ProcessHeartbeat(c *gin.Context) {
	var req models.HeartbeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.heartbeatService.ProcessHeartbeat(req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.HeartbeatResponse{Success: true})
}
