package model

import (
	"time"
)

type ServiceStatus struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ServiceName  string    `gorm:"uniqueIndex;not null" json:"service_name"`
	Status       string    `gorm:"not null;default:healthy" json:"status"` // healthy, degraded, unhealthy
	LastCheckAt  time.Time `gorm:"not null" json:"last_check_at"`
	LatencyMs    float64   `json:"latency_ms"`
	ErrorMessage *string   `json:"error_message"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ServiceStatusUpdate struct {
	ServiceName  string  `json:"service_name" binding:"required"`
	Status       string  `json:"status" binding:"required,oneof=healthy degraded unhealthy"`
	LatencyMs    float64 `json:"latency_ms"`
	ErrorMessage *string `json:"error_message"`
}
