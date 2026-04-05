package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/TheTuxis/gondor-monitoring/internal/model"
	"github.com/TheTuxis/gondor-monitoring/internal/service"
)

type ServiceStatusHandler struct {
	serviceStatusService *service.ServiceStatusService
}

func NewServiceStatusHandler(serviceStatusService *service.ServiceStatusService) *ServiceStatusHandler {
	return &ServiceStatusHandler{serviceStatusService: serviceStatusService}
}

func (h *ServiceStatusHandler) List(c *gin.Context) {
	statuses, err := h.serviceStatusService.ListAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "failed to list service statuses",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": statuses})
}

func (h *ServiceStatusHandler) UpdateStatus(c *gin.Context) {
	var input model.ServiceStatusUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "invalid request body",
			"details": gin.H{"error": err.Error()},
		})
		return
	}

	status, err := h.serviceStatusService.UpdateStatus(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "failed to update service status",
		})
		return
	}

	c.JSON(http.StatusOK, status)
}
