package handlers

import (
	"errors"
	"time"

	"github.com/FreakyGranny/launchpad-api/internal/app"
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

func parseDate(value string) (time.Time, error) {
	if value != "" {
		return time.Parse(app.DateLayout, value)
	}

	return time.Time{}, nil
}

func parseDateTime(value string) (time.Time, error) {
	if value != "" {
		return time.Parse(app.DateTimeLayout, value)
	}

	return time.Time{}, nil
}
