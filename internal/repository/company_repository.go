package repository

import (
	"github.com/AliUmarov/team-find-me-job/internal/models"
	"gorm.io/gorm"
)

type CompanyRepository interface {
	Get(uint) (*models.Company, error)
	List() ([]models.Company, error)
	Create(company *models.Company) error
}

type companyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) CompanyRepository {
	return &companyRepository{db: db}
}

func (r *companyRepository) Get(id uint) (*models.Company, error) {
	var company models.Company
	if err := r.db.First(&company, id).Error; err != nil {
		return nil, err
	}

	return &company, nil
}

func (r *companyRepository) List() ([]models.Company, error) {
	var companies []models.Company
	if err := r.db.Find(&companies).Error; err != nil {
		return companies, err
	}

	return companies, nil
}

func (r *companyRepository) Create(company *models.Company) error {
	return r.db.Create(&company).Error
}
