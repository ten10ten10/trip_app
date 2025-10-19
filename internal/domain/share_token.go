package domain

import (
	"time"

	"github.com/google/uuid"
)

type ShareToken struct {
	TripID    uuid.UUID `gorm:"column:trip_id;type:uuid;not null;primaryKey"`
	TokenHash string    `gorm:"column:token_hash;size:255;not null;uniqueIndex"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamptz;not null;autoCreateTime:false"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamptz;not null;autoUpdateTime:false"`
}
