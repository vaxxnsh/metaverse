package auth

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Subject string
	Role    string
}

// VerifyJWT validates the signing method, expiry, and required subject claim.
func VerifyJWT(tokenString, secret string) (Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return Claims{}, errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return Claims{}, errors.New("invalid token claims")
	}
	subject, ok := claims["sub"].(string)
	if !ok || subject == "" {
		return Claims{}, errors.New("invalid subject claim")
	}
	role, _ := claims["role"].(string)
	return Claims{Subject: subject, Role: role}, nil
}
