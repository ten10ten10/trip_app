package security

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

type PasswordGenerator interface {
	GeneratePassword() (string, string, error)
	ComparePassword(hashedPassword string, rawPassword string) error
}

type passwordGenerator struct{}

func NewPasswordGenerator() PasswordGenerator {
	return &passwordGenerator{}
}

func (r *passwordGenerator) GeneratePassword() (string, string, error) {
	bytes := make([]byte, 12)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", err
	}
	initialPassword := base64.URLEncoding.EncodeToString(bytes)

	hash, err := bcrypt.GenerateFromPassword([]byte(initialPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}
	hashPassword := string(hash)

	return initialPassword, hashPassword, nil
}

func (r *passwordGenerator) ComparePassword(hashedPassword string, rawPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(rawPassword))
}
