package model

import (
	"time"
)

type Alert struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	RuleID         uint       `gorm:"index;not null" json:"rule_id"`
	Status         string     `gorm:"not null;default:firing;index" json:"status"` // firing, resolved
	Value          float64    `gorm:"not null" json:"value"`
	Message        string     `gorm:"not null" json:"message"`
	FiredAt        time.Time  `gorm:"not null" json:"fired_at"`
	ResolvedAt     *time.Time `json:"resolved_at"`
	AcknowledgedBy *uint      `json:"acknowledged_by"`
	CreatedAt      time.Time  `json:"created_at"`

	Rule AlertRule `gorm:"foreignKey:RuleID" json:"rule,omitempty"`
}
