package transport

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/AliUmarov/team-find-me-job/internal/constants"
	"github.com/AliUmarov/team-find-me-job/internal/dto"
	"github.com/AliUmarov/team-find-me-job/internal/middlewares"
	"github.com/AliUmarov/team-find-me-job/internal/services"
	"github.com/gin-gonic/gin"
)

type ApplicantHandler struct {
	service     services.ApplicantService
	authService services.AuthService
	logger      *slog.Logger
}

func NewApplicantHandler(service services.ApplicantService, authService services.AuthService, logger *slog.Logger) *ApplicantHandler {
	return &ApplicantHandler{service: service, authService: authService, logger: logger}
}

func (h *ApplicantHandler) RegisterRoutes(r *gin.Engine) {
	jwtService := h.authService.GetJWTService()
	applicant := r.Group("/applicant")
	{
		applicant.GET("", h.List)
		applicant.GET("/:id", middlewares.Authenticate(*jwtService), h.GetByID)
		applicant.PUT("/:id", middlewares.Authenticate(*jwtService), h.Update)
		applicant.DELETE("/:id", middlewares.Authenticate(*jwtService), h.Delete)
	}
}

func (h *ApplicantHandler) List(c *gin.Context) {
	applicants, err := h.service.List()
	if err != nil {
		h.logger.Error("ошибка при получении списка соискателей",
			slog.String("error", err.Error()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ERR_CAN_NOT_GET_APPLICANT})
		return
	}

	h.logger.Info("список соискателей успешно получен",
		slog.Int("count", len(applicants)),
	)

	c.JSON(http.StatusOK, gin.H{"data": applicants})
}

func (h *ApplicantHandler) Create(c *gin.Context) {
	var req dto.CreateApplicant
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("некорректный JSON",
			slog.String("method", c.Request.Method),
			slog.String("path", c.FullPath()),
			slog.String("error", err.Error()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ERR_INCORRECT_PAYLOAD})
		return
	}
	applicant, err := h.service.Create(&req)
	if err != nil {
		h.logger.Error("ошибка при создании соискателя",
			slog.String("error", err.Error()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ERR_CAN_NOT_CREATE_APPLICANT})
		return
	}

	h.logger.Info("соискатель успешно создан",
		slog.Any("data", applicant),
	)
	c.JSON(http.StatusCreated, gin.H{"data": applicant})
}

func (h *ApplicantHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	var id uint
	if _, err := fmt.Sscan(idParam, &id); err != nil {
		h.logger.Warn("некорректный ID",
			slog.String("error", err.Error()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ERR_INCORRECT_ID})
		return
	}

	applicant, err := h.service.GetByID(id)
	if err != nil {
		h.logger.Error("ошибка при получении соискателя по ID",
			slog.Uint64("id", uint64(id)),
			slog.String("error", err.Error()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ERR_CAN_NOT_GET_APPLICANT})
		return
	}

	h.logger.Info("соискатель успешно получен",
		slog.Any("data", applicant),
	)
	c.JSON(http.StatusOK, gin.H{"data": applicant})
}

func (h *ApplicantHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	var id uint
	if _, err := fmt.Sscan(idParam, &id); err != nil {
		h.logger.Warn("некорректный ID",
			slog.String("error", err.Error()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ERR_INCORRECT_ID})
		return
	}

	var req dto.UpdateApplicant
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("некорректный JSON",
			slog.String("method", c.Request.Method),
			slog.String("path", c.FullPath()),
			slog.String("error", err.Error()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ERR_INCORRECT_PAYLOAD})
		return
	}

	updatedApplicant, err := h.service.Update(id, req)
	if err != nil {
		h.logger.Error("ошибка при обновлении соискателя",
			slog.Uint64("id", uint64(id)),
			slog.String("error", err.Error()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ERR_CAN_NOT_UPDATE_APPLICANT})
		return
	}

	h.logger.Info("соискатель успешно обновлен",
		slog.Any("data", updatedApplicant),
	)
	c.JSON(http.StatusOK, gin.H{"data": updatedApplicant})
}

func (h *ApplicantHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	var id uint
	if _, err := fmt.Sscan(idParam, &id); err != nil {
		h.logger.Warn("некорректный ID",
			slog.Any("error", err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ERR_INCORRECT_ID})
		return
	}

	if err := h.service.Delete(id); err != nil {
		h.logger.Error("ошибка при удалении соискателя",
			slog.Uint64("id", uint64(id)),
			slog.String("error", err.Error()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ERR_CAN_NOT_DELETE_APPLICANT})
		return
	}

	h.logger.Info("соискатель успешно удален",
		slog.Uint64("id", uint64(id)),
	)
	c.JSON(http.StatusOK, gin.H{"message": constants.ERR_CAN_NOT_DELETE_APPLICANT})
}
