package security

import (
	"time"

	"trip_app/internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type AuthTokenGenerator interface {
	GenerateAccessToken(user *domain.User) (string, error)
}

type jwtGenerator struct {
	jwtSecret string
}

func NewAuthTokenGenerator(jwtSecret string) AuthTokenGenerator {
	return &jwtGenerator{jwtSecret: jwtSecret}
}

func (g *jwtGenerator) GenerateAccessToken(user *domain.User) (string, error) {
	// jwtに埋め込むデータ(クレーム)を作成
	claims := &JwtCustomClaims{
		UserID: user.ID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(g.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
