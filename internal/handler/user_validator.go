package handler

import (
	"github.com/go-playground/validator/v10"
)

type UserHandlerValidator interface {
	ValidateSignUp(name, email string) error
	ValidateLogin(email, password string) error
	ValidateChangePassword(currentPassword, newPassword string) error
}

type userHandlerValidator struct {
	validate *validator.Validate
}

func NewUserHandlerValidator() UserHandlerValidator {
	return &userHandlerValidator{validate: validator.New()}
}

func (uv *userHandlerValidator) ValidateSignUp(name, email string) error {
	type signUpRequest struct {
		Name  string `validate:"required"`
		Email string `validate:"required,email"`
	}
	req := signUpRequest{Name: name, Email: email}
	return uv.validate.Struct(req)
}

func (uv *userHandlerValidator) ValidateLogin(email, password string) error {
	type loginRequest struct {
		Email    string `validate:"required,email"`
		Password string `validate:"required,min=8"`
	}
	req := loginRequest{Email: email, Password: password}
	return uv.validate.Struct(req)
}

func (uv *userHandlerValidator) ValidateChangePassword(currentPassword, newPassword string) error {
	type changePasswordRequest struct {
		CurrentPassword string `validate:"required,min=8"`
		NewPassword     string `validate:"required,min=8"`
	}
	req := changePasswordRequest{CurrentPassword: currentPassword, NewPassword: newPassword}
	return uv.validate.Struct(req)
}
