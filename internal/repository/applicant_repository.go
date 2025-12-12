package repository

import (
	"context"
	"log/slog"

	"github.com/AliUmarov/team-find-me-job/internal/dto"
	"github.com/AliUmarov/team-find-me-job/internal/models"
	"gorm.io/gorm"
)

type ApplicantRepositoryInterface interface {
	List() ([]models.Applicant, error)
	Create(applicant *models.Applicant) (*models.Applicant, error)
	GetByID(id uint) (*dto.ApplicantResponse, error)
	Update(id uint, applicant *dto.UpdateApplicant) (*models.Applicant, error)
	Delete(id uint) error
}

type ApplicantRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewApplicantRepository(db *gorm.DB, logger *slog.Logger) ApplicantRepository {
	return ApplicantRepository{
		db:     db,
		logger: logger,
	}
}

func (r *ApplicantRepository) List() ([]models.Applicant, error) {
	var applicants []models.Applicant
	if err := r.db.Find(&applicants).Error; err != nil {
		r.logger.Error("не удалось получить список соискателей",
			slog.String("error", err.Error()),
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
			slog.String("error", err.Error()),
		)
		return &models.Applicant{}, err
	}
	r.logger.Info("соискатель успешно создан",
		slog.Uint64("id", uint64(applicant.ID)),
		slog.String("full_name", applicant.FullName),
	)
	return applicant, nil
}

func (r *ApplicantRepository) GetByID(id uint) (*dto.ApplicantResponse, error) {
	var applicant models.Applicant
	if err := r.db.Model(&applicant).First(&applicant, id).Error; err != nil {
		r.logger.Error("не удалось получить соискателя",
			slog.Uint64("id", uint64(id)),
			slog.String("error", err.Error()),
		)
		return &dto.ApplicantResponse{}, err
	}

	resp := dto.ApplicantResponse{
		Base:       applicant.Base,
		FullName:   applicant.FullName,
		Email:      applicant.Email,
		Phone:      applicant.Phone,
		Role:       applicant.Role,
		IsVerified: applicant.IsVerified,
	}

	r.logger.Info("соискатель успешно получен",
		slog.Uint64("id", uint64(applicant.ID)),
		slog.String("full_name", applicant.FullName),
	)
	return &resp, nil
}

func (r *ApplicantRepository) Update(id uint, applicant *dto.ApplicantResponse) (*models.Applicant, error) {
	if err := r.db.Model(&models.Applicant{}).Where("id = ?", id).Updates(&applicant).Error; err != nil {
		r.logger.Error("не удалось обновить соискателя",
			slog.Uint64("id", uint64(id)),
			slog.String("error", err.Error()),
		)
		return &models.Applicant{}, err
	}

	r.logger.Info("соискатель успешно обновлен",
		slog.Uint64("id", uint64(id)),
		slog.Any("updates", applicant),
	)

	var existingApplicant models.Applicant

	if err := r.db.Model(&existingApplicant).First(&existingApplicant, id).Error; err != nil {
		r.logger.Error("не удалось получить обновленного соискателя",
			slog.Uint64("id", uint64(id)),
			slog.String("error", err.Error()),
		)
		return &models.Applicant{}, err
	}

	return &existingApplicant, nil
}

func (r *ApplicantRepository) Delete(id uint) error {
	r.logger.Info("удаление соискателя",
		slog.Uint64("id", uint64(id)),
	)
	return r.db.Model(&models.Applicant{}).Where("id = ?", id).Delete(&models.Applicant{}).Error
}

func (r *ApplicantRepository) Register(ctx context.Context, tx *gorm.DB, user models.Applicant) (*models.Applicant, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *ApplicantRepository) GetApplicantByEmail(ctx context.Context, tx *gorm.DB, email string) (*models.Applicant, error) {
	if tx == nil {
		tx = r.db
	}

	var user models.Applicant
	if err := tx.WithContext(ctx).Where("email = ?", email).Take(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *ApplicantRepository) CheckEmail(ctx context.Context, tx *gorm.DB, email string) (models.Applicant, bool, error) {
	if tx == nil {
		tx = r.db
	}

	var user models.Applicant
	if err := tx.WithContext(ctx).Where("email = ?", email).Take(&user).Error; err != nil {
		return models.Applicant{}, false, err
	}

	return user, true, nil
}
