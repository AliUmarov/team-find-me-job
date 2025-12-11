package services

import (
	"fmt"
	"log/slog"

	"github.com/AliUmarov/team-find-me-job/internal/gigachat"
	"github.com/AliUmarov/team-find-me-job/internal/models"
	"github.com/AliUmarov/team-find-me-job/internal/repository"
)

type ResumeService interface {
	Create(id uint, req models.ResumeCreateRequest) (*models.Resume, error)
	GetAllResumes() ([]models.Resume, error)
	GetByID(id uint) (*models.Resume, error)
	Update(id uint, req models.ResumeUpdateRequest) (*models.Resume, error)
	Delete(id uint) error
	ImproveResume(id uint) (*models.Resume, error)
}

type resumeService struct {
	repo   repository.ResumeRepository
	logger *slog.Logger
	client *gigachat.Client
}

func NewResumeService(repo repository.ResumeRepository, logger *slog.Logger, client *gigachat.Client) ResumeService {
	return &resumeService{
		repo:   repo,
		logger: logger,
		client: client,
	}
}

func (s *resumeService) Create(id uint, req models.ResumeCreateRequest) (*models.Resume, error) {
	var resume = models.Resume{
		Position:    req.Position,
		Summary:     req.Summary,
		Skills:      req.Skills,
		Experience:  req.Experience,
		Portfolio:   req.Portfolio,
		Salary:      req.Salary,
		ApplicantID: id,
	}

	if err := s.repo.Create(&resume); err != nil {
		s.logger.Error("ошибка при добавлении резюме",
			slog.String("position", req.Position),
			slog.Any("error", err),
		)
		return nil, err
	}

	return &resume, nil
}

func (s *resumeService) GetAllResumes() ([]models.Resume, error) {
	resumes, err := s.repo.GetAllResumes()
	if err != nil {
		s.logger.Error("не удалось получить список резюме",
			slog.Any("error", err),
		)
		return nil, err
	}

	return resumes, nil
}

func (s *resumeService) GetByID(id uint) (*models.Resume, error) {
	resume, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("резюме не найдено",
			slog.Any("error", err),
		)
		return nil, err
	}

	return resume, nil
}

func (s *resumeService) Update(id uint, req models.ResumeUpdateRequest) (*models.Resume, error) {
	resume, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("резюме не найдено",
			slog.Any("error", err),
		)
		return nil, err
	}

	if req.Position != nil {
		resume.Position = *req.Position
	}

	if req.Summary != nil {
		resume.Summary = *req.Summary
	}

	if req.Skills != nil {
		resume.Skills = *req.Skills
	}

	if req.Experience != nil {
		resume.Experience = *req.Experience
	}

	if req.Portfolio != nil {
		resume.Portfolio = *req.Portfolio
	}

	if req.Salary != nil {
		resume.Salary = *req.Salary
	}

	if err := s.repo.Update(id, resume); err != nil {
		s.logger.Error("не удалось сохранить изменения",
			slog.Uint64("resume_id", uint64(id)),
			slog.Any("error", err),
		)
		return nil, err
	}

	return resume, nil
}

func (s *resumeService) Delete(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		s.logger.Error("ошибка при удалении резюме",
			slog.Uint64("resume_id", uint64(id)),
			slog.Any("error", err),
		)
		return err
	}

	return nil
}

func (s *resumeService) ImproveResume(id uint) (*models.Resume, error) {
	resume, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	fullText := fmt.Sprintf(`
		Position: %s
		Summary: %s
		Skills: %s
		Experience: %s
		Portfolio: %s 

		Salary: %d
	`, resume.Position, resume.Summary, resume.Skills, resume.Experience, resume.Portfolio, resume.Salary)

	improved, score, err := gigachat.ImproveResume(fullText, s.client)
	if err != nil {
		return nil, err
	}

	resume.AIImproved = improved
	resume.AIScore = score

	if err := s.repo.Save(resume); err != nil {
		return nil, err
	}

	return resume, nil
}
