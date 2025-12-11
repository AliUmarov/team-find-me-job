package services

import (
	"context"
	"log/slog"

	"github.com/AliUmarov/team-find-me-job/internal/constants"
	"github.com/AliUmarov/team-find-me-job/internal/dto"
	"github.com/AliUmarov/team-find-me-job/internal/models"
	"github.com/AliUmarov/team-find-me-job/internal/pkg/helpers"
	"github.com/AliUmarov/team-find-me-job/internal/pkg/utils"
	"github.com/AliUmarov/team-find-me-job/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService interface {
	RegisterApplicant(ctx context.Context, req dto.ApplicantRegisterRequest) (dto.ApplicantResponse, error)
	Login(ctx context.Context, req dto.ApplicantLoginRequest) (dto.TokenResponse, error)
	RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (dto.TokenResponse, error)
	Logout(ctx context.Context, userId string) error
	SendVerificationEmail(ctx context.Context, req dto.SendVerificationEmailRequest) error
	VerifyEmail(ctx context.Context, req dto.VerifyEmailRequest) (dto.VerifyEmailResponse, error)
	SendPasswordReset(ctx context.Context, req dto.SendPasswordResetRequest) error
	ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error
}

type authService struct {
	applicantRepo          repository.ApplicantRepository
	companyRepo            repository.CompanyRepository
	refreshTokenRepository repository.RefreshTokenRepository
	jwtService             JWTService
	db                     *gorm.DB
}

func NewAuthService(
	applicantRepo repository.ApplicantRepository,
	companyRepo repository.CompanyRepository,
	logger *slog.Logger,
	refreshTokenRepo repository.RefreshTokenRepository,
	jwtService JWTService,
	db *gorm.DB,
) AuthService {
	return &authService{
		applicantRepo:          applicantRepo,
		companyRepo:            companyRepo,
		refreshTokenRepository: refreshTokenRepo,
		jwtService:             jwtService,
		db:                     db,
	}
}

func (s *authService) RegisterApplicant(ctx context.Context, req dto.ApplicantRegisterRequest) (dto.ApplicantResponse, error) {
	_, isExist, err := s.applicantRepo.CheckEmail(ctx, s.db, req.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return dto.ApplicantResponse{}, err
	}

	if isExist {
		return dto.ApplicantResponse{}, dto.ErrEmailAlreadyExists
	}

	hashedPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		return dto.ApplicantResponse{}, err
	}

	user := models.Applicant{
		FullName:   req.FullName,
		Email:      req.Email,
		Phone:      req.Phone,
		Password:   hashedPassword,
		Role:       "APPLICANT",
		IsVerified: false,
	}

	createdUser, err := s.applicantRepo.Register(ctx, s.db, user)
	if err != nil {
		return dto.ApplicantResponse{}, err
	}

	return dto.ApplicantResponse{
		FullName:   createdUser.FullName,
		Email:      createdUser.Email,
		Phone:      createdUser.Phone,
		Role:       createdUser.Role,
		IsVerified: createdUser.IsVerified,
	}, nil
}

