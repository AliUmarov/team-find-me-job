package transport

import (
	"log/slog"

	"github.com/AliUmarov/team-find-me-job/internal/services"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.Engine,
	logger *slog.Logger,
	companyService services.CompanyService,
	applicantService services.ApplicantService,
	resumeService services.ResumeService,
	vacancyService services.VacancyService,
	authService services.AuthService,
	applicationService services.ApplicationService,
) {
	authHandler := NewAuthHandler(authService, logger)

	companyHandler := NewCompanyHandler(companyService)
	resumeHandler := NewResumeHandler(resumeService, logger)
	applicantHandler := NewApplicantHandler(applicantService, logger)
	vacancyHandler := NewVacancyHandler(vacancyService)
	applicationHandler := NewApplicationHandler(applicationService)

	companyHandler.RegisterRoutes(router)
	applicantHandler.RegisterRoutes(router)
	resumeHandler.RegisterRoutes(router)
	vacancyHandler.RegisterRoutes(router)

	authHandler.RegisterRoutes(router)
	applicationHandler.RegisterRoutes(router)
}
