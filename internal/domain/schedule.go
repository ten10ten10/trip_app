package domain

import (
	"time"

	"github.com/google/uuid"
)

type Schedule struct {
	ID            uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v7();primaryKey"`
	TripID        uuid.UUID `gorm:"column:trip_id;type:uuid;not null;index"`
	Title         string    `gorm:"column:title;size:255;not null"`
	StartDateTime time.Time `gorm:"column:start_date_time;type:timestamptz;not null"`
	EndDateTime   time.Time `gorm:"column:end_date_time;type:timestamptz;not null"`
	Memo          string    `gorm:"column:memo;type:text"`
	CreatedAt     time.Time `gorm:"column:created_at;type:timestamptz;not null;autoCreateTime:false"`
	UpdatedAt     time.Time `gorm:"column:updated_at;type:timestamptz;not null;autoUpdateTime:false"`
}
