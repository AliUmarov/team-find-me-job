package transport

import (
	"net/http"

	"github.com/AliUmarov/team-find-me-job/internal/models"
	"github.com/AliUmarov/team-find-me-job/internal/services"
	"github.com/gin-gonic/gin"
)

type ApplicationHandler struct {
	service services.ApplicationService
}

func NewApplicationHandler(service services.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{service: service}
}

func (h *ApplicationHandler) RegisterRoutes(r *gin.Engine) {
	application := r.Group("/applications")
	{
		application.POST("", h.Create)
	}
}

func (h *ApplicationHandler) Create(c *gin.Context) {
	var req models.CreateApplication
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.service.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success"})
}
