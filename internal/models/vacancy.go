package models

import "github.com/lib/pq"

type Vacancy struct {
	Base

	Title       string  `json:"title" gorm:"type:varchar(255);not null"`
	Description string  `json:"description" gorm:"type:text;not null"`
	Salary      int     `json:"salary" gorm:"type:int;not null"`
	Rating      float64 `json:"rating" gorm:"type:float;not null"`
	Requirements     pq.StringArray `json:"requirements" gorm:"type:text[];not null"`
	Responsibilities pq.StringArray `json:"responsibilities" gorm:"type:text[];not null"`
	NiceToHave       pq.StringArray `json:"nice_to_have" gorm:"type:text[];not null"`

	CompanyID uint `json:"company_id" binding:"required" gorm:"not null"`
	Company   Company `json:"-"`
}

type VacancyCreateRequest struct {
	Title            string   `json:"title" binding:"required"`
	Description      string   `json:"description" binding:"required"`
	Salary           int      `json:"salary" binding:"required,gt=0"`
	CompanyID        uint     `json:"company_id" binding:"required"`
	Requirements     []string `json:"requirements" binding:"required"`
	Responsibilities []string `json:"responsibilities" binding:"required"`
	NiceToHave       []string `json:"nice_to_have" binding:"required"`
}

type VacancyFilter struct {
	Title *string `form:"title"`
}
