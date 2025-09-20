package domain

import (
	"time"

	"github.com/google/uuid"
)

type Schedule struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v7();primaryKey"`
	TripID        uuid.UUID `gorm:"type:uuid;not null;index"`
	Title         string    `gorm:"size:255;not null"`
	StartDateTime time.Time `gorm:"type:timestamptz;not null"`
	EndDateTime   time.Time `gorm:"type:timestamptz;not null"`
	Memo          string    `gorm:"type:text"`
	CreatedAt     time.Time `gorm:"column:createdAt;type:timestamptz;not null;autoCreateTime:false"`
	UpdatedAt     time.Time `gorm:"column:updatedAt;type:timestamptz;not null;autoUpdateTime:false"`
}
