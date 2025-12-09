package transport

import (
	"github.com/AliUmarov/team-find-me-job/internal/services"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.Engine,
	companyService services.CompanyService,
	vacancyService services.VacancyService,
) {
	companyHandler := NewCompanyHandler(companyService)
	vacancyHandler := NewVacancyHandler(vacancyService)

	companyHandler.RegisterRoutes(router)
	vacancyHandler.RegisterRoutes(router)
}
