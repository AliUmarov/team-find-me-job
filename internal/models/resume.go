package models

import "gorm.io/gorm"

type Resume struct {
	gorm.Model

	Position    string `json:"position"`
	Summary     string `json:"summary"`
	Skills      string `json:"skills"`
	Experience  string `json:"experience"`
	Portfolio   string `json:"portfolio"`
	Salary      int    `json:"salary"`
	AIImproved  string `json:"ai_improved"`
	AIScore     int    `json:"ai_score"`
	ApplicantID uint   `json:"applicant_id"`
}

type ResumeCreateRequest struct {
	Position   string `json:"position"`
	Summary    string `json:"summary"`
	Skills     string `json:"skills"`
	Experience string `json:"experience"`
	Portfolio  string `json:"portfolio"`
	Salary     int    `json:"salary"`
}

type ResumeUpdateRequest struct {
	Position   string `json:"position"`
	Summary    string `json:"summary"`
	Skills     string `json:"skills"`
	Experience string `json:"experience"`
	Portfolio  string `json:"portfolio"`
	Salary     int    `json:"salary"`
}
