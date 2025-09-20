package domain

import (
	"time"

	"github.com/google/uuid"
)

type ShareToken struct {
	TripID    uuid.UUID `gorm:"type:uuid;not null;primaryKey"`
	TokenHash string    `gorm:"size:255;not null;uniqueIndex;column:token_hash"`
	CreatedAt time.Time `gorm:"column:createdAt;type:timestamptz;not null;autoCreateTime:false"`
	UpdatedAt time.Time `gorm:"column:updatedAt;type:timestamptz;not null;autoUpdateTime:false"`
}
