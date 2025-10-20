package repository

import (
	"context"

	"gorm.io/gorm"

	"trip_app/internal/domain"
)

type ShareTokenRepository interface {
	Create(ctx context.Context, shareToken *domain.ShareToken) error
	Update(ctx context.Context, shareToken *domain.ShareToken) error
}

type shareTokenRepository struct {
	db *gorm.DB
}

func NewShareTokenRepository(db *gorm.DB) ShareTokenRepository {
	return &shareTokenRepository{db}
}

func (r *shareTokenRepository) Create(ctx context.Context, shareToken *domain.ShareToken) error {
	return r.db.WithContext(ctx).Create(shareToken).Error
}

func (r *shareTokenRepository) Update(ctx context.Context, shareToken *domain.ShareToken) error {
	return r.db.WithContext(ctx).Save(shareToken).Error
}
