package transport

import (
	"github.com/AliUmarov/team-find-me-job/internal/services"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.Engine,
	companyService services.CompanyService,
) {
	companyHandler := NewCompanyHandler(companyService)

	companyHandler.RegisterRoutes(router)
}
