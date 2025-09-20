// internal/domain/user.go

package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v7();primaryKey"`
	Name                       string    `gorm:"size:255;not null"`
	Email                      string    `gorm:"size:255;uniqueIndex;not null"`
	PasswordHash               string    `gorm:"size:255;not null"`
	IsActive                   bool      `gorm:"not null;default:false"`
	VerificationTokenHash      *string   `gorm:"size:255;uniqueIndex"`
	VerificationTokenExpiresAt *time.Time
	CreatedAt                  time.Time `gorm:"column:createdAt;type:timestamptz;not null;autoCreateTime:false"`
	UpdatedAt                  time.Time `gorm:"column:updatedAt;type:timestamptz;not null;autoUpdateTime:false"`
	Trips                      []Trip    `gorm:"constraint:OnDelete:CASCADE"`
}
