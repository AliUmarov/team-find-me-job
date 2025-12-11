package transport

import (
	"log/slog"
	"net/http"

	"github.com/AliUmarov/team-find-me-job/internal/constants"
	"github.com/AliUmarov/team-find-me-job/internal/dto"
	"github.com/AliUmarov/team-find-me-job/internal/pkg/utils"
	"github.com/AliUmarov/team-find-me-job/internal/services"
	"github.com/AliUmarov/team-find-me-job/internal/validation"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service    services.AuthService
	logger     *slog.Logger
	validation *validation.AuthValidation
}

func (h *AuthHandler) RegisterRoutes(r *gin.Engine) {
	authRoutes := r.Group("/api/auth")
	{
		authRoutes.POST("/register", h.Register)
		authRoutes.POST("/login", h.Login)
		authRoutes.POST("/refresh", h.RefreshToken)
		authRoutes.POST("/logout", h.Logout)
		authRoutes.POST("/send-verification-email", h.SendVerificationEmail)
		authRoutes.POST("/verify-email", h.VerifyEmail)
		authRoutes.POST("/send-password-reset", h.SendPasswordReset)
		authRoutes.POST("/reset-password", h.ResetPassword)
	}
}

func NewAuthHandler(service services.AuthService, logger *slog.Logger) *AuthHandler {
	authValidation := validation.NewAuthValidation()
	return &AuthHandler{
		service:    service,
		logger:     logger,
		validation: authValidation,
	}
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	var req dto.ApplicantRegisterRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": constants.MESSAGE_FAILED_GET_DATA_FROM_BODY})
		return
	}

	// Validate request
	if err := h.validation.ValidateRegisterRequest(req); err != nil {
		res := utils.BuildResponseFailed("Validation failed", err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := h.service.RegisterApplicant(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_REGISTER_USER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_REGISTER_USER, result)
	ctx.JSON(http.StatusOK, res)
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var req dto.ApplicantLoginRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// Validate request
	if err := h.validation.ValidateLoginRequest(req); err != nil {
		res := utils.BuildResponseFailed("Validation failed", err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := h.service.Login(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_LOGIN, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_LOGIN, result)
	ctx.JSON(http.StatusOK, res)
}

func (h *AuthHandler) RefreshToken(ctx *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := h.service.RefreshToken(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_REFRESH_TOKEN, err.Error(), nil)
		ctx.JSON(http.StatusUnauthorized, res)
		return
	}

	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_REFRESH_TOKEN, result)
	ctx.JSON(http.StatusOK, res)
}

func (h *AuthHandler) Logout(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(string)

	err := h.service.Logout(ctx.Request.Context(), userId)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_LOGOUT, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_LOGOUT, nil)
	ctx.JSON(http.StatusOK, res)
}

func (h *AuthHandler) SendVerificationEmail(ctx *gin.Context) {
	var req dto.SendVerificationEmailRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	err := h.service.SendVerificationEmail(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_PROSES_REQUEST, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(constants.MESSAGE_SEND_VERIFICATION_EMAIL_SUCCESS, nil)
	ctx.JSON(http.StatusOK, res)
}

func (h *AuthHandler) VerifyEmail(ctx *gin.Context) {
	var req dto.VerifyEmailRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := h.service.VerifyEmail(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_VERIFY_EMAIL, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_VERIFY_EMAIL, result)
	ctx.JSON(http.StatusOK, res)
}

func (h *AuthHandler) SendPasswordReset(ctx *gin.Context) {
	var req dto.SendPasswordResetRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	err := h.service.SendPasswordReset(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_SEND_PASSWORD_RESET, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_SEND_PASSWORD_RESET, nil)
	ctx.JSON(http.StatusOK, res)
}

func (h *AuthHandler) ResetPassword(ctx *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	err := h.service.ResetPassword(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_RESET_PASSWORD, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(constants.MESSAGE_SUCCESS_RESET_PASSWORD, nil)
	ctx.JSON(http.StatusOK, res)
}
