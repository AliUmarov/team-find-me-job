package models

import "gorm.io/gorm"

type Vacancy struct {
	gorm.Model

	Title       string
	Description string
	Salary      int
	Rating      float64
	CompanyID   uint
	// Company     Company
}
