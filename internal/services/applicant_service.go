package services

import (
	"log/slog"

	"github.com/AliUmarov/team-find-me-job/internal/dto"
	"github.com/AliUmarov/team-find-me-job/internal/models"
	"github.com/AliUmarov/team-find-me-job/internal/repository"
)

type ApplicantServiceInterface interface {
	Create(applicant models.Applicant) (models.Applicant, error)
	GetByID(id uint) (models.Applicant, error)
	Update(id uint, applicant dto.UpdateApplicant) (models.Applicant, error)
	Delete(id uint) error
}

type ApplicantService struct {
	applicantRepo repository.ApplicantRepository
	companyRepo   repository.CompanyRepository
	logger        *slog.Logger
}

func NewApplicantService(applicantRepo repository.ApplicantRepository, logger *slog.Logger) ApplicantService {
	return ApplicantService{applicantRepo: applicantRepo, logger: logger}
}

func (s *ApplicantService) List() ([]models.Applicant, error) {
	applicants, err := s.applicantRepo.List()
	if err != nil {
		s.logger.Error("ошибка при получении списка соискателей",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	s.logger.Info("список соискателей успешно получен",
		slog.Int("count", len(applicants)),
	)

	return applicants, nil
}

func (s *ApplicantService) Create(applicant *dto.CreateApplicant) (*models.Applicant, error) {
	if applicant == nil {
		s.logger.Error("пустой соискатель не может быть создан")
		return &models.Applicant{}, nil
	}

	createdApplicant := &models.Applicant{
		FullName: applicant.FullName,
		Email:    applicant.Email,
		Phone:    applicant.Phone,
	}

	return s.applicantRepo.Create(createdApplicant)
}

func (s *ApplicantService) GetByID(id uint) (*models.Applicant, error) {
	applicant, err := s.applicantRepo.GetByID(id)
	if err != nil {
		s.logger.Error("ошибка при получении соискателя по ID",
			slog.Uint64("id", uint64(id)),
			slog.String("error", err.Error()),
		)
		return &models.Applicant{}, err
	}
	return applicant, nil
}

func (s *ApplicantService) Update(id uint, updApplicant dto.UpdateApplicant) (*models.Applicant, error) {
	existing, err := s.applicantRepo.GetByID(id)
	if err != nil {
		s.logger.Error("ошибка при получении соискателя для обновления",
			slog.Uint64("id", uint64(id)),
			slog.String("error", err.Error()),
		)
		return &models.Applicant{}, err
	}


    // Применяем только те поля, которые пришли
    if updApplicant.FullName != nil {
        existing.FullName = *updApplicant.FullName
    }
    if updApplicant.Email != nil {
        existing.Email = *updApplicant.Email
    }
    if updApplicant.Phone != nil {
        existing.Phone = *updApplicant.Phone
    }

	updatedApplicant, err := s.applicantRepo.Update(id, existing)
	if err != nil {
		s.logger.Error("ошибка при обновлении соискателя",
			slog.Uint64("id", uint64(id)),
			slog.String("error", err.Error()),
		)
		return &models.Applicant{}, err
	}
	
	return updatedApplicant, nil
}

func (s *ApplicantService) Delete(id uint) error {
	if err := s.applicantRepo.Delete(id); err != nil {
		s.logger.Error("ошибка при удалении соискателя",
			slog.Uint64("id", uint64(id)),
			slog.String("error", err.Error()),
		)
		return err
	}
	return nil
}
