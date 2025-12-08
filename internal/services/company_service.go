package services

import (
	"github.com/AliUmarov/team-find-me-job/internal/models"
	"github.com/AliUmarov/team-find-me-job/internal/repository"
)

type CompanyService struct {
	companyRepo repository.CompanyRepository
}

func NewCompanyService(repo repository.CompanyRepository) *CompanyService {
	return &CompanyService{companyRepo: repo}
}

func (s *CompanyService) List() ([]models.Company, error) {
	return s.companyRepo.List()
}

func (s *CompanyService) Create(dto models.CompanyCreateRequest) (*models.Company, error) {
	company := &models.Company{
		Name:        dto.Name,
		Description: dto.Description,
		Website:     dto.Website,
	}
	if err := s.companyRepo.Create(company); err != nil {
		return company, err
	}

	return company, nil
}
