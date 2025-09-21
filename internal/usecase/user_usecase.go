package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/o_ten/trip_app/internal/domain"
	"github.com/o_ten/trip_app/internal/infrastructure/email"
	"github.com/o_ten/trip_app/internal/repository"
	"github.com/o_ten/trip_app/internal/security"
	"github.com/o_ten/trip_app/internal/validator"
	"gorm.io/gorm"
)

type UserUsecase interface {
	SignUp(ctx context.Context, name, email string) (*domain.User, error)
	VerifyEmail(ctx context.Context, token string) error
	Login(ctx context.Context, email, password string) (*domain.User, error)
	Logout(ctx context.Context, userID uuid.UUID) error
	GetProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) error
}

type userUsecase struct {
	ur repository.UserRepository
	uv validator.UserValidator
	up security.PasswordGenerator
	us security.TokenGenerator
	ue email.Sender
}

func NewUserUsecase(ur repository.UserRepository, uv validator.UserValidator, up security.PasswordGenerator, us security.TokenGenerator, ue email.Sender) UserUsecase {
	return &userUsecase{ur, uv, up, us, ue}
}

func (uu *userUsecase) SignUp(ctx context.Context, name, email string) (*domain.User, error) {
	// validate input
	if err := uu.uv.ValidateSignUp(name, email); err != nil {
		return nil, err
	}

	// check if email already exists
	_, err := uu.ur.FindByEmail(ctx, email)
	if err == nil {
		// find no error means email already exists
		return nil, errors.New("email already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// some other error occurred
		return nil, err
	}

	// initPassword & hashPassword
	rawPassword, hashPassword, err := uu.up.GeneratePassword()
	if err != nil {
		return nil, err
	}

	// generate verification token
	rawToken, hashToken, err := uu.us.GenerateToken()
	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().Add(30 * time.Minute)

	// create user
	user := &domain.User{
		ID:                         uuid.New(),
		Name:                       name,
		Email:                      email,
		PasswordHash:               string(hashPassword),
		IsActive:                   false,
		VerificationTokenHash:      &hashToken,
		VerificationTokenExpiresAt: &expiresAt,
		CreatedAt:                  time.Now(),
		UpdatedAt:                  time.Now(),
	}
	if err := uu.ur.Create(ctx, user); err != nil {
		return nil, err
	}

	// send verification email
	if err := uu.ue.SendVerificationEmail(ctx, user.Email, rawToken, rawPassword); err != nil {
		return nil, err
	}

	return user, nil
}
