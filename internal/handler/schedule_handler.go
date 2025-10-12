package handler

import (
	"errors"
	"net/http"
	"trip_app/api"
	"trip_app/internal/usecase"

	"github.com/labstack/echo/v4"
)

type scheduleHandler struct {
	su usecase.ScheduleUsecase
	sv ScheduleHandlerValidator
}

func NewScheduleHandler(su usecase.ScheduleUsecase, sv ScheduleHandlerValidator) *scheduleHandler {
	return &scheduleHandler{su, sv}
}

// (POST /trips/{tripId}/schedules)
func (h *scheduleHandler) AddScheduleToTrip(ctx echo.Context, tripId api.TripId) error {
	var req api.NewSchedule
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
	}

	if err := h.sv.ValidateAddSchedule(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	var memo string
	if req.Memo != nil {
		memo = *req.Memo
	}

	createdSchedule, err := h.su.CreateSchedule(
		ctx.Request().Context(),
		tripId,
		*req.Title,
		*req.StartDateTime,
		*req.EndDateTime,
		memo,
	)
	if err != nil {
		if errors.Is(err, usecase.ErrValidation) {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
	}

	// domain.Scheduleからapi.Scheduleへの変換
	res := api.Schedule{
		Id:            &createdSchedule.ID,
		Title:         &createdSchedule.Title,
		StartDateTime: &createdSchedule.StartDateTime,
		EndDateTime:   &createdSchedule.EndDateTime,
		Memo:          &createdSchedule.Memo,
		CreatedAt:     &createdSchedule.CreatedAt,
		UpdatedAt:     &createdSchedule.UpdatedAt,
	}

	return ctx.JSON(http.StatusCreated, res)
}

// (GET /trips/{tripId}/schedules)
func (h *scheduleHandler) GetSchedulesForTrip(ctx echo.Context, tripId api.TripId) error {
	schedules, err := h.su.GetSchedulesByTripID(ctx.Request().Context(), tripId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
	}

	res := make([]api.Schedule, len(schedules))
	for i, schedule := range schedules {
		res[i] = api.Schedule{
			Id:            &schedule.ID,
			Title:         &schedule.Title,
			StartDateTime: &schedule.StartDateTime,
			EndDateTime:   &schedule.EndDateTime,
			Memo:          &schedule.Memo,
			CreatedAt:     &schedule.CreatedAt,
			UpdatedAt:     &schedule.UpdatedAt,
		}
	}

	return ctx.JSON(http.StatusOK, res)
}

// (GET /trips/{tripId}/schedules/{scheduleId})
func (h *scheduleHandler) GetScheduleForTrip(ctx echo.Context, tripId api.TripId, scheduleId api.ScheduleId) error {
	schedule, err := h.su.GetScheduleByID(ctx.Request().Context(), scheduleId)
	if err != nil {
		if errors.Is(err, usecase.ErrScheduleNotFound) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"message": "Schedule not found"})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
	}

	res := api.Schedule{
		Id:            &schedule.ID,
		Title:         &schedule.Title,
		StartDateTime: &schedule.StartDateTime,
		EndDateTime:   &schedule.EndDateTime,
		Memo:          &schedule.Memo,
		CreatedAt:     &schedule.CreatedAt,
		UpdatedAt:     &schedule.UpdatedAt,
	}

	return ctx.JSON(http.StatusOK, res)
}

// (PATCH /trips/{tripId}/schedules/{scheduleId})
func (h *scheduleHandler) UpdateScheduleForTrip(ctx echo.Context, tripId api.TripId, scheduleId api.ScheduleId) error {
	var req api.UpdateSchedule
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
	}

	updatedSchedule, err := h.su.UpdateSchedule(
		ctx.Request().Context(),
		scheduleId,
		usecase.UpdateScheduleParams{
			Title:         req.Title,
			StartDateTime: req.StartDateTime,
			EndDateTime:   req.EndDateTime,
			Memo:          req.Memo,
		},
	)
	if err != nil {
		if errors.Is(err, usecase.ErrScheduleNotFound) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"message": "Schedule not found"})
		}
		if errors.Is(err, usecase.ErrValidation) {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
	}

	res := api.Schedule{
		Id:            &updatedSchedule.ID,
		Title:         &updatedSchedule.Title,
		StartDateTime: &updatedSchedule.StartDateTime,
		EndDateTime:   &updatedSchedule.EndDateTime,
		Memo:          &updatedSchedule.Memo,
		CreatedAt:     &updatedSchedule.CreatedAt,
		UpdatedAt:     &updatedSchedule.UpdatedAt,
	}

	return ctx.JSON(http.StatusOK, res)
}

// (DELETE /trips/{tripId}/schedules/{scheduleId})
func (h *scheduleHandler) DeleteScheduleForTrip(ctx echo.Context, tripId api.TripId, scheduleId api.ScheduleId) error {
	if err := h.su.DeleteSchedule(ctx.Request().Context(), scheduleId); err != nil {
		if errors.Is(err, usecase.ErrScheduleNotFound) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"message": "Schedule not found"})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
	}

	return ctx.NoContent(http.StatusNoContent)
}
