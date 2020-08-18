package handlers

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

func errorResponse(message string) map[string]string {
	return map[string]string{
		"error": message,
	}
}

func getUserIDFromToken(t interface{}) (int, error) {
	userToken, ok := t.(*jwt.Token)
	if !ok {
		return 0, errors.New("invalid token")
	}
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["id"].(float64)

	return int(userID), nil
}
