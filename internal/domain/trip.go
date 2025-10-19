package domain

import (
	"time"

	"github.com/google/uuid"
)

type Trip struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v7();primaryKey"`
	UserID    uuid.UUID `gorm:"column:user_id;type:uuid;not null;index"`
	Title     string    `gorm:"column:title;size:255;not null"`
	StartDate time.Time `gorm:"column:start_date;type:date;not null"`
	EndDate   time.Time `gorm:"column:end_date;type:date;not null"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamptz;not null;autoCreateTime:false"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamptz;not null;autoUpdateTime:false"`

	Members    []Member   `gorm:"foreignKey:trip_id;constraint:OnDelete:CASCADE"`
	Schedules  []Schedule `gorm:"foreignKey:trip_id;constraint:OnDelete:CASCADE"`
	ShareToken ShareToken `gorm:"foreignKey:trip_id;constraint:OnDelete:CASCADE"`
}
