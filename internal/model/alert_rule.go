package model

import (
	"time"
)

type AlertRule struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	CompanyID        uint      `gorm:"index;not null" json:"company_id"`
	Name             string    `gorm:"not null" json:"name"`
	Description      string    `json:"description"`
	MetricName       string    `gorm:"not null" json:"metric_name"`
	Condition        string    `gorm:"not null" json:"condition"`         // gt, lt, eq, gte, lte
	Threshold        float64   `gorm:"not null" json:"threshold"`
	DurationSeconds  int       `gorm:"not null;default:60" json:"duration_seconds"`
	Severity         string    `gorm:"not null;default:info;index" json:"severity"` // info, warning, critical
	IsActive         bool      `gorm:"not null;default:true" json:"is_active"`
	NotifyChannels   string    `gorm:"type:text" json:"notify_channels"` // JSON array
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	Alerts []Alert `gorm:"foreignKey:RuleID" json:"alerts,omitempty"`
}

type AlertRuleCreate struct {
	CompanyID       uint    `json:"company_id" binding:"required"`
	Name            string  `json:"name" binding:"required"`
	Description     string  `json:"description"`
	MetricName      string  `json:"metric_name" binding:"required"`
	Condition       string  `json:"condition" binding:"required,oneof=gt lt eq gte lte"`
	Threshold       float64 `json:"threshold" binding:"required"`
	DurationSeconds int     `json:"duration_seconds"`
	Severity        string  `json:"severity" binding:"required,oneof=info warning critical"`
	IsActive        *bool   `json:"is_active"`
	NotifyChannels  string  `json:"notify_channels"`
}

type AlertRuleUpdate struct {
	Name            *string  `json:"name"`
	Description     *string  `json:"description"`
	MetricName      *string  `json:"metric_name"`
	Condition       *string  `json:"condition"`
	Threshold       *float64 `json:"threshold"`
	DurationSeconds *int     `json:"duration_seconds"`
	Severity        *string  `json:"severity"`
	IsActive        *bool    `json:"is_active"`
	NotifyChannels  *string  `json:"notify_channels"`
}
