package handler

import (
	"errors"
	"net/http"
	"trip_app/api"
	"trip_app/internal/domain"
	"trip_app/internal/usecase"

	"github.com/labstack/echo/v4"
)

type publicScheduleHandler struct {
	su usecase.ScheduleUsecase
	sv ScheduleHandlerValidator
}

func NewPublicScheduleHandler(su usecase.ScheduleUsecase, sv ScheduleHandlerValidator) *publicScheduleHandler {
	return &publicScheduleHandler{su, sv}
}

// (POST /public/trips/{shareToken}/schedules)
func (h *publicScheduleHandler) AddScheduleToPublicTrip(ctx echo.Context, shareToken api.ShareToken) error {
	trip := ctx.Get("trip").(*domain.Trip)

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
		trip.ID,
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

// (GET /public/trips/{shareToken}/schedules)
func (h *publicScheduleHandler) GetSchedulesForPublicTrip(ctx echo.Context, shareToken api.ShareToken) error {
	trip := ctx.Get("trip").(*domain.Trip)

	schedules, err := h.su.GetSchedulesByTripID(ctx.Request().Context(), trip.ID)
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

// (GET /public/trips/{shareToken}/schedules/{scheduleId})
func (h *publicScheduleHandler) GetScheduleForPublicTrip(ctx echo.Context, shareToken api.ShareToken, scheduleId api.ScheduleId) error {
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

// (PATCH /public/trips/{shareToken}/schedules/{scheduleId})
func (h *publicScheduleHandler) UpdateScheduleForPublicTrip(ctx echo.Context, shareToken api.ShareToken, scheduleId api.ScheduleId) error {
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

// (DELETE /public/trips/{shareToken}/schedules/{scheduleId})
func (h *publicScheduleHandler) DeleteScheduleForPublicTrip(ctx echo.Context, shareToken api.ShareToken, scheduleId api.ScheduleId) error {
	if err := h.su.DeleteSchedule(ctx.Request().Context(), scheduleId); err != nil {
		if errors.Is(err, usecase.ErrScheduleNotFound) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"message": "Schedule not found"})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
	}

	return ctx.NoContent(http.StatusNoContent)
}
