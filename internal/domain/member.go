package domain

import "github.com/google/uuid"

type Member struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v7();primaryKey"`
	TripID uuid.UUID `gorm:"type:uuid;not null;index"`
	Name   string    `gorm:"size:255;not null"`
}
