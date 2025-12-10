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
	resumeService services.ResumeService,
	vacancyService services.VacancyService,
) {
	companyHandler := NewCompanyHandler(companyService)
	resumeHandler := NewResumeHandler(resumeService, logger)
	vacancyHandler := NewVacancyHandler(vacancyService)

	companyHandler.RegisterRoutes(router)
	resumeHandler.RegisterRoutes(router)
	vacancyHandler.RegisterRoutes(router)
}
