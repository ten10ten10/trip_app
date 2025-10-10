package handler

import (
	"errors"
	"net/http"
	"trip_app/api"
	"trip_app/internal/domain"
	"trip_app/internal/usecase"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type tripHandler struct {
	tu usecase.TripUsecase
}

func NewTripHandler(tu usecase.TripUsecase) *tripHandler {
	return &tripHandler{tu}
}

// --- Model Conversion Helper Functions ---

func toAPIMembers(members []domain.Member) *[]api.Member {
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

func toAPITrip(trip *domain.Trip) *api.Trip {
	if trip == nil {
		return nil
	}
	return &api.Trip{
		Id:        &trip.ID,
		Title:     &trip.Title,
		StartDate: &openapi_types.Date{Time: trip.StartDate},
		EndDate:   &openapi_types.Date{Time: trip.EndDate},
		Members:   toAPIMembers(trip.Members),
		CreatedAt: &trip.CreatedAt,
		UpdatedAt: &trip.UpdatedAt,
	}
}

func toAPITrips(trips []domain.Trip) []api.Trip {
	apiTrips := make([]api.Trip, len(trips))
	for i, t := range trips {
		apiTrips[i] = *toAPITrip(&t)
	}
	return apiTrips
}

func toAPISchedules(schedules []domain.Schedule) *[]api.Schedule {
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

func (h *tripHandler) CreateUserTrip(ctx echo.Context) error {
	userID, ok := ctx.Get("user_id").(uuid.UUID)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
	}

	var req api.NewTripRequest
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

	createdTrip, err := h.tu.CreateTrip(
		ctx.Request().Context(),
		userID,
		req.Title,
		req.StartDate.Time,
		req.EndDate.Time,
		members,
	)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create trip"})
	}

	return ctx.JSON(http.StatusCreated, toAPITrip(createdTrip))
}

func (h *tripHandler) GetUserTrips(ctx echo.Context) error {
	userID, ok := ctx.Get("user_id").(uuid.UUID)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
	}

	trips, err := h.tu.GetTripsByUserID(ctx.Request().Context(), userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
	}

	return ctx.JSON(http.StatusOK, toAPITrips(trips))
}

func (h *tripHandler) GetUserTrip(ctx echo.Context, tripId api.TripId) error {
	trip := ctx.Get("trip").(*domain.Trip)
	return ctx.JSON(http.StatusOK, toAPITrip(trip))
}

func (h *tripHandler) UpdateUserTrip(ctx echo.Context, tripId api.TripId) error {
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

	updatedTrip, err := h.tu.UpdateTrip(
		ctx.Request().Context(),
		tripId,
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

	return ctx.JSON(http.StatusOK, toAPITrip(updatedTrip))
}

func (h *tripHandler) GetTripDetails(ctx echo.Context, tripId api.TripId) error {
	tripWithSchedules, err := h.tu.GetTripDetailsByID(ctx.Request().Context(), tripId)
	if err != nil {
		if errors.Is(err, usecase.ErrTripNotFound) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"message": "Trip not found"})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
	}

	res := api.TripDetailView{
		Trip:      toAPITrip(tripWithSchedules),
		Schedules: toAPISchedules(tripWithSchedules.Schedules),
	}

	return ctx.JSON(http.StatusOK, res)
}

func (h *tripHandler) DeleteUserTrip(ctx echo.Context, tripId api.TripId) error {
	if err := h.tu.DeleteTrip(ctx.Request().Context(), tripId); err != nil {
		if errors.Is(err, usecase.ErrTripNotFound) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"message": "Trip not found"})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
	}

	return ctx.NoContent(http.StatusNoContent)
}
