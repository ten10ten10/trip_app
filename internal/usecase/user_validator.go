package usecase

import (
	"github.com/go-playground/validator/v10"
)

type UserUsecaseValidator interface {
	ValidateChangePassword(currentPassword, newPassword string) error
}

type userUsecaseValidator struct {
	validate *validator.Validate
}

func NewUserUsecaseValidator() UserUsecaseValidator {
	return &userUsecaseValidator{validate: validator.New()}
}

func (uv *userUsecaseValidator) ValidateChangePassword(currentPassword, newPassword string) error {
	type changePasswordRequest struct {
		CurrentPassword string
		NewPassword     string `validate:"nefield=CurrentPassword"`
	}
	req := changePasswordRequest{CurrentPassword: currentPassword, NewPassword: newPassword}
	return uv.validate.Struct(req)
}
