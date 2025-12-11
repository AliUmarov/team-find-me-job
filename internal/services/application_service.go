package services

import (
	"errors"

	"github.com/AliUmarov/team-find-me-job/internal/models"
	"github.com/AliUmarov/team-find-me-job/internal/repository"
)

type ApplicationService interface {
	Create(models.CreateApplication) (*models.Application, error)
}

type applicationService struct {
	applicationRepo repository.ApplicationRepository
	vacancyRepo     repository.VacancyRepository
	resumeRepo      repository.ResumeRepository
}

func NewApplicationService(
	applicationRepo repository.ApplicationRepository,
	vacancyRepo repository.VacancyRepository,
	resumeRepo repository.ResumeRepository,
) ApplicationService {
	return &applicationService{
		applicationRepo: applicationRepo,
		vacancyRepo:     vacancyRepo,
		resumeRepo:      resumeRepo,
	}
}

func (s *applicationService) Create(dto models.CreateApplication) (*models.Application, error) {
	isVacancyExists, err := s.vacancyRepo.IsVacancyExists(dto.VacancyID)
	if err != nil {
		return nil, err
	}
	if isVacancyExists == false {
		return nil, errors.New("vacancy is not exists")
	}

	isResumeExists, err := s.resumeRepo.IsResumeExists(dto.ResumeID)
	if err != nil {
		return nil, err
	}
	if isResumeExists == false {
		return nil, errors.New("resume is not exists")
	}

	application := &models.Application{
		Status:    models.StatusPending,
		VacancyID: dto.VacancyID,
		ResumeID:  dto.ResumeID,
	}

	if err := s.applicationRepo.Create(application); err != nil {
		return nil, err
	}

	return application, nil
}
