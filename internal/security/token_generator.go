package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

type TokenGenerator interface {
	GenerateToken() (rawToken string, hashedToken string, err error)
	HashToken(rawToken string) string
}

type tokenGenerator struct{}

func NewTokenGenerator() TokenGenerator {
	return &tokenGenerator{}
}

func (s *tokenGenerator) GenerateToken() (string, string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", err
	}
	rawToken := base64.URLEncoding.EncodeToString(bytes)

	hash := sha256.Sum256([]byte(rawToken))
	hashedToken := fmt.Sprintf("%x", hash)

	return rawToken, hashedToken, nil
}

func (s *tokenGenerator) HashToken(rawToken string) string {
	hash := sha256.Sum256([]byte(rawToken))
	return fmt.Sprintf("%x", hash)
}
