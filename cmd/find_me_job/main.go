package main

import (
	"log/slog"
	"os"

	"github.com/AliUmarov/team-find-me-job/internal/config"
	"github.com/AliUmarov/team-find-me-job/internal/models"
	"github.com/AliUmarov/team-find-me-job/internal/repository"
	"github.com/AliUmarov/team-find-me-job/internal/services"
	"github.com/AliUmarov/team-find-me-job/internal/transport"
	"github.com/gin-gonic/gin"
)

func main() {
	log := config.InitLogger()

	config.SetEnv(log)
	db := config.Connect(log)

	if err := db.AutoMigrate(
		&models.Company{},
		&models.Vacancy{},
		&models.Resume{},
		&models.Application{},
	); err != nil {
		log.Error("failed to migrate database", "error", err)
		os.Exit(1)
	}

	log.Info("migrations completed")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	companyRepo := repository.NewCompanyRepository(db)
	applicantRepo := repository.NewApplicantRepository(db, log)
	vacancyRepo := repository.NewVacancyRepository(db)
	resumeRepo := repository.NewResumeRepository(db, log)
	applicationRepo := repository.NewApplicationRepository(db)

	applicantService := services.NewApplicantService(applicantRepo, log)
	resumeService := services.NewResumeService(resumeRepo, log)
	companyService := services.NewCompanyService(companyRepo, vacancyRepo)
	vacancyService := services.NewVacancyService(vacancyRepo)
	applicationService := services.NewApplicationService(applicationRepo, vacancyRepo, resumeRepo)

	r := gin.Default()

	transport.RegisterRoutes(r, log, companyService, applicantService, resumeService, vacancyService, applicationService)

	log.Info("server started",
		slog.String("addr", port))

	if err := r.Run(":" + port); err != nil {
		log.Error("не удалось запустить сервер", slog.Any("error", err))
	}
}
