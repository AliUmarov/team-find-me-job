package models

type Vacancy struct {
	Base

	Title       string
	Description string
	Salary      int
	Rating      float64
	CompanyID   uint
}

type CreateVacancy struct {
	Title       string
	Description string
	Salary      int
	Rating      float64
	CompanyID   uint
}
