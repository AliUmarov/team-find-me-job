package models

type Company struct {
	Base

	Name        string  `json:"name" gorm:"type:varchar(100);not null"`
	Description string  `json:"description" gorm:"type:varchar(1000);not null"`
	Website     string  `json:"website" gorm:"type:varchar(255);not null"`
	Rating      float64 `json:"rating" gorm:"type:double precision"`
	ReviewCount int     `json:"review_count"`

	Vacancies []Vacancy `json:"vacancies" gorm:"constraint:OnDelete:RESTRICT;"`
}

type CompanyCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Website     string `json:"website"`
}
