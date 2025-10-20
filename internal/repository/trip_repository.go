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
	FindWithSchedulesByID(ctx context.Context, tripID uuid.UUID) (*domain.Trip, error)
	Delete(ctx context.Context, tripID uuid.UUID) error
}

type tripRepository struct {
	db *gorm.DB
}

func NewTripRepository(db *gorm.DB) TripRepository {
	return &tripRepository{db}
}

func (r *tripRepository) Create(ctx context.Context, trip *domain.Trip) error {
	if err := r.db.WithContext(ctx).Create(trip).Error; err != nil {
		return err
	}
	return nil
}

func (r *tripRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Trip, error) {
	var trips []domain.Trip
	if err := r.db.WithContext(ctx).Preload("Members").Where("user_id = ?", userID).Find(&trips).Error; err != nil {
		return nil, err
	}
	return trips, nil
}

func (r *tripRepository) FindByID(ctx context.Context, tripID uuid.UUID) (*domain.Trip, error) {
	var trip domain.Trip
	if err := r.db.WithContext(ctx).Preload("Members").First(&trip, "id = ?", tripID).Error; err != nil {
		return nil, err
	}
	return &trip, nil
}

func (r *tripRepository) FindByShareToken(ctx context.Context, shareTokenHash string) (*domain.Trip, error) {
	var trip domain.Trip
	if err := r.db.WithContext(ctx).Preload("Members").Preload("Schedules").Preload("ShareToken").
		Joins("ShareToken").Where("ShareToken.token_hash = ?", shareTokenHash).First(&trip).Error; err != nil {
		return nil, err
	}
	return &trip, nil
}

func (r *tripRepository) Update(ctx context.Context, trip *domain.Trip) error {
	if err := r.db.WithContext(ctx).Save(trip).Error; err != nil {
		return err
	}
	return nil
}

func (r *tripRepository) FindWithSchedulesByID(ctx context.Context, tripID uuid.UUID) (*domain.Trip, error) {
	var trip domain.Trip
	if err := r.db.WithContext(ctx).Preload("Members").Preload("Schedules").First(&trip, "id = ?", tripID).Error; err != nil {
		return nil, err
	}
	return &trip, nil
}

func (r *tripRepository) Delete(ctx context.Context, tripID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&domain.Trip{}, "id = ?", tripID).Error; err != nil {
		return err
	}
	return nil
}
