package repository

import (
	"context"
	"trip_app/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScheduleRepository interface {
	Create(ctx context.Context, schedule *domain.Schedule) error
	FindByTripID(ctx context.Context, tripID uuid.UUID) ([]domain.Schedule, error)
	FindByID(ctx context.Context, scheduleID uuid.UUID) (*domain.Schedule, error)
	Update(ctx context.Context, schedule *domain.Schedule) error
	Delete(ctx context.Context, scheduleID uuid.UUID) error
}

type scheduleRepository struct {
	db *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) ScheduleRepository {
	return &scheduleRepository{db}
}

func (r *scheduleRepository) Create(ctx context.Context, schedule *domain.Schedule) error {
	if err := r.db.WithContext(ctx).Create(schedule).Error; err != nil {
		return err
	}
	return nil
}

func (r *scheduleRepository) FindByTripID(ctx context.Context, tripID uuid.UUID) ([]domain.Schedule, error) {
	var schedules []domain.Schedule
	if err := r.db.WithContext(ctx).Where("tripId = ?", tripID).Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *scheduleRepository) FindByID(ctx context.Context, scheduleID uuid.UUID) (*domain.Schedule, error) {
	var schedule domain.Schedule
	if err := r.db.WithContext(ctx).First(&schedule, "id = ?", scheduleID).Error; err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *scheduleRepository) Update(ctx context.Context, schedule *domain.Schedule) error {
	if err := r.db.WithContext(ctx).Save(schedule).Error; err != nil {
		return err
	}
	return nil
}

func (r *scheduleRepository) Delete(ctx context.Context, scheduleID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&domain.Schedule{}, "id = ?", scheduleID).Error; err != nil {
		return err
	}
	return nil
}
