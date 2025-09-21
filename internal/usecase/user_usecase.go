package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/o_ten/trip_app/internal/domain"
	"github.com/o_ten/trip_app/internal/repository"
	"golang.org/x/crypto/bcrypt"
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
}

func NewUserUsecase(ur repository.UserRepository, uv validator.UserValidator) UserUsecase {
	return &userUsecase{ur, uv}
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
	initPassword := generateInitialPassword()
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(initPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// generate verification token
	Token, err := generateVerificationToken()
	if err != nil {
		return nil, err
	}
	tokenHash, err := hashToken(Token)
	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().Add(30 * time.Minute)

	// create user
	user := &domain.User{
		ID:                         uuid.New(),
		Name:                       name,
		Email:                      email,
		PasswordHash:               string(passwordHash),
		IsActive:                   false,
		VerificationTokenHash:      string(tokenHash),
		VerificationTokenExpiresAt: &expiresAt,
		CreatedAt:                  time.Now(),
		UpdatedAt:                  time.Now(),
	}
	if err := uu.ur.Create(ctx, user); err != nil {
		return nil, err
	}

	// send verification email
	if err := sendVerificationEmail(user.Email, initPassword, Token); err != nil {
		return nil, err
	}

	return user, nil
}
