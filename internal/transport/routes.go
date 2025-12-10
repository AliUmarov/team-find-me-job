package transport

import (
	"log/slog"

	"github.com/AliUmarov/team-find-me-job/internal/services"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	logger *slog.Logger,
	router *gin.Engine,
	companyService services.CompanyService,
	applicantService services.ApplicantService,
) {
	companyHandler := NewCompanyHandler(companyService)
	applicantHandler := NewApplicantHandler(applicantService, logger)

	companyHandler.RegisterRoutes(router)
	applicantHandler.RegisterRoutes(router)
}
