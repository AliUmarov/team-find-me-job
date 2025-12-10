package models

type Applicant struct {
	Base

	FullName string
	Email    string
	Phone    string
}

type CreateApplicant struct {
	FullName string
	Email    string
	Phone    string
}

type UpdateApplicant struct {
	FullName *string
	Email    *string
	Phone    *string
}

type Application struct {
	Base

	VacancyID   uint
	ResumeID    uint
	ApplicantID uint
}

type CreateApplication struct {
	VacancyID   uint
	ResumeID    uint
	ApplicantID uint
}
