package validator

import (
	"github.com/go-playground/validator/v10"
)

type UserValidator interface {
	ValidateSignUp(name, email string) error
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
