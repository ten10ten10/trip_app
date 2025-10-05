package usecase

import (
	"context"
	"errors"
	"time"
	"trip_app/internal/domain"
	"trip_app/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrTripNotFound = errors.New("trip not found")

type TripUsecase interface {
	CreateTrip(ctx context.Context, userID uuid.UUID, title string, startDate, endDate time.Time, members []domain.Member) (*domain.Trip, error)
	GetTripsByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Trip, error)
	GetTripByTripID(ctx context.Context, tripID uuid.UUID) (*domain.Trip, error)
	GetTripByShareToken(ctx context.Context, shareToken string) (*domain.Trip, error)
	UpdateTrip(ctx context.Context, tripID uuid.UUID, title string, startDate, endDate time.Time, members []domain.Member) (*domain.Trip, error)
	GetTripDetailsByID(ctx context.Context, tripID uuid.UUID) (*domain.Trip, error)
	DeleteTrip(ctx context.Context, tripID uuid.UUID) error
}

type tripUsecase struct {
	tr repository.TripRepository
}

func NewTripUsecase(tr repository.TripRepository) TripUsecase {
	return &tripUsecase{tr}
}

func (tu *tripUsecase) CreateTrip(ctx context.Context, userID uuid.UUID, title string, startDate, endDate time.Time, members []domain.Member) (*domain.Trip, error) {
	trip := &domain.Trip{
		UserID:    userID,
		Title:     title,
		StartDate: startDate,
		EndDate:   endDate,
		Members:   members,
	}

	if err := tu.tr.Create(ctx, trip); err != nil {
		return nil, err
	}

	return trip, nil
}

func (tu *tripUsecase) GetTripsByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Trip, error) {
	trips, err := tu.tr.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return trips, nil
}

func (tu *tripUsecase) GetTripByTripID(ctx context.Context, tripID uuid.UUID) (*domain.Trip, error) {
	trip, err := tu.tr.FindByID(ctx, tripID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTripNotFound
		}
		return nil, err
	}
	return trip, nil
}

func (tu *tripUsecase) GetTripByShareToken(ctx context.Context, shareToken string) (*domain.Trip, error) {
	trip, err := tu.tr.FindByShareToken(ctx, shareToken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTripNotFound
		}
		return nil, err
	}
	return trip, nil
}

func (tu *tripUsecase) UpdateTrip(ctx context.Context, tripID uuid.UUID, title string, startDate, endDate time.Time, members []domain.Member) (*domain.Trip, error) {
	trip, err := tu.tr.FindByID(ctx, tripID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTripNotFound
		}
		return nil, err
	}

	trip.Title = title
	trip.StartDate = startDate
	trip.EndDate = endDate
	trip.Members = members

	if err := tu.tr.Update(ctx, trip); err != nil {
		return nil, err
	}

	return trip, nil
}

func (tu *tripUsecase) GetTripDetailsByID(ctx context.Context, tripID uuid.UUID) (*domain.Trip, error) {
	trip, err := tu.tr.FindWithSchedulesByID(ctx, tripID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTripNotFound
		}
		return nil, err
	}
	return trip, nil
}

func (tu *tripUsecase) DeleteTrip(ctx context.Context, tripID uuid.UUID) error {
	_, err := tu.tr.FindByID(ctx, tripID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTripNotFound
		}
		return err
	}

	if err := tu.tr.Delete(ctx, tripID); err != nil {
		return err
	}

	return nil
}
