package repository

import (
	"github.com/TheTuxis/gondor-monitoring/internal/model"
	"gorm.io/gorm"
)

type ServiceStatusRepository struct {
	db *gorm.DB
}

func NewServiceStatusRepository(db *gorm.DB) *ServiceStatusRepository {
	return &ServiceStatusRepository{db: db}
}

func (r *ServiceStatusRepository) ListAll() ([]model.ServiceStatus, error) {
	var statuses []model.ServiceStatus
	err := r.db.Order("service_name asc").Find(&statuses).Error
	return statuses, err
}

func (r *ServiceStatusRepository) FindByServiceName(name string) (*model.ServiceStatus, error) {
	var status model.ServiceStatus
	if err := r.db.Where("service_name = ?", name).First(&status).Error; err != nil {
		return nil, err
	}
	return &status, nil
}

func (r *ServiceStatusRepository) Upsert(status *model.ServiceStatus) error {
	return r.db.Where("service_name = ?", status.ServiceName).
		Assign(model.ServiceStatus{
			Status:       status.Status,
			LastCheckAt:  status.LastCheckAt,
			LatencyMs:    status.LatencyMs,
			ErrorMessage: status.ErrorMessage,
		}).
		FirstOrCreate(status).Error
}
