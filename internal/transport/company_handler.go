package transport

import (
	"net/http"
	"strconv"

	"github.com/AliUmarov/team-find-me-job/internal/models"
	"github.com/AliUmarov/team-find-me-job/internal/services"
	"github.com/gin-gonic/gin"
)

type CompanyHandler struct {
	service services.CompanyService
}

func NewCompanyHandler(service services.CompanyService) *CompanyHandler {
	return &CompanyHandler{service: service}
}

func (h *CompanyHandler) RegisterRoutes(r *gin.Engine) {
	company := r.Group("/companies")
	{
		company.GET(":id/applications/:app/accept", h.AcceptApplication)
		company.GET(":id/applications/:app/reject", h.RejectApplication)
		company.GET(":id/applications", h.Applications)
		company.GET("", h.List)
		company.POST("", h.Create)
		company.GET(":id/vacancies", h.GetVacanciesByCompanyId)
	}
}

func (h *CompanyHandler) RejectApplication(c *gin.Context) {
	idStr := c.Param("id")
	appStr := c.Param("app")
	companyId, err1 := strconv.ParseUint(appStr, 10, 64)
	appId, err2 := strconv.ParseUint(idStr, 10, 64)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return
	}
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
		return
	}
	if err := h.service.RejectApplication(uint(companyId), uint(appId)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h *CompanyHandler) AcceptApplication(c *gin.Context) {
	idStr := c.Param("id")
	appStr := c.Param("app")
	companyId, err1 := strconv.ParseUint(appStr, 10, 64)
	appId, err2 := strconv.ParseUint(idStr, 10, 64)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return
	}
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
		return
	}
	if err := h.service.AcceptApplication(uint(companyId), uint(appId)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h *CompanyHandler) Applications(c *gin.Context) {
	var applicationsFilter models.ApplicationFilter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindQuery(&applicationsFilter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	applications, err := h.service.Applications(uint(id), applicationsFilter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": applications})
}

func (h *CompanyHandler) GetVacanciesByCompanyId(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	companies, err := h.service.GetVacanciesByCompanyId(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": companies})
}

func (h *CompanyHandler) List(c *gin.Context) {
	companies, err := h.service.List()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": companies})
}

func (h *CompanyHandler) Create(c *gin.Context) {
	var req models.CompanyCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	company, err := h.service.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": company})
}
