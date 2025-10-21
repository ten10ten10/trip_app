package handler

import (
	"errors"
	"net/http"
	"trip_app/api"
	"trip_app/internal/domain"
	"trip_app/internal/usecase"

	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type publicTripHandler struct {
	ptu usecase.PublicTripUsecase
}

func NewPublicTripHandler(ptu usecase.PublicTripUsecase) *publicTripHandler {
	return &publicTripHandler{ptu}
}

// --- Model Conversion Helper Functions ---

func toAPIPublicTrip(trip *domain.Trip) *api.Trip {
	if trip == nil {
		return nil
	}
	return &api.Trip{
		Id:        &trip.ID,
		Title:     &trip.Title,
		StartDate: &openapi_types.Date{Time: trip.StartDate},
		EndDate:   &openapi_types.Date{Time: trip.EndDate},
		Members:   toAPIPublicMembers(trip.Members),
		CreatedAt: &trip.CreatedAt,
		UpdatedAt: &trip.UpdatedAt,
	}
}

func toAPIPublicMembers(members []domain.Member) *[]api.Member {
	if members == nil {
		return nil
	}
	apiMembers := make([]api.Member, len(members))
	for i, m := range members {
		apiMembers[i] = api.Member{
			Id:   &m.ID,
			Name: &m.Name,
		}
	}
	return &apiMembers
}

func toAPIPublicSchedules(schedules []domain.Schedule) *[]api.Schedule {
	if schedules == nil {
		return nil
	}
	apiSchedules := make([]api.Schedule, len(schedules))
	for i, s := range schedules {
		apiSchedules[i] = api.Schedule{
			Id:            &s.ID,
			Title:         &s.Title,
			StartDateTime: &s.StartDateTime,
			EndDateTime:   &s.EndDateTime,
			Memo:          &s.Memo,
			CreatedAt:     &s.CreatedAt,
			UpdatedAt:     &s.UpdatedAt,
		}
	}
	return &apiSchedules
}

// --- Handlers ---

func (h *publicTripHandler) GetPublicTripByShareToken(ctx echo.Context, shareToken api.ShareToken) error {
	trip := ctx.Get("trip").(*domain.Trip)
	return ctx.JSON(http.StatusOK, toAPIPublicTrip(trip))
}

func (h *publicTripHandler) UpdatePublicTripByShareToken(ctx echo.Context, shareToken api.ShareToken) error {
	var req api.UpdateTripRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
	}

	var members []domain.Member
	if req.Members != nil {
		for _, m := range *req.Members {
			if m.Name != nil {
				members = append(members, domain.Member{Name: *m.Name})
			}
		}
	}

	updatedTrip, err := h.ptu.UpdateTripByShareToken(
		ctx.Request().Context(),
		shareToken,
		req.Title,
		req.StartDate.Time,
		req.EndDate.Time,
		members,
	)
	if err != nil {
		if errors.Is(err, usecase.ErrTripNotFound) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"message": "Trip not found"})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
	}

	return ctx.JSON(http.StatusOK, toAPIPublicTrip(updatedTrip))
}

func (h *publicTripHandler) GetTripDetailsForPublicTrip(ctx echo.Context, shareToken api.ShareToken) error {
	tripWithSchedules, err := h.ptu.GetTripDetailsByShareToken(ctx.Request().Context(), shareToken)
	if err != nil {
		if errors.Is(err, usecase.ErrTripNotFound) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"message": "Trip not found"})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
	}

	res := api.TripDetailView{
		Trip:      toAPIPublicTrip(tripWithSchedules),
		Schedules: toAPIPublicSchedules(tripWithSchedules.Schedules),
	}

	return ctx.JSON(http.StatusOK, res)
}
