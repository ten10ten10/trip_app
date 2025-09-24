package security

import (
	"time"

	"trip_app/internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

type AuthTokenGenerator interface {
	GenerateAccessToken(user *domain.User) (string, error)
}

// jwtSecretは.envなどで安全に管理
type jwtGenerator struct {
	jwtSecret string
}

func NewAuthTokenGenerator(jwtSecret string) AuthTokenGenerator {
	return &jwtGenerator{jwtSecret: jwtSecret}
}

func (g *jwtGenerator) GenerateAccessToken(user *domain.User) (string, error) {
	// jwtに埋め込むデータ(クレーム)を作成
	// ここではユーザーIDと有効期限を設定
	claims := jwt.MapClaims{
		"user_id":    user.ID.String(),
		"expired_at": time.Now().Add(time.Minute * 30).Unix(),
		"issued_at":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(g.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
