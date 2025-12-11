package repository

import (
	"github.com/AliUmarov/team-find-me-job/internal/models"
	"gorm.io/gorm"
)

type ApplicationRepository interface {
	Create(*models.Application) error
}

type applicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) ApplicationRepository {
	return &applicationRepository{db: db}
}

func (r *applicationRepository) Create(application *models.Application) error {
	return r.db.Create(&application).Error
}
