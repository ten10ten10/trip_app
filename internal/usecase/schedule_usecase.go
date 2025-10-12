package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"
	"trip_app/internal/domain"
	"trip_app/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrScheduleNotFound = errors.New("schedule not found")
var ErrValidation = errors.New("input validation failed")

type UpdateScheduleParams struct {
	Title         *string
	StartDateTime *time.Time
	EndDateTime   *time.Time
	Memo          *string
}

type ScheduleUsecase interface {
	CreateSchedule(ctx context.Context, tripID uuid.UUID, title string, startDateTime, endDateTime time.Time, memo string) (*domain.Schedule, error)
	GetSchedulesByTripID(ctx context.Context, tripID uuid.UUID) ([]domain.Schedule, error)
	GetScheduleByID(ctx context.Context, scheduleID uuid.UUID) (*domain.Schedule, error)
	UpdateSchedule(ctx context.Context, scheduleID uuid.UUID, params UpdateScheduleParams) (*domain.Schedule, error)
	DeleteSchedule(ctx context.Context, scheduleID uuid.UUID) error
}

type scheduleUsecase struct {
	sr repository.ScheduleRepository
	sv ScheduleUsecaseValidator
}

func NewScheduleUsecase(sr repository.ScheduleRepository, sv ScheduleUsecaseValidator) ScheduleUsecase {
	return &scheduleUsecase{sr, sv}
}

func (su *scheduleUsecase) CreateSchedule(ctx context.Context, tripID uuid.UUID, title string, startDateTime, endDateTime time.Time, memo string) (*domain.Schedule, error) {
	if err := su.sv.ValidateCreateSchedule(startDateTime, endDateTime); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidation, err)
	}

	schedule := &domain.Schedule{
		TripID:        tripID,
		Title:         title,
		StartDateTime: startDateTime,
		EndDateTime:   endDateTime,
		Memo:          memo,
	}

	if err := su.sr.Create(ctx, schedule); err != nil {
		return nil, err
	}

	return schedule, nil
}

func (su *scheduleUsecase) GetSchedulesByTripID(ctx context.Context, tripID uuid.UUID) ([]domain.Schedule, error) {
	schedules, err := su.sr.FindByTripID(ctx, tripID)
	if err != nil {
		return nil, err
	}
	return schedules, nil
}

func (su *scheduleUsecase) GetScheduleByID(ctx context.Context, scheduleID uuid.UUID) (*domain.Schedule, error) {
	schedule, err := su.sr.FindByID(ctx, scheduleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrScheduleNotFound
		}
		return nil, err
	}
	return schedule, nil
}

func (su *scheduleUsecase) UpdateSchedule(ctx context.Context, scheduleID uuid.UUID, params UpdateScheduleParams) (*domain.Schedule, error) {
	schedule, err := su.sr.FindByID(ctx, scheduleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrScheduleNotFound
		}
		return nil, err
	}

	// Determine the final values for start and end times for validation
	newStart := schedule.StartDateTime
	if params.StartDateTime != nil {
		newStart = *params.StartDateTime
	}
	newEnd := schedule.EndDateTime
	if params.EndDateTime != nil {
		newEnd = *params.EndDateTime
	}

	// Validate the relationship of the final values
	if err := su.sv.ValidateCreateSchedule(newStart, newEnd); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidation, err)
	}

	// Update fields if new values are provided
	if params.Title != nil {
		schedule.Title = *params.Title
	}
	if params.StartDateTime != nil {
		schedule.StartDateTime = *params.StartDateTime
	}
	if params.EndDateTime != nil {
		schedule.EndDateTime = *params.EndDateTime
	}
	if params.Memo != nil {
		schedule.Memo = *params.Memo
	}

	if err := su.sr.Update(ctx, schedule); err != nil {
		return nil, err
	}

	return schedule, nil
}

func (su *scheduleUsecase) DeleteSchedule(ctx context.Context, scheduleID uuid.UUID) error {
	_, err := su.sr.FindByID(ctx, scheduleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrScheduleNotFound
		}
		return err
	}
	if err := su.sr.Delete(ctx, scheduleID); err != nil {
		return err
	}
	return nil
}
