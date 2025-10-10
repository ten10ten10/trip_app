package middleware

import (
	"trip_app/internal/security"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// JWTトークンを検証し、ユーザーIDをコンテキストに設定するEchoミドルウェアを生成
func AuthMiddleware(secret string) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		// JWT トークンの署名に使用するキー
		SigningKey: []byte(secret),
		// Claims の型を指定
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(security.JwtCustomClaims)
		},
		// SuccessHandlerはトークンが有効な場合に呼び出される関数
		// トークンからユーザーIDを抽出してコンテキストに保存
		SuccessHandler: func(c echo.Context) {
			if user, ok := c.Get("user").(*jwt.Token); ok {
				if claims, ok := user.Claims.(*security.JwtCustomClaims); ok {
					userID := claims.UserID
					parsedUserID, err := uuid.Parse(userID)
					if err == nil {
						c.Set("user_id", parsedUserID)
					}
				}
			}
		},
	})
}
