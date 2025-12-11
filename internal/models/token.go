package models

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID    uint `gorm:"not null" json:"user_id"`
	Token     string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"token"`
	ExpiresAt time.Time `gorm:"type:timestamp with time zone;not null" json:"expires_at"`
	User      Applicant `gorm:"foreignKey:ApplicantID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`

	Timestamp
}
