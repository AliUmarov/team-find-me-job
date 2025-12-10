package models

type Vacancy struct {
	Base

	Title       string  `json:"title" gorm:"type:varchar(255);not null"`
	Description string  `json:"description" gorm:"type:text;not null"`
	Salary      int     `json:"salary" gorm:"type:int;not null"`
	Rating      float64 `json:"rating" gorm:"type:float;not null"`
	CompanyID   uint    `json:"company_id" gorm:"not null"`
	
	Company *Company `json:"-" gorm:"foreignKey:CompanyID"`
}

type CreateVacancy struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Salary      int     `json:"salary"`
	Rating      float64 `json:"rating"`
	CompanyID   uint    `json:"company_id"`
}
