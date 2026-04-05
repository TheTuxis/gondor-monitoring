package service

import (
	"errors"

	"go.uber.org/zap"

	"github.com/TheTuxis/gondor-monitoring/internal/model"
	"github.com/TheTuxis/gondor-monitoring/internal/repository"
)

var (
	ErrAlertRuleNotFound = errors.New("alert rule not found")
)

type AlertRuleService struct {
	alertRuleRepo *repository.AlertRuleRepository
	logger        *zap.Logger
}

func NewAlertRuleService(alertRuleRepo *repository.AlertRuleRepository, logger *zap.Logger) *AlertRuleService {
	return &AlertRuleService{alertRuleRepo: alertRuleRepo, logger: logger}
}

func (s *AlertRuleService) List(params model.ListParams) (*model.PaginatedResult, error) {
	rules, total, err := s.alertRuleRepo.List(params)
	if err != nil {
		return nil, err
	}

	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 20
	}

	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize > 0 {
		totalPages++
	}

	return &model.PaginatedResult{
		Data: rules,
		Pagination: model.Pagination{
			Page:       params.Page,
			PageSize:   params.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
			HasNext:    params.Page < totalPages,
			HasPrev:    params.Page > 1,
		},
	}, nil
}

func (s *AlertRuleService) GetByID(id uint) (*model.AlertRule, error) {
	rule, err := s.alertRuleRepo.FindByID(id)
	if err != nil {
		return nil, ErrAlertRuleNotFound
	}
	return rule, nil
}

func (s *AlertRuleService) Create(input model.AlertRuleCreate) (*model.AlertRule, error) {
	durationSeconds := input.DurationSeconds
	if durationSeconds == 0 {
		durationSeconds = 60
	}
	isActive := true
	if input.IsActive != nil {
		isActive = *input.IsActive
	}

	rule := &model.AlertRule{
		CompanyID:       input.CompanyID,
		Name:            input.Name,
		Description:     input.Description,
		MetricName:      input.MetricName,
		Condition:       input.Condition,
		Threshold:       input.Threshold,
		DurationSeconds: durationSeconds,
		Severity:        input.Severity,
		IsActive:        isActive,
		NotifyChannels:  input.NotifyChannels,
	}

	if err := s.alertRuleRepo.Create(rule); err != nil {
		return nil, err
	}

	return rule, nil
}

func (s *AlertRuleService) Update(id uint, input model.AlertRuleUpdate) (*model.AlertRule, error) {
	rule, err := s.alertRuleRepo.FindByID(id)
	if err != nil {
		return nil, ErrAlertRuleNotFound
	}

	if input.Name != nil {
		rule.Name = *input.Name
	}
	if input.Description != nil {
		rule.Description = *input.Description
	}
	if input.MetricName != nil {
		rule.MetricName = *input.MetricName
	}
	if input.Condition != nil {
		rule.Condition = *input.Condition
	}
	if input.Threshold != nil {
		rule.Threshold = *input.Threshold
	}
	if input.DurationSeconds != nil {
		rule.DurationSeconds = *input.DurationSeconds
	}
	if input.Severity != nil {
		rule.Severity = *input.Severity
	}
	if input.IsActive != nil {
		rule.IsActive = *input.IsActive
	}
	if input.NotifyChannels != nil {
		rule.NotifyChannels = *input.NotifyChannels
	}

	if err := s.alertRuleRepo.Update(rule); err != nil {
		return nil, err
	}

	return rule, nil
}

func (s *AlertRuleService) Delete(id uint) error {
	if _, err := s.alertRuleRepo.FindByID(id); err != nil {
		return ErrAlertRuleNotFound
	}
	return s.alertRuleRepo.Delete(id)
}
