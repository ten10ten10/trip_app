package handler

import (
	"net/http"

	"trip_app/api"
	"trip_app/internal/domain"
	"trip_app/internal/usecase"

	"github.com/labstack/echo/v4"
)

type ShareTokenHandler interface {
	CreateShareToken(ctx echo.Context) error
}

type shareTokenHandler struct {
	u usecase.ShareTokenUsecase
}

func NewShareTokenHandler(u usecase.ShareTokenUsecase) ShareTokenHandler {
	return &shareTokenHandler{u}
}

func (h *shareTokenHandler) CreateShareToken(ctx echo.Context) error {
	trip, ok := ctx.Get("trip").(*domain.Trip)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get trip from context")
	}

	regenerateStr := ctx.QueryParam("regenerate")
	regenerate := regenerateStr == "true"

	shareToken, token, err := h.u.CreateShareToken(ctx.Request().Context(), trip.ID, regenerate)
	if err != nil {
		return err
	}

	shareUrl := "/public/trips/" + token

	res := api.ShareLinkResponse{
		ShareToken: &token,
		ShareUrl:   &shareUrl,
		CreatedAt:  &shareToken.CreatedAt,
		UpdatedAt:  &shareToken.UpdatedAt,
	}

	return ctx.JSON(http.StatusCreated, res)
}
