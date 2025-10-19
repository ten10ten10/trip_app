package domain

import "github.com/google/uuid"

type Member struct {
	ID     uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v7();primaryKey"`
	TripID uuid.UUID `gorm:"column:trip_id;type:uuid;not null;index"`
	Name   string    `gorm:"column:name;size:255;not null"`
}
