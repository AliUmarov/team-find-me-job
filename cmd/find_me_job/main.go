package main

import (
	"log/slog"
	"os"

	"github.com/AliUmarov/team-find-me-job/internal/config"
	"github.com/AliUmarov/team-find-me-job/internal/gigachat"
	"github.com/AliUmarov/team-find-me-job/internal/middlewares"
	"github.com/AliUmarov/team-find-me-job/internal/models"
	"github.com/AliUmarov/team-find-me-job/internal/repository"
	"github.com/AliUmarov/team-find-me-job/internal/services"
	"github.com/AliUmarov/team-find-me-job/internal/transport"
	"github.com/AliUmarov/team-find-me-job/internal/validators"
	"github.com/gin-gonic/gin"
)

func main() {
	log := config.InitLogger()
	validators.RegisterValidators()

	config.SetEnv(log)
	db := config.Connect(log)

	if err := db.AutoMigrate(
		&models.Company{},
		&models.Vacancy{},
		&models.Resume{},
		&models.Applicant{},
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

	tokenProvider, err := gigachat.NewTokenProvider()
	if err != nil {
		log.Error("failed to init GigaChat TokenProvider", slog.Any("error", err))
		os.Exit(1)
	}

	gigaClient, err := gigachat.NewClient(tokenProvider)
	if err != nil {
		log.Error("failed to init GigaChat client", slog.Any("error", err))
		os.Exit(1)
	}

	companyRepo := repository.NewCompanyRepository(db)
	applicantRepo := repository.NewApplicantRepository(db, log)
	vacancyRepo := repository.NewVacancyRepository(db)
	resumeRepo := repository.NewResumeRepository(db, log)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)

	jwtService := services.NewJWTService()
	authService := services.NewAuthService(applicantRepo, companyRepo, log, refreshTokenRepo, jwtService, db)
	applicationRepo := repository.NewApplicationRepository(db)

	applicantService := services.NewApplicantService(applicantRepo, log)
	resumeService := services.NewResumeService(resumeRepo, applicantRepo, log, gigaClient)
	companyService := services.NewCompanyService(companyRepo, vacancyRepo, applicationRepo)
	vacancyService := services.NewVacancyService(vacancyRepo)
	applicationService := services.NewApplicationService(applicationRepo, vacancyRepo, resumeRepo, db)

	r := gin.Default()
	r.Use(middlewares.CORSMiddleware())

	transport.RegisterRoutes(r, log, companyService, applicantService, resumeService, vacancyService, applicationService, authService, db)

	log.Info("server started",
		slog.String("addr", port))

	if err := r.Run(":" + port); err != nil {
		log.Error("не удалось запустить сервер", slog.Any("error", err))
	}
}
