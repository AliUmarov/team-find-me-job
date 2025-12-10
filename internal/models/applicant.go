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
