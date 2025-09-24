package handler

import (
	"errors"
	"net/http"

	"trip_app/api"
	"trip_app/internal/usecase"

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
		// handle errors
		if err.Error() == "invalid verification token" || err.Error() == "verification token has expired" {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": message})
}
