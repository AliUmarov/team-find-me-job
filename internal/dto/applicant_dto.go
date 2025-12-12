package dto

import (
	"errors"

	"github.com/AliUmarov/team-find-me-job/internal/models"
)

const (
	// Failed
	MESSAGE_FAILED_GET_DATA_FROM_BODY = "failed get data from body"
	MESSAGE_FAILED_REGISTER_USER      = "failed create user"
	MESSAGE_FAILED_GET_LIST_USER      = "failed get list user"
	MESSAGE_FAILED_TOKEN_NOT_VALID    = "token not valid"
	MESSAGE_FAILED_TOKEN_NOT_FOUND    = "token not found"
	MESSAGE_FAILED_GET_USER           = "failed get user"
	MESSAGE_FAILED_LOGIN              = "failed login"
	MESSAGE_FAILED_UPDATE_USER        = "failed update user"
	MESSAGE_FAILED_DELETE_USER        = "failed delete user"
	MESSAGE_FAILED_PROSES_REQUEST     = "failed proses request"
	MESSAGE_FAILED_DENIED_ACCESS      = "denied access"
	MESSAGE_FAILED_VERIFY_EMAIL       = "failed verify email"

	// Success
	MESSAGE_SUCCESS_REGISTER_USER           = "success create user"
	MESSAGE_SUCCESS_GET_LIST_USER           = "success get list user"
	MESSAGE_SUCCESS_GET_USER                = "success get user"
	MESSAGE_SUCCESS_LOGIN                   = "success login"
	MESSAGE_SUCCESS_UPDATE_USER             = "success update user"
	MESSAGE_SUCCESS_DELETE_USER             = "success delete user"
	MESSAGE_SEND_VERIFICATION_EMAIL_SUCCESS = "success send verification email"
	MESSAGE_SUCCESS_VERIFY_EMAIL            = "success verify email"
)

var (
	ErrCreateUser             = errors.New("failed to create user")
	ErrGetUserById            = errors.New("failed to get user by id")
	ErrGetUserByEmail         = errors.New("failed to get user by email")
	ErrEmailAlreadyExists     = errors.New("email already exist")
	ErrUpdateUser             = errors.New("failed to update user")
	ErrUserNotFound           = errors.New("user not found")
	ErrEmailNotFound          = errors.New("email not found")
	ErrDeleteUser             = errors.New("failed to delete user")
	ErrTokenInvalid           = errors.New("token invalid")
	ErrTokenExpired           = errors.New("token expired")
	ErrAccountAlreadyVerified = errors.New("account already verified")
)

type (
	CreateApplicant struct {
		FullName string `json:"full_name"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
	}

	UpdateApplicant struct {
		FullName *string `json:"full_name"`
		Email    *string `json:"email"`
		Phone    *string `json:"phone"`
	}
	
	ApplicantRegisterRequest struct {
		FullName string `json:"full_name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Phone    string `json:"phone" binding:"required"`
		Password string `json:"password" binding:"required,min=8"`
	}

	ApplicantResponse struct {
		models.Base
		FullName   string `json:"full_name" binding:"required"`
		Email      string `json:"email" binding:"required,email"`
		Phone      string `json:"phone" binding:"required"`
		Role       string `json:"role"`
		IsVerified bool   `json:"is_verified"`
	}

	ApplicantUpdateRequest struct {
		FullName string `json:"full_name" binding:"omitempty"`
		Email    string `json:"email" binding:"omitempty,email"`
		Phone    string `json:"phone" binding:"omitempty"`
	}

	ApplicantUpdateResponse struct {
		models.Base
		FullName   string `json:"full_name" binding:"omitempty"`
		Email      string `json:"email" binding:"omitempty,email"`
		Phone      string `json:"phone" binding:"omitempty"`
		Role       string `json:"role"`
		IsVerified bool   `json:"is_verified"`
	}

	SendVerificationEmailRequest struct {
		Email string `json:"email" form:"email" binding:"required"`
	}

	VerifyEmailRequest struct {
		Token string `json:"token" form:"token" binding:"required"`
	}

	VerifyEmailResponse struct {
		Email      string `json:"email"`
		IsVerified bool   `json:"is_verified"`
	}

	ApplicantLoginRequest struct {
		Email    string `json:"email" form:"email" binding:"required"`
		Password string `json:"password" form:"password" binding:"required"`
	}
)
