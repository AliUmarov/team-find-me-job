package models

import (
	"gorm.io/gorm"
)

type Applicant struct {
	Base

	FullName string `json:"full_name" gorm:"type:varchar(255);not null"`
	Email    string `json:"email" gorm:"type:varchar(255);not null;uniqueIndex"`
	Phone    string `json:"phone" gorm:"type:varchar(255);not null;uniqueIndex"`
	Password string `gorm:"type:varchar(255);not null" json:"password"`
	Role     string `gorm:"type:varchar(50);not null;default:'APPLICANT'" json:"role"`

	IsVerified bool `gorm:"default:false" json:"is_verified"`
}



// BeforeCreate hook to hash password and set defaults
func (u *Applicant) BeforeCreate(_ *gorm.DB) (err error) {

	// Set default role if not specified
	if u.Role == "" {
		u.Role = "APPLICANT"
	}

	return nil
}
