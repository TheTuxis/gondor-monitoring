package service

import (
	"errors"
	"time"

	"go.uber.org/zap"

	"github.com/TheTuxis/gondor-monitoring/internal/model"
	"github.com/TheTuxis/gondor-monitoring/internal/repository"
)

var (
	ErrAlertNotFound = errors.New("alert not found")
)

type AlertService struct {
	alertRepo *repository.AlertRepository
	logger    *zap.Logger
}

func NewAlertService(alertRepo *repository.AlertRepository, logger *zap.Logger) *AlertService {
	return &AlertService{alertRepo: alertRepo, logger: logger}
}

func (s *AlertService) List(params model.ListParams) (*model.PaginatedResult, error) {
	alerts, total, err := s.alertRepo.List(params)
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
		Data: alerts,
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

func (s *AlertService) Acknowledge(id uint, userID uint) (*model.Alert, error) {
	alert, err := s.alertRepo.FindByID(id)
	if err != nil {
		return nil, ErrAlertNotFound
	}

	alert.AcknowledgedBy = &userID
	now := time.Now()
	if alert.Status == "firing" {
		alert.Status = "resolved"
		alert.ResolvedAt = &now
	}

	if err := s.alertRepo.Update(alert); err != nil {
		return nil, err
	}

	return alert, nil
}
