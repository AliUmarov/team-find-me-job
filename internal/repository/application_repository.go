package repository

import (
	"github.com/AliUmarov/team-find-me-job/internal/models"
	"gorm.io/gorm"
)

type ApplicationRepository interface {
	Create(*models.Application) error
	Applications(uint, models.ApplicationFilter) ([]models.Application, error)
	AcceptApplication(appId uint) error
	RejectApplication(appId uint) error
}

type applicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) ApplicationRepository {
	return &applicationRepository{db: db}
}

func (r *applicationRepository) RejectApplication(appId uint) error {
	if err := r.db.Model(&models.Application{}).
		Where("id = ? AND status <> ?", appId, models.StatusRejected).
		Update("status", models.StatusRejected).
		Error; err != nil {
		return err
	}

	return nil
}

func (r *applicationRepository) AcceptApplication(appId uint) error {
	if err := r.db.Model(&models.Application{}).
		Where("id = ? AND status <> ?", appId, models.StatusAccepted).
		Update("status", models.StatusAccepted).
		Error; err != nil {
		return err
	}

	return nil
}

func (r *applicationRepository) Applications(id uint, filter models.ApplicationFilter) ([]models.Application, error) {
	var apps []models.Application
	query := r.db.Model(&models.Application{}).
		Joins("JOIN vacancies ON vacancies.id = applications.vacancy_id ").
		Where("vacancies.company_id = ?", id)
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if err := query.Find(&apps).Error; err != nil {
		return nil, err
	}

	return apps, nil
}

func (r *applicationRepository) Create(application *models.Application) error {
	return r.db.Create(&application).Error
}
