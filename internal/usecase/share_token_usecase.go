package usecase

import (
	"context"
	"time"
	"trip_app/internal/domain"
	"trip_app/internal/repository"
	"trip_app/internal/security"

	"github.com/google/uuid"
)

type ShareTokenUsecase interface {
	CreateShareToken(ctx context.Context, tripID uuid.UUID, regenerate bool) (*domain.ShareToken, string, error)
}

type shareTokenUsecase struct {
	ur repository.ShareTokenRepository
	us security.TokenGenerator
}

func NewShareTokenUsecase(ur repository.ShareTokenRepository, us security.TokenGenerator) ShareTokenUsecase {
	return &shareTokenUsecase{ur, us}
}

func (u *shareTokenUsecase) CreateShareToken(ctx context.Context, tripID uuid.UUID, regenerate bool) (*domain.ShareToken, string, error) {
	rawShareToken, hashToken, err := u.us.GenerateToken()
	if err != nil {
		return nil, "", err
	}

	shareToken := &domain.ShareToken{
		TripID:    tripID,
		TokenHash: hashToken,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if regenerate {
		err = u.ur.Update(ctx, shareToken)
	} else {
		err = u.ur.Create(ctx, shareToken)
	}

	if err != nil {
		return nil, "", err
	}

	return shareToken, rawShareToken, nil
}
