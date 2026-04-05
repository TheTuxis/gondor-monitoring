package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/TheTuxis/gondor-monitoring/internal/model"
	"github.com/TheTuxis/gondor-monitoring/internal/service"
)

type AlertHandler struct {
	alertService *service.AlertService
}

func NewAlertHandler(alertService *service.AlertService) *AlertHandler {
	return &AlertHandler{alertService: alertService}
}

func (h *AlertHandler) List(c *gin.Context) {
	var params model.ListParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "invalid query parameters",
			"details": gin.H{"error": err.Error()},
		})
		return
	}

	result, err := h.alertService.List(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "failed to list alerts",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *AlertHandler) Acknowledge(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	userIDVal, _ := c.Get("user_id")
	userID, _ := userIDVal.(uint)

	alert, err := h.alertService.Acknowledge(id, userID)
	if err != nil {
		if err == service.ErrAlertNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "alert not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "failed to acknowledge alert",
		})
		return
	}

	c.JSON(http.StatusOK, alert)
}
