package api

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/FreakyGranny/launchpad-api/db"
)

// GetDonation return list of users
func GetDonation(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["id"].(float64)

	dbClient := db.GetDbClient()
	var donations []db.Donation
	dbClient.Where("user_id = ?", int(userID)).Find(&donations)

	return c.JSON(http.StatusOK, donations)
}
