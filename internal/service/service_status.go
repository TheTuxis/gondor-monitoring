package service

import (
	"time"

	"go.uber.org/zap"

	"github.com/TheTuxis/gondor-monitoring/internal/model"
	"github.com/TheTuxis/gondor-monitoring/internal/repository"
)

type ServiceStatusService struct {
	serviceStatusRepo *repository.ServiceStatusRepository
	logger            *zap.Logger
}

func NewServiceStatusService(serviceStatusRepo *repository.ServiceStatusRepository, logger *zap.Logger) *ServiceStatusService {
	return &ServiceStatusService{serviceStatusRepo: serviceStatusRepo, logger: logger}
}

func (s *ServiceStatusService) ListAll() ([]model.ServiceStatus, error) {
	return s.serviceStatusRepo.ListAll()
}

func (s *ServiceStatusService) UpdateStatus(input model.ServiceStatusUpdate) (*model.ServiceStatus, error) {
	status := &model.ServiceStatus{
		ServiceName:  input.ServiceName,
		Status:       input.Status,
		LastCheckAt:  time.Now(),
		LatencyMs:    input.LatencyMs,
		ErrorMessage: input.ErrorMessage,
	}

	if err := s.serviceStatusRepo.Upsert(status); err != nil {
		return nil, err
	}

	return status, nil
}
