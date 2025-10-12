package handler

import (
	"trip_app/api"
	"trip_app/internal/usecase"
)

// Handler holds all handlers
type Handler struct {
	*userHandler
	*tripHandler
	*scheduleHandler
}

func NewHandler(
	userUsecase usecase.UserUsecase,
	tripUsecase usecase.TripUsecase,
	scheduleUsecase usecase.ScheduleUsecase,
	userHandlerValidator UserHandlerValidator,
	scheduleHandlerValidator ScheduleHandlerValidator,
) api.ServerInterface {
	return &Handler{
		userHandler:     NewUserHandler(userUsecase, userHandlerValidator),
		tripHandler:     NewTripHandler(tripUsecase),
		scheduleHandler: NewScheduleHandler(scheduleUsecase, scheduleHandlerValidator),
	}
}
