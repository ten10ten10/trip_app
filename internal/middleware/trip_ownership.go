package middleware

import (
	"errors"
	"net/http"
	"trip_app/internal/usecase"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func TripOwnershipMiddleware(tripUsecase usecase.TripUsecase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 認証済みユーザーIDを取得
			userID, ok := c.Get("user_id").(uuid.UUID)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
			}

			// URLからtripIdを取得
			tripIDStr := c.Param("tripId")
			tripID, err := uuid.Parse(tripIDStr)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid trip ID"})
			}

			// Usecaseを使って旅行情報を取得
			trip, err := tripUsecase.GetTripByTripID(c.Request().Context(), tripID)
			if err != nil {
				if errors.Is(err, usecase.ErrTripNotFound) {
					return c.JSON(http.StatusNotFound, map[string]string{"message": "Trip not found"})
				}
				return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
			}

			// 所有者かどうかチェック
			if trip.UserID != userID {
				return c.JSON(http.StatusForbidden, map[string]string{"message": "Forbidden"})
			}

			// 取得した旅行情報をctxに保存
			c.Set("trip", trip)

			// handlerへ処理を渡す
			return next(c)
		}
	}
}
