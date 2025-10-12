package handler

import (
	"time"
	"trip_app/api"

	"github.com/go-playground/validator/v10"
)

type ScheduleHandlerValidator interface {
	ValidateAddSchedule(req api.NewSchedule) error
}

type scheduleHandlerValidator struct {
	validate *validator.Validate
}

func NewScheduleHandlerValidator() ScheduleHandlerValidator {
	return &scheduleHandlerValidator{validate: validator.New()}
}

func (sv *scheduleHandlerValidator) ValidateAddSchedule(req api.NewSchedule) error {
	// バリデーション用の構造体を定義
	type addScheduleRequest struct {
		Title         *string    `validate:"required"`
		StartDateTime *time.Time `validate:"required"`
		EndDateTime   *time.Time `validate:"required"`
	}

	// リクエストをバリデーション用構造体にマッピング
	validateReq := addScheduleRequest{
		Title:         req.Title,
		StartDateTime: req.StartDateTime,
		EndDateTime:   req.EndDateTime,
	}

	return sv.validate.Struct(validateReq)
}
