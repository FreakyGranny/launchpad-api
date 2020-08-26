package handlers

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	dateLayout     = "2006-01-02"
	dateTimeLayout = "2006-01-02 15:04:05"
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
		return time.Parse(dateLayout, value)
	}

	return time.Time{}, nil
}

func parseDateTime(value string) (time.Time, error) {
	if value != "" {
		return time.Parse(dateTimeLayout, value)
	}

	return time.Time{}, nil
}
