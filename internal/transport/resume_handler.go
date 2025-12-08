package transport

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/AliUmarov/team-find-me-job/internal/models"
	"github.com/AliUmarov/team-find-me-job/internal/services"
	"github.com/gin-gonic/gin"
)

type ResumeHandler struct {
	service services.ResumeService
	logger  *slog.Logger
}

func NewResumeHandler(service services.ResumeService, logger *slog.Logger) *ResumeHandler {
	return &ResumeHandler{
		service: service,
		logger:  logger,
	}
}

func (h *ResumeHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/resumes")
	{
		api.GET("/", h.GetAll)
		api.PATCH("/:id", h.Update)
		api.DELETE("/:id", h.Delete)
	}

	r.POST("/applicants/:id/resumes", h.Create)
}

func (h *ResumeHandler) Create(c *gin.Context) {
	h.logger.Info("handler called",
		slog.String("method", c.Request.Method),
		slog.String("path", c.FullPath()),
	)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Warn("некорректный ID",
			slog.Any("error", err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID"})
		return
	}

	var req models.ResumeCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("некорректный JSON",
			slog.String("method", c.Request.Method),
			slog.String("path", c.FullPath()),
			slog.Any("error", err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resume, err := h.service.Create(uint(id), req)
	if err != nil {
		h.logger.Error("не удалось создать резюме",
			slog.Any("error", err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось создать резюме"})
		return
	}

	h.logger.Info("резюме успешно добавлено")
	c.JSON(http.StatusOK, resume)
}

func (h *ResumeHandler) GetAll(c *gin.Context) {
	h.logger.Info("handler called",
		slog.String("method", c.Request.Method),
		slog.String("path", c.FullPath()),
	)

	resume, err := h.service.GetAllResumes()
	if err != nil {
		h.logger.Error("не удалось получить список резюме",
			slog.Any("error", err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось получить список резюме"})
		return
	}

	c.JSON(http.StatusOK, resume)
}

func (h *ResumeHandler) Update(c *gin.Context) {
	h.logger.Info("handler called",
		slog.String("method", c.Request.Method),
		slog.String("path", c.FullPath()),
	)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Error("некорректный ID",
			slog.Any("error", err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID"})
		return
	}

	var req models.ResumeUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("некорректный JSON",
			slog.String("method", c.Request.Method),
			slog.String("path", c.FullPath()),
			slog.Any("error", err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный JSON"})
		return
	}

	updated, err := h.service.Update(uint(id), req)
	if err != nil {
		h.logger.Error("не удалось сохранить изменения",
			slog.String("method", c.Request.Method),
			slog.String("path", c.FullPath()),
			slog.Any("error", err),
		)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось сохранить изменения"})
		return
	}

	h.logger.Info("изменения сохранены")

	c.JSON(http.StatusOK, updated)
}

func (h *ResumeHandler) Delete(c *gin.Context) {
	h.logger.Info("handler called",
		slog.String("method", c.Request.Method),
		slog.String("path", c.FullPath()),
	)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Warn("некорректный ID",
			slog.String("id", idStr),
			slog.Any("error", err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID"})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		h.logger.Error("не удалось удалить резюме",
			slog.Uint64("resume_id", uint64(id)),
			slog.Any("error", err),
		)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось удалить резюме"})
		return
	}

	h.logger.Info("резюме удалено",
		slog.Uint64("resume_id", uint64(id)),
	)

	c.JSON(http.StatusOK, gin.H{"message": "резюме удалено"})
}
