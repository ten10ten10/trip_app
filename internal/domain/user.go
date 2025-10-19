// internal/domain/user.go

package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                         uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v7();primaryKey"`
	Name                       string    `gorm:"column:name;size:255;not null"`
	Email                      string    `gorm:"column:email;size:255;uniqueIndex;not null"`
	PasswordHash               string    `gorm:"column:password_hash;size:255;not null"`
	IsActive                   bool      `gorm:"column:is_active;not null;default:false"`
	VerificationTokenHash      *string   `gorm:"column:verification_token_hash;size:255;uniqueIndex"`
	VerificationTokenExpiresAt *time.Time `gorm:"column:verification_token_expires_at"`
	CreatedAt                  time.Time `gorm:"column:created_at;type:timestamptz;not null;autoCreateTime:false"`
	UpdatedAt                  time.Time `gorm:"column:updated_at;type:timestamptz;not null;autoUpdateTime:false"`
	Trips                      []Trip    `gorm:"foreignKey:user_id;constraint:OnDelete:CASCADE"`
}
