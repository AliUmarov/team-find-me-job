package repository

import (
	"strings"

	"github.com/AliUmarov/team-find-me-job/internal/models"
	"gorm.io/gorm"
)

type VacancyRepository interface {
	Search(models.VacancyFilter) ([]models.Vacancy, error)
	Create(*models.Vacancy) error
	GetByCompanyId(uint) ([]models.Vacancy, error)
	IsVacancyExists(id uint) (bool, error)
}

type vacancyRepository struct {
	db *gorm.DB
}

func NewVacancyRepository(db *gorm.DB) VacancyRepository {
	return &vacancyRepository{db: db}
}

func (r *vacancyRepository) Search(filter models.VacancyFilter) ([]models.Vacancy, error) {
	var vacancies []models.Vacancy
	query := r.db.Model(&models.Vacancy{})

	title := strings.ReplaceAll(*filter.Title, "%", "\\%")
	title = strings.ReplaceAll(title, "_", "\\_")
	if filter.Title != nil {
		query = query.Where("title ILIKE ? ESCAPE '\\'", "%"+title+"%")
	}
	if err := query.Find(&vacancies).Error; err != nil {
		return nil, err
	}

	return vacancies, nil
}

func (r *vacancyRepository) GetByCompanyId(id uint) ([]models.Vacancy, error) {
	var vacancies []models.Vacancy

	if err := r.db.Where("company_id = ?", id).Find(&vacancies).Error; err != nil {
		return nil, err
	}

	return vacancies, nil
}

func (r *vacancyRepository) Create(vacancy *models.Vacancy) error {
	return r.db.Create(&vacancy).Error
}

func (r *vacancyRepository) IsVacancyExists(id uint) (bool, error) {
	var count int64
	if err := r.db.Model(&models.Vacancy{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
