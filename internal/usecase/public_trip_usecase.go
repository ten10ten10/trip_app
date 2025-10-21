package usecase

import (
	"context"
	"errors"
	"time"
	"trip_app/internal/domain"
	"trip_app/internal/repository"
	"trip_app/internal/security"

	"gorm.io/gorm"
)

type PublicTripUsecase interface {
	GetTripByShareToken(ctx context.Context, shareToken string) (*domain.Trip, error)
	UpdateTripByShareToken(ctx context.Context, shareToken string, title string, startDate, endDate time.Time, members []domain.Member) (*domain.Trip, error)
	GetTripDetailsByShareToken(ctx context.Context, shareToken string) (*domain.Trip, error)
}

type publicTripUsecase struct {
	pt repository.PublicTripRepository
	tg security.TokenGenerator
}

func NewPublicTripUsecase(pt repository.PublicTripRepository, tg security.TokenGenerator) PublicTripUsecase {
	return &publicTripUsecase{pt, tg}
}

func (pu *publicTripUsecase) GetTripByShareToken(ctx context.Context, shareToken string) (*domain.Trip, error) {
	tokenHash := pu.tg.HashToken(shareToken)
	trip, err := pu.pt.FindByShareToken(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTripNotFound
		}
		return nil, err
	}
	return trip, nil
}

func (pu *publicTripUsecase) UpdateTripByShareToken(ctx context.Context, shareToken string, title string, startDate, endDate time.Time, members []domain.Member) (*domain.Trip, error) {
	tokenHash := pu.tg.HashToken(shareToken)
	trip, err := pu.pt.FindByShareToken(ctx, tokenHash)
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

	if err := pu.pt.Update(ctx, trip); err != nil {
		return nil, err
	}

	return trip, nil
}

func (pu *publicTripUsecase) GetTripDetailsByShareToken(ctx context.Context, shareToken string) (*domain.Trip, error) {
	tokenHash := pu.tg.HashToken(shareToken)
	trip, err := pu.pt.FindWithSchedulesByShareToken(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTripNotFound
		}
		return nil, err
	}
	return trip, nil
}
