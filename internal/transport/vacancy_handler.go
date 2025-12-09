package transport

import (
	"net/http"

	"github.com/AliUmarov/team-find-me-job/internal/models"
	"github.com/AliUmarov/team-find-me-job/internal/services"
	"github.com/gin-gonic/gin"
)

type VacancyHandler struct {
	service services.VacancyService
}

func NewVacancyHandler(service services.VacancyService) *VacancyHandler {
	return &VacancyHandler{service: service}
}

func (h *VacancyHandler) RegisterRoutes(r *gin.Engine) {
	vacancy := r.Group("/vacancy")
	{
		vacancy.GET("", h.Search)
		vacancy.POST("", h.Create)
	}
}

func (h *VacancyHandler) Search(c *gin.Context) {
	var filter models.VacancyFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	vacancies, err := h.service.Search(filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": vacancies})
}

func (h *VacancyHandler) Create(c *gin.Context) {
	var req models.VacancyCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	vacancy, err := h.service.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": vacancy})
}
