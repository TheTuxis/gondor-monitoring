package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/TheTuxis/gondor-monitoring/internal/model"
	"github.com/TheTuxis/gondor-monitoring/internal/service"
)

type AlertRuleHandler struct {
	alertRuleService *service.AlertRuleService
}

func NewAlertRuleHandler(alertRuleService *service.AlertRuleService) *AlertRuleHandler {
	return &AlertRuleHandler{alertRuleService: alertRuleService}
}

func (h *AlertRuleHandler) List(c *gin.Context) {
	var params model.ListParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "invalid query parameters",
			"details": gin.H{"error": err.Error()},
		})
		return
	}

	result, err := h.alertRuleService.List(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "failed to list alert rules",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *AlertRuleHandler) GetByID(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	rule, err := h.alertRuleService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "not_found",
			"message": "alert rule not found",
		})
		return
	}

	c.JSON(http.StatusOK, rule)
}

func (h *AlertRuleHandler) Create(c *gin.Context) {
	var input model.AlertRuleCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "invalid request body",
			"details": gin.H{"error": err.Error()},
		})
		return
	}

	rule, err := h.alertRuleService.Create(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "failed to create alert rule",
		})
		return
	}

	c.JSON(http.StatusCreated, rule)
}

func (h *AlertRuleHandler) Update(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	var input model.AlertRuleUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "invalid request body",
			"details": gin.H{"error": err.Error()},
		})
		return
	}

	rule, err := h.alertRuleService.Update(id, input)
	if err != nil {
		if err == service.ErrAlertRuleNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "alert rule not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "failed to update alert rule",
		})
		return
	}

	c.JSON(http.StatusOK, rule)
}

func (h *AlertRuleHandler) Delete(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	if err := h.alertRuleService.Delete(id); err != nil {
		if err == service.ErrAlertRuleNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "alert rule not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "failed to delete alert rule",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "alert rule deleted successfully"})
}
