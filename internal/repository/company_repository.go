package repository

import (
	"github.com/AliUmarov/team-find-me-job/internal/models"
	"gorm.io/gorm"
)

type CompanyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

func (r *CompanyRepository) List() ([]models.Company, error) {
	var companies []models.Company
	if err := r.db.Find(&companies).Error; err != nil {
		return companies, err
	}

	return companies, nil
}

func (r *CompanyRepository) Create(company *models.Company) error {
	return r.db.Create(&company).Error
}
