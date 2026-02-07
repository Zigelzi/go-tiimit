package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const PasswordMinLength = 6

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

func IsWeakPassword(password string) bool {
	return len(password) < PasswordMinLength
}

func CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
