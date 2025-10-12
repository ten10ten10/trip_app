package usecase

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type ScheduleUsecaseValidator interface {
	ValidateCreateSchedule(startDateTime, endDateTime time.Time) error
}

type scheduleUsecaseValidator struct {
	validate *validator.Validate
}

func NewScheduleUsecaseValidator() ScheduleUsecaseValidator {
	return &scheduleUsecaseValidator{validate: validator.New()}
}

func (sv *scheduleUsecaseValidator) ValidateCreateSchedule(startDateTime, endDateTime time.Time) error {
	type createRequest struct {
		StartDateTime time.Time
		EndDateTime   time.Time `validate:"gtfield=StartDateTime"`
	}

	req := createRequest{
		StartDateTime: startDateTime,
		EndDateTime:   endDateTime,
	}

	return sv.validate.Struct(req)
}
