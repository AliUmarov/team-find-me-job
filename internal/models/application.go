package models

type Application struct {
	Base

	Status    string `json:"status" gorm:"type:varchar(100);not null;oneof('pending','reviewed','accepted','rejected');default:'pending'"`
	VacancyID uint   `json:"vacancy_id" gorm:"not null"`
	ResumeID  uint   `json:"resume_id" gorm:"not null"`

	Vacancy *Vacancy `json:"-" gorm:"foreignKey:VacancyID"`
	Resume  *Resume  `json:"-" gorm:"foreignKey:ResumeID"`
}

type CreateApplication struct {
	VacancyID uint `json:"vacancy_id"`
	ResumeID  uint `json:"resume_id"`
}
