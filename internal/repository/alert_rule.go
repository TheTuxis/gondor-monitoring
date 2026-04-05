package repository

import (
	"github.com/TheTuxis/gondor-monitoring/internal/model"
	"gorm.io/gorm"
)

type AlertRuleRepository struct {
	db *gorm.DB
}

func NewAlertRuleRepository(db *gorm.DB) *AlertRuleRepository {
	return &AlertRuleRepository{db: db}
}

func (r *AlertRuleRepository) FindByID(id uint) (*model.AlertRule, error) {
	var rule model.AlertRule
	if err := r.db.First(&rule, id).Error; err != nil {
		return nil, err
	}
	return &rule, nil
}

func (r *AlertRuleRepository) List(params model.ListParams) ([]model.AlertRule, int64, error) {
	var rules []model.AlertRule
	var total int64

	query := r.db.Model(&model.AlertRule{})

	if params.Search != "" {
		search := "%" + params.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", search, search)
	}
	if params.CompanyID != nil {
		query = query.Where("company_id = ?", *params.CompanyID)
	}
	if params.Severity != "" {
		query = query.Where("severity = ?", params.Severity)
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
	sortOrder := "asc"
	if params.SortOrder == "desc" {
		sortOrder = "desc"
	}

	offset := (params.Page - 1) * params.PageSize
	err := query.
		Order(sortBy + " " + sortOrder).
		Offset(offset).Limit(params.PageSize).
		Find(&rules).Error

	return rules, total, err
}

func (r *AlertRuleRepository) Create(rule *model.AlertRule) error {
	return r.db.Create(rule).Error
}

func (r *AlertRuleRepository) Update(rule *model.AlertRule) error {
	return r.db.Save(rule).Error
}

func (r *AlertRuleRepository) Delete(id uint) error {
	return r.db.Delete(&model.AlertRule{}, id).Error
}
