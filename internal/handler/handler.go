package handler

import (
	"trip_app/api"
	"trip_app/internal/usecase"
)

// Handler holds all handlers
type Handler struct {
	*userHandler
	*tripHandler
}

func NewHandler(userUsecase usecase.UserUsecase, tripUsecase usecase.TripUsecase) api.ServerInterface {
	return &Handler{
		userHandler: NewUserHandler(userUsecase),
		tripHandler: NewTripHandler(tripUsecase),
	}
}
