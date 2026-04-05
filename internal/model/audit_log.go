package model

import (
	"time"
)

type AuditLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	CompanyID    uint      `gorm:"index;not null" json:"company_id"`
	UserID       uint      `gorm:"index;not null" json:"user_id"`
	Action       string    `gorm:"not null;index" json:"action"`
	ResourceType string    `gorm:"not null;index" json:"resource_type"`
	ResourceID   string    `gorm:"not null" json:"resource_id"`
	Details      string    `gorm:"type:text" json:"details"` // JSON
	IPAddress    string    `json:"ip_address"`
	CreatedAt    time.Time `json:"created_at"`
}

type AuditLogCreate struct {
	CompanyID    uint   `json:"company_id" binding:"required"`
	UserID       uint   `json:"user_id" binding:"required"`
	Action       string `json:"action" binding:"required"`
	ResourceType string `json:"resource_type" binding:"required"`
	ResourceID   string `json:"resource_id" binding:"required"`
	Details      string `json:"details"`
	IPAddress    string `json:"ip_address"`
}
