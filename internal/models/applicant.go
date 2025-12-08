package models

import "gorm.io/gorm"

type Applicant struct {
	gorm.Model

	FullName string
	Email    string
	Phone    string
}

type Application struct {
	gorm.Model

	VacancyID   uint
	ResumeID    uint
	ApplicantID uint
}
