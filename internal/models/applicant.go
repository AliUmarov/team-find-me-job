package models

type Applicant struct {
	Base

	FullName string `json:"full_name" gorm:"type:varchar(255);not null"`
	Email    string `json:"email" gorm:"type:varchar(255);not null;uniqueIndex"`
	Phone    string `json:"phone" gorm:"type:varchar(255);not null;uniqueIndex"`
}

type CreateApplicant struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type UpdateApplicant struct {
	FullName *string `json:"full_name"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
}

type Application struct {
	Base

	Status    string `json:"status" gorm:"type:varchar(100);not null;oneof('pending','reviewed','accepted','rejected');default:'pending'"`
	VacancyID uint   `json:"vacancy_id" gorm:"not null"`
	ResumeID  uint   `json:"resume_id" gorm:"not null"`
}

type CreateApplication struct {
	VacancyID uint `json:"vacancy_id"`
	ResumeID  uint `json:"resume_id"`

	Vacancy *Vacancy `json:"vacancy,omitempty" gorm:"-"`
	Resume  *Resume  `json:"resume,omitempty" gorm:"-"`
}
