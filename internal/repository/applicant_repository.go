package repository

import (
	"log/slog"

	"github.com/AliUmarov/team-find-me-job/internal/models"
	"gorm.io/gorm"
)

type ApplicantRepositoryInterface interface {
	List() ([]models.Applicant, error)
	Create(applicant *models.Applicant) (*models.Applicant, error)
	GetByID(id uint) (*models.Applicant, error)
	Update(id uint, applicant *models.UpdateApplicant) (*models.Applicant, error)
	Delete(id uint) error
}

type ApplicantRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewApplicantRepository(db *gorm.DB, logger *slog.Logger) *ApplicantRepository {
	return &ApplicantRepository{
		db:     db,
		logger: logger,
	}
}

func (r *ApplicantRepository) List() ([]models.Applicant, error) {
	var applicants []models.Applicant
	if err := r.db.Model(&applicants).Find(&applicants).Error; err != nil {
		r.logger.Error("не удалось получить список соискателей",
			slog.Any("error", err),
		)
		return nil, err
	}

	r.logger.Info("список соискателей успешно получен",
		slog.Int("count", len(applicants)),
	)
	return applicants, nil
}

func (r *ApplicantRepository) Create(applicant *models.Applicant) (*models.Applicant, error) {
	if err := r.db.Model(&applicant).Create(&applicant).Error; err != nil {
		r.logger.Error("ошибка при создании соискателя",
			slog.String("full_name", applicant.FullName),
			slog.Any("error", err),
		)
		return &models.Applicant{}, err
	}
	r.logger.Info("соискатель успешно создан",
		slog.Uint64("id", uint64(applicant.ID)),
		slog.String("full_name", applicant.FullName),
	)
	return applicant, nil
}

func (r *ApplicantRepository) GetByID(id uint) (*models.Applicant, error) {
	var applicant models.Applicant
	if err := r.db.Model(&applicant).Where("id = ?", id).First(&applicant).Error; err != nil {
		r.logger.Error("не удалось получить соискателя",
			slog.Uint64("id", uint64(id)),
			slog.Any("error", err),
		)
		return &models.Applicant{}, err
	}
	r.logger.Info("соискатель успешно получен",
		slog.Uint64("id", uint64(applicant.ID)),
		slog.String("full_name", applicant.FullName),
	)
	return &applicant, nil
}

func (r *ApplicantRepository) Update(id uint, applicant *models.UpdateApplicant) (*models.Applicant, error) {
	var existingApplicant models.Applicant
	if err := r.db.Model(&existingApplicant).Where("id = ?", id).Updates(&applicant).Error; err != nil {
		r.logger.Error("не удалось обновить соискателя",
			slog.Uint64("id", uint64(id)),
			slog.Any("error", err),
		)
		return &models.Applicant{}, err
	}

	r.logger.Info("соискатель успешно обновлен",
		slog.Uint64("id", uint64(existingApplicant.ID)),
		slog.String("full_name", existingApplicant.FullName),
	)

	return &existingApplicant, nil
}

func (r *ApplicantRepository) Delete(id uint) error {
	var existingApplicant models.Applicant
	if err := r.db.Model(&existingApplicant).Where("id = ?", id).First(&existingApplicant).Error; err != nil {
		r.logger.Error("не удалось найти соискателя для удаления",
			slog.Uint64("id", uint64(id)),
			slog.Any("error", err),
		)
		return err
	}

	r.logger.Info("соискатель успешно удален",
		slog.Uint64("id", uint64(existingApplicant.ID)),
		slog.String("full_name", existingApplicant.FullName),
	)
	return r.db.Model(&existingApplicant).Delete(&existingApplicant).Error
}
