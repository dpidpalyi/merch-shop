package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID int, secret string, tokenExpiry time.Duration) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id": userID,
		"iat":     now.Unix(),
		"exp":     now.Add(tokenExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
