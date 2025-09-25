package repository

import (
	"context"

	"trip_app/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TripRepository interface {
	Create(ctx context.Context, trip *domain.Trip) error
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Trip, error)
	FindByID(ctx context.Context, tripID uuid.UUID) (*domain.Trip, error)
	FindByShareToken(ctx context.Context, shareToken string) (*domain.Trip, error)
	Update(ctx context.Context, trip *domain.Trip) error
	Delete(ctx context.Context, tripID uuid.UUID) error
}

type tripRepository struct {
	db *gorm.DB
}

func NewTripRepository(db *gorm.DB) TripRepository {
	return &tripRepository{db}
}
