package handler

import (
	"errors"
	"net/http"

	"trip_app/api"
	"trip_app/internal/usecase"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type userHandler struct {
	uu usecase.UserUsecase
}

func NewUserHandler(uu usecase.UserUsecase) api.ServerInterface {
	return &userHandler{uu}
}

func (h *userHandler) CreateUser(ctx echo.Context) error {
	var req api.NewUser
	// Bind and validate request
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
	}

	// send request data to usecase from handler
	createdUser, err := h.uu.SignUp(ctx.Request().Context(), *req.Name, string(*req.Email))
	if err != nil {
		if errors.Is(err, usecase.ErrEmailConflict) {
			// if email already exists and active
			return ctx.JSON(http.StatusConflict, map[string]string{"message": err.Error()})
		}
		// other errors
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	// prepare response
	emailDTO := openapi_types.Email(createdUser.Email)
	res := api.User{
		Id:        &createdUser.ID,
		Name:      &createdUser.Name,
		Email:     &emailDTO,
		IsActive:  &createdUser.IsActive,
		CreatedAt: &createdUser.CreatedAt,
		UpdatedAt: &createdUser.UpdatedAt,
	}

	return ctx.JSON(http.StatusCreated, res)
}

func (h *userHandler) VerifyUser(ctx echo.Context) error {
	// get verificationToken from path parameter
	verificationToken := ctx.Param("verificationToken")

	// call usecase to verify email
	message, err := h.uu.VerifyEmail(ctx.Request().Context(), verificationToken)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidVerificationToken) || errors.Is(err, usecase.ErrVerificationTokenExpired) {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": message})
}

func (h *userHandler) Login(ctx echo.Context) error {
	var req api.LoginRequest

	// Bind and validate request
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
	}

	// send request data to usecase from handler
	user, token, err := h.uu.Login(ctx.Request().Context(), string(req.Email), req.Password)
	if err != nil {
		if errors.Is(err, usecase.ErrValidation) {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}
		if errors.Is(err, usecase.ErrInvalidCredentials) || errors.Is(err, usecase.ErrUserNotActive) {
			return ctx.JSON(http.StatusUnauthorized, map[string]string{"message": err.Error()})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
	}

	// prepare response
	emailDTO := openapi_types.Email(user.Email)
	userResponse := api.User{
		Id:        &user.ID,
		Name:      &user.Name,
		Email:     &emailDTO,
		IsActive:  &user.IsActive,
		CreatedAt: &user.CreatedAt,
		UpdatedAt: &user.UpdatedAt,
	}

	res := api.AuthResponse{
		User:  &userResponse,
		Token: &token,
	}

	return ctx.JSON(http.StatusOK, res)
}

func (h *userHandler) Logout(ctx echo.Context) error {
	// currently, client-side just deletes the token, so nothing to do server-side
	// in the future, we might want to implement token blacklisting or expiration by redis
	return nil
}

func (h *userHandler) GetMe(ctx echo.Context) error {
	userID, ok := ctx.Get("user_id").(uuid.UUID)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
	}

	user, err := h.uu.GetProfile(ctx.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"message": err.Error()})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
	}

	// prepare response
	emailDTO := openapi_types.Email(user.Email)
	res := api.User{
		Id:        &user.ID,
		Name:      &user.Name,
		Email:     &emailDTO,
		IsActive:  &user.IsActive,
		CreatedAt: &user.CreatedAt,
		UpdatedAt: &user.UpdatedAt,
	}

	return ctx.JSON(http.StatusOK, res)
}

func (h *userHandler) ChangePassword(ctx echo.Context) error {
	var req api.PasswordChangeRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
	}

	userID, ok := ctx.Get("user_id").(uuid.UUID)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
	}

	err := h.uu.ChangePassword(ctx.Request().Context(), userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		if errors.Is(err, usecase.ErrValidation) {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}
		if errors.Is(err, usecase.ErrUserNotFound) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"message": err.Error()})
		}
		if errors.Is(err, usecase.ErrIncorrectCurrentPassword) {
			return ctx.JSON(http.StatusUnauthorized, map[string]string{"message": err.Error()})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
	}

	return ctx.JSON(http.StatusNoContent, map[string]string{"message": "Password changed successfully"})
}
