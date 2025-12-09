package repository

import (
	"log/slog"

	"github.com/AliUmarov/team-find-me-job/internal/models"
	"gorm.io/gorm"
)

type ResumeRepository interface {
	Create(resume *models.Resume) error
	GetAllResumes() ([]models.Resume, error)
	GetByID(id uint) (*models.Resume, error)
	Update(id uint, resume *models.Resume) error
	Delete(id uint) error
}

type gormResumeRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewResumeRepository(db *gorm.DB, logger *slog.Logger) ResumeRepository {
	return &gormResumeRepository{
		db:     db,
		logger: logger,
	}
}

func (r *gormResumeRepository) Create(resume *models.Resume) error {
	op := "repo.resume.create"

	r.logger.Debug("db call",
		slog.String("op", op),
		slog.String("position", resume.Position),
	)

	err := r.db.Create(resume).Error

	if err != nil {
		r.logger.Error("db error",
			slog.String("op", op),
			slog.Any("error", err),
		)
		return err
	}

	return nil
}

func (r *gormResumeRepository) GetAllResumes() ([]models.Resume, error) {
	op := "repo.resume.getall"

	r.logger.Debug("db call",
		slog.String("op", op),
	)

	var resumes []models.Resume

	if err := r.db.Find(&resumes).Error; err != nil {
		r.logger.Error("db error",
			slog.String("op", op),
			slog.Any("error", err),
		)
		return nil, err
	}

	return resumes, nil
}

func (r *gormResumeRepository) GetByID(id uint) (*models.Resume, error) {
	op := "repo.resume.get_by_id"

	r.logger.Debug("db call",
		slog.String("op", op),
	)

	var resume models.Resume

	if err := r.db.First(&resume, id).Error; err != nil {
		r.logger.Error("db error",
			slog.String("op", op),
			slog.Any("error", err),
		)
		return nil, err
	}

	return &resume, nil

}

func (r *gormResumeRepository) Update(id uint, resume *models.Resume) error {
	op := "repo.resume.update"

	r.logger.Debug("db call",
		slog.String("op", op),
	)

	if err := r.db.Model(&models.Resume{}).Where("id = ?", id).Updates(resume).Error; err != nil {
		r.logger.Error("db error",
			slog.String("op", op),
			slog.Any("error", err),
		)
		return err
	}

	return nil
}

func (r *gormResumeRepository) Delete(id uint) error {
	op := "repo.resume.delete"

	r.logger.Debug("db call",
		slog.String("op", op),
		slog.Uint64("resume_id", uint64(id)),
	)

	if err := r.db.Delete(&models.Resume{}, id).Error; err != nil {
		r.logger.Error("db error",
			slog.String("op", op),
			slog.Uint64("resume_id", uint64(id)),
			slog.Any("error", err),
		)
		return err
	}

	return nil
}
