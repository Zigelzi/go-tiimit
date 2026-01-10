package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func MakeToken() (string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	hexToken := hex.EncodeToString(tokenBytes)
	return hexToken, nil
}
