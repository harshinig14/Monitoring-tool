package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"MONITORING-TOOL/internal/models"
	"MONITORING-TOOL/internal/service"
)

type RegistrationHandler struct {
	regService *service.RegistrationService
}

func NewRegistrationHandler(regService *service.RegistrationService) *RegistrationHandler {
	return &RegistrationHandler{regService: regService}
}

func (h *RegistrationHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.regService.RegisterDevice(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := models.RegisterResponse{UserID: userID}
	c.JSON(http.StatusOK, res)
}
