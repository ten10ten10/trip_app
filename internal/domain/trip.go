package domain

import (
	"time"

	"github.com/google/uuid"
)

type Trip struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v7();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Title     string    `gorm:"size:255;not null"`
	StartDate time.Time `gorm:"type:date;not null"`
	EndDate   time.Time `gorm:"type:date;not null"`
	CreatedAt time.Time `gorm:"column:createdAt;type:timestamptz;not null;autoCreateTime:false"`
	UpdatedAt time.Time `gorm:"column:updatedAt;type:timestamptz;not null;autoUpdateTime:false"`

	Members    []Member   `gorm:"foreignKey:TripID;constraint:OnDelete:CASCADE"`
	Schedules  []Schedule `gorm:"foreignKey:TripID;constraint:OnDelete:CASCADE"`
	ShareToken ShareToken `gorm:"foreignKey:TripID;constraint:OnDelete:CASCADE"`
}
