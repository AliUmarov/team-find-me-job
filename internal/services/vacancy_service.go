package services

import (
	"github.com/AliUmarov/team-find-me-job/internal/models"
	"github.com/AliUmarov/team-find-me-job/internal/repository"
)

type VacancyService interface {
	Search(models.VacancyFilter) ([]models.Vacancy, error)
	Create(dto models.VacancyCreateRequest) (*models.Vacancy, error)
}

type vacancyService struct {
	vacancyRepo repository.VacancyRepository
}

func NewVacancyService(repo repository.VacancyRepository) VacancyService {
	return &vacancyService{vacancyRepo: repo}
}

func (s *vacancyService) Search(filter models.VacancyFilter) ([]models.Vacancy, error) {
	return s.vacancyRepo.Search(filter)
}

func (s *vacancyService) Create(dto models.VacancyCreateRequest) (*models.Vacancy, error) {
	vacancy := &models.Vacancy{
		Title:            dto.Title,
		Description:      dto.Description,
		Salary:           dto.Salary,
		Requirements:     dto.Requirements,
		Responsibilities: dto.Responsibilities,
		NiceToHave:       dto.NiceToHave,
		CompanyID:        dto.CompanyID,
	}
	if err := s.vacancyRepo.Create(vacancy); err != nil {
		return nil, err
	}

	return vacancy, nil
}
