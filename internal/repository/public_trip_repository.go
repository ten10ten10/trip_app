package repository

import (
	"context"
	"trip_app/internal/domain"

	"gorm.io/gorm"
)

type PublicTripRepository interface {
	FindByShareToken(ctx context.Context, shareToken string) (*domain.Trip, error)
	Update(ctx context.Context, trip *domain.Trip) error
	FindWithSchedulesByShareToken(ctx context.Context, shareToken string) (*domain.Trip, error)
}

type publicTripRepository struct {
	db *gorm.DB
}

func NewPublicTripRepository(db *gorm.DB) PublicTripRepository {
	return &publicTripRepository{db}
}

func (r *publicTripRepository) FindByShareToken(ctx context.Context, shareToken string) (*domain.Trip, error) {
	var trip domain.Trip
	var token domain.ShareToken
	if err := r.db.WithContext(ctx).First(&token, "token_hash = ?", shareToken).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Preload("Members").First(&trip, "id = ?", token.TripID).Error; err != nil {
		return nil, err
	}
	return &trip, nil
}

func (r *publicTripRepository) Update(ctx context.Context, trip *domain.Trip) error {
	return r.db.WithContext(ctx).Save(trip).Error
}

func (r *publicTripRepository) FindWithSchedulesByShareToken(ctx context.Context, shareToken string) (*domain.Trip, error) {
	var trip domain.Trip
	var token domain.ShareToken
	if err := r.db.WithContext(ctx).First(&token, "token_hash = ?", shareToken).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Preload("Members").Preload("Schedules").First(&trip, "id = ?", token.TripID).Error; err != nil {
		return nil, err
	}
	return &trip, nil
}
