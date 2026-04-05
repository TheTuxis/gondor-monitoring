package service

import (
	"go.uber.org/zap"

	"github.com/TheTuxis/gondor-monitoring/internal/model"
	"github.com/TheTuxis/gondor-monitoring/internal/repository"
)

type AuditLogService struct {
	auditLogRepo *repository.AuditLogRepository
	logger       *zap.Logger
}

func NewAuditLogService(auditLogRepo *repository.AuditLogRepository, logger *zap.Logger) *AuditLogService {
	return &AuditLogService{auditLogRepo: auditLogRepo, logger: logger}
}

func (s *AuditLogService) List(params model.ListParams) (*model.PaginatedResult, error) {
	logs, total, err := s.auditLogRepo.List(params)
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
		Data: logs,
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

func (s *AuditLogService) Create(input model.AuditLogCreate) (*model.AuditLog, error) {
	log := &model.AuditLog{
		CompanyID:    input.CompanyID,
		UserID:       input.UserID,
		Action:       input.Action,
		ResourceType: input.ResourceType,
		ResourceID:   input.ResourceID,
		Details:      input.Details,
		IPAddress:    input.IPAddress,
	}

	if err := s.auditLogRepo.Create(log); err != nil {
		return nil, err
	}

	return log, nil
}
