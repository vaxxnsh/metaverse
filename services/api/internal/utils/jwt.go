package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vaxxnsh/metaverse/api/internal/config"
)

var secret = config.Load().JWTSecret

func GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
