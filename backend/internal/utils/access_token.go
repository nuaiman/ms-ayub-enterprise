package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(userID int64, secret string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256, jwt.MapClaims{
			"id":  userID,
			"exp": time.Now().Add(1 * time.Hour).Unix(),
		})

	return token.SignedString([]byte(secret))
}

func VerifyAccessToken(tokenString, secret string) (*int64, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invlaid claims")
	}

	id, ok := claims["id"].(float64)
	if !ok {
		return nil, errors.New("invalid id claim")
	}

	userId := int64(id)

	return &userId, nil
}
