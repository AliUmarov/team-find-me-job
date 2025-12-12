package services

import (
	"errors"

	"github.com/AliUmarov/team-find-me-job/internal/models"
	"github.com/AliUmarov/team-find-me-job/internal/repository"
	"gorm.io/gorm"
)

type ApplicationService interface {
	Create(models.CreateApplication) error
}

type applicationService struct {
	applicationRepo repository.ApplicationRepository
	vacancyRepo     repository.VacancyRepository
	resumeRepo      repository.ResumeRepository
	db              *gorm.DB
}

func NewApplicationService(
	applicationRepo repository.ApplicationRepository,
	vacancyRepo repository.VacancyRepository,
	resumeRepo repository.ResumeRepository,
	db *gorm.DB,
) ApplicationService {
	return &applicationService{
		applicationRepo: applicationRepo,
		vacancyRepo:     vacancyRepo,
		resumeRepo:      resumeRepo,
		db:              db,
	}
}

func (s *applicationService) Create(dto models.CreateApplication) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		applicationRepo := s.applicationRepo.WithDb(tx)
		vacancyRepo := s.vacancyRepo.WithDb(tx)
		resumeRepo := s.resumeRepo.WithDb(tx)
		isVacancyExists, err := vacancyRepo.IsVacancyExists(dto.VacancyID)
		if err != nil {
			return err
		}
		if isVacancyExists == false {
			return errors.New("vacancy is not exists")
		}

		isResumeExists, err := resumeRepo.IsResumeExists(dto.ResumeID)
		if err != nil {
			return err
		}
		if isResumeExists == false {
			return errors.New("resume is not exists")
		}

		application := &models.Application{
			Status:    models.StatusPending,
			VacancyID: dto.VacancyID,
			ResumeID:  dto.ResumeID,
		}

		if err := applicationRepo.Create(application); err != nil {
			return err
		}

		return nil
	})
}
