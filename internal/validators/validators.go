package validators

import (
	"github.com/AliUmarov/team-find-me-job/internal/models"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func RegisterValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("application_status", validateApplicationStatus)
	}
}

func validateApplicationStatus(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	switch models.ApplicationStatus(value) {
	case models.StatusPending, models.StatusReviewed, models.StatusAccepted, models.StatusRejected:
		return true
	}
	return false
}
