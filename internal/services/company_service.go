package services

import (
	"github.com/AliUmarov/team-find-me-job/internal/models"
	"github.com/AliUmarov/team-find-me-job/internal/repository"
)

type CompanyService interface {
	List() ([]models.Company, error)
	Create(models.CompanyCreateRequest) (*models.Company, error)
	GetVacanciesByCompanyId(uint) ([]models.Vacancy, error)
	Applications(uint, models.ApplicationFilter) ([]models.Application, error)
}

type companyService struct {
	companyRepo     repository.CompanyRepository
	vacancyRepo     repository.VacancyRepository
	applicationRepo repository.ApplicationRepository
}

func NewCompanyService(
	companyRepo repository.CompanyRepository,
	vacancyRepo repository.VacancyRepository,
	applicationRepo repository.ApplicationRepository,
) CompanyService {
	return &companyService{
		companyRepo:     companyRepo,
		vacancyRepo:     vacancyRepo,
		applicationRepo: applicationRepo,
	}
}

func (s *companyService) Applications(id uint, filter models.ApplicationFilter) ([]models.Application, error) {
	_, err := s.companyRepo.Get(id)
	if err != nil {
		return nil, err
	}
	return s.applicationRepo.Applications(id, filter)
}

func (s *companyService) GetVacanciesByCompanyId(id uint) ([]models.Vacancy, error) {
	return s.vacancyRepo.GetByCompanyId(id)
}

func (s *companyService) List() ([]models.Company, error) {
	return s.companyRepo.List()
}

func (s *companyService) Create(dto models.CompanyCreateRequest) (*models.Company, error) {
	company := &models.Company{
		Name:        dto.Name,
		Description: dto.Description,
		Website:     dto.Website,
	}
	if err := s.companyRepo.Create(company); err != nil {
		return nil, err
	}

	return company, nil
}
