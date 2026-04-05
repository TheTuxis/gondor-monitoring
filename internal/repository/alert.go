package repository

import (
	"github.com/TheTuxis/gondor-monitoring/internal/model"
	"gorm.io/gorm"
)

type AlertRepository struct {
	db *gorm.DB
}

func NewAlertRepository(db *gorm.DB) *AlertRepository {
	return &AlertRepository{db: db}
}

func (r *AlertRepository) FindByID(id uint) (*model.Alert, error) {
	var alert model.Alert
	if err := r.db.Preload("Rule").First(&alert, id).Error; err != nil {
		return nil, err
	}
	return &alert, nil
}

func (r *AlertRepository) List(params model.ListParams) ([]model.Alert, int64, error) {
	var alerts []model.Alert
	var total int64

	query := r.db.Model(&model.Alert{})

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.Severity != "" {
		query = query.Joins("JOIN alert_rules ON alert_rules.id = alerts.rule_id").
			Where("alert_rules.severity = ?", params.Severity)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 20
	}

	sortBy := "id"
	if params.SortBy != "" {
		sortBy = params.SortBy
	}
	sortOrder := "desc"
	if params.SortOrder == "asc" {
		sortOrder = "asc"
	}

	offset := (params.Page - 1) * params.PageSize
	err := r.db.Preload("Rule").
		Order(sortBy + " " + sortOrder).
		Offset(offset).Limit(params.PageSize).
		Find(&alerts).Error

	return alerts, total, err
}

func (r *AlertRepository) Update(alert *model.Alert) error {
	return r.db.Save(alert).Error
}
