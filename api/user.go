package api

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/FreakyGranny/launchpad-api/db"
)

// GetUsers return list of users
func GetUsers(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["id"].(float64)

	dbClient := db.GetDbClient()
	var user db.User
	dbClient.First(&user, int(userID))

	return c.JSON(http.StatusOK, user)
}
