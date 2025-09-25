package validator

import (
	"github.com/go-playground/validator/v10"
)

type UserValidator interface {
	ValidateSignUp(name, email string) error
	ValidateLogin(email, password string) error
	ValidateChangePassword(currentPassword, newPassword string) error
}

type userValidator struct {
	validate *validator.Validate
}

func NewUserValidator() UserValidator {
	return &userValidator{validate: validator.New()}
}

func (uv *userValidator) ValidateSignUp(name, email string) error {
	type signUpRequest struct {
		Name  string `validate:"required"`
		Email string `validate:"required,email"`
	}

	req := signUpRequest{
		Name:  name,
		Email: email,
	}

	return uv.validate.Struct(req)
}

func (uv *userValidator) ValidateLogin(email, password string) error {
	type LoginRequest struct {
		Email    string `validate:"required,email"`
		Password string `validate:"required,min=8"`
	}

	req := LoginRequest{
		Email:    email,
		Password: password,
	}

	return uv.validate.Struct(req)
}

func (uv *userValidator) ValidateChangePassword(currentPassword, newPassword string) error {
	type ChangePasswordRequest struct {
		CurrentPassword string `validate:"required,min=8"`
		NewPassword     string `validate:"required,min=8,nefield=CurrentPassword"`
	}

	req := ChangePasswordRequest{
		CurrentPassword: currentPassword,
		NewPassword:     newPassword,
	}

	return uv.validate.Struct(req)
}
