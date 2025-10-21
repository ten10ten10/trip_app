package middleware

import (
	"errors"
	"net/http"
	"trip_app/internal/usecase"

	"github.com/labstack/echo/v4"
)

func ShareTokenOwnershipMiddleware(publicTripUsecase usecase.PublicTripUsecase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// URLからshareTokenを取得
			shareToken := c.Param("shareToken")
			if shareToken == "" {
				return c.JSON(http.StatusBadRequest, map[string]string{"message": "Share token is required"})
			}

			// Usecaseを使ってshareTokenから旅行情報を取得
			trip, err := publicTripUsecase.GetTripByShareToken(c.Request().Context(), shareToken)
			if err != nil {
				if errors.Is(err, usecase.ErrTripNotFound) {
					return c.JSON(http.StatusNotFound, map[string]string{"message": "Trip not found"})
				}
				return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
			}

			// 取得した旅行情報をctxに保存
			c.Set("trip", trip)

			// handlerへ処理を渡す
			return next(c)
		}
	}
}
