package models

type ApplicationStatus string

const (
	StatusPending  ApplicationStatus = "pending"
	StatusReviewed ApplicationStatus = "reviewed"
	StatusAccepted ApplicationStatus = "accepted"
	StatusRejected ApplicationStatus = "rejected"
)

type Application struct {
	Base

	Status ApplicationStatus `json:"status" gorm:"type:varchar(100);not null;default:'pending'"`

	VacancyID uint `json:"vacancy_id" gorm:"not null"`
	ResumeID  uint `json:"resume_id" gorm:"not null"`

	Vacancy *Vacancy `json:"vacancy,omitempty" gorm:"foreignKey:VacancyID"`
	Resume  *Resume  `json:"resume,omitempty" gorm:"foreignKey:ResumeID"`
}

type CreateApplication struct {
	VacancyID uint `json:"vacancy_id" binding:"required"`
	ResumeID  uint `json:"resume_id" binding:"required"`
}