func (s *authService) Login(ctx context.Context, req dto.ApplicantLoginRequest) (dto.TokenResponse, error) {
	user, err := s.applicantRepo.GetApplicantByEmail(ctx, s.db, req.Email)
	if err != nil {
		return dto.TokenResponse{}, dto.ErrEmailNotFound
	}

	isValid, err := helpers.CheckPassword(user.Password, []byte(req.Password))
	if err != nil || !isValid {
		return dto.TokenResponse{}, constants.ErrInvalidCredentials
	}

	accessToken := s.jwtService.GenerateAccessToken(user.ID, user.Role)
	refreshTokenString, expiresAt := s.jwtService.GenerateRefreshToken()

	refreshToken := models.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     refreshTokenString,
		ExpiresAt: expiresAt,
	}

	_, err = s.refreshTokenRepository.Create(ctx, s.db, refreshToken)
	if err != nil {
		return dto.TokenResponse{}, err
	}

	return dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
		Role:         user.Role,
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (dto.TokenResponse, error) {
	refreshToken, err := s.refreshTokenRepository.FindByToken(ctx, s.db, req.RefreshToken)
	if err != nil {
		return dto.TokenResponse{}, constants.ErrRefreshTokenNotFound
	}

	accessToken := s.jwtService.GenerateAccessToken(refreshToken.UserID, refreshToken.User.Role)
	newRefreshTokenString, expiresAt := s.jwtService.GenerateRefreshToken()

	err = s.refreshTokenRepository.DeleteByToken(ctx, s.db, req.RefreshToken)
	if err != nil {
		return dto.TokenResponse{}, err
	}

	newRefreshToken := models.RefreshToken{
		ID:        uuid.New(),
		UserID:    refreshToken.UserID,
		Token:     newRefreshTokenString,
		ExpiresAt: expiresAt,
	}

	_, err = s.refreshTokenRepository.Create(ctx, s.db, newRefreshToken)
	if err != nil {
		return dto.TokenResponse{}, err
	}

	return dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshTokenString,
		Role:         refreshToken.User.Role,
	}, nil
}

func (s *authService) Logout(ctx context.Context, userId string) error {
	return s.refreshTokenRepository.DeleteByUserID(ctx, s.db, userId)
}

func (s *authService) SendVerificationEmail(ctx context.Context, req dto.SendVerificationEmailRequest) error {
	user, err := s.applicantRepo.GetApplicantByEmail(ctx, s.db, req.Email)
	if err != nil {
		return dto.ErrEmailNotFound
	}

	if user.IsVerified {
		return dto.ErrAccountAlreadyVerified
	}

	verificationToken := s.jwtService.GenerateAccessToken(user.ID, "verification")

	subject := "Email Verification"
	body := "Please verify your email using this token: " + verificationToken

	return utils.SendMail(user.Email, subject, body)
}

func (s *authService) VerifyEmail(ctx context.Context, req dto.VerifyEmailRequest) (dto.VerifyEmailResponse, error) {
	token, err := s.jwtService.ValidateToken(req.Token)
	if err != nil || !token.Valid {
		return dto.VerifyEmailResponse{}, dto.ErrTokenInvalid
	}

	userId, err := s.jwtService.GetUserIDByToken(req.Token)
	if err != nil {
		return dto.VerifyEmailResponse{}, dto.ErrTokenInvalid
	}

	user, err := s.applicantRepo.GetByID(userId)
	if err != nil {
		return dto.VerifyEmailResponse{}, dto.ErrUserNotFound
	}

	user.IsVerified = true
	updatedUser, err := s.applicantRepo.Update(userId, user)
	if err != nil {
		return dto.VerifyEmailResponse{}, err
	}

	return dto.VerifyEmailResponse{
		Email:      updatedUser.Email,
		IsVerified: updatedUser.IsVerified,
	}, nil
}

func (s *authService) SendPasswordReset(ctx context.Context, req dto.SendPasswordResetRequest) error {
	user, err := s.applicantRepo.GetApplicantByEmail(ctx, s.db, req.Email)
	if err != nil {
		return dto.ErrEmailNotFound
	}

	resetToken := s.jwtService.GenerateAccessToken(user.ID, "password_reset")

	subject := "Password Reset"
	body := "Please reset your password using this token: " + resetToken

	return utils.SendMail(user.Email, subject, body)
}

func (s *authService) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error {
	token, err := s.jwtService.ValidateToken(req.Token)
	if err != nil || !token.Valid {
		return constants.ErrPasswordResetToken
	}

	userId, err := s.jwtService.GetUserIDByToken(req.Token)
	if err != nil {
		return constants.ErrPasswordResetToken
	}

	user, err := s.applicantRepo.GetByID(userId)
	if err != nil {
		return dto.ErrUserNotFound
	}

	hashedPassword, err := helpers.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	_, err = s.applicantRepo.Update(userId, user)
	if err != nil {
		return err
	}

	return nil
}
