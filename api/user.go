package api

import (
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
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

// GetUser return specific user
func GetUser(c echo.Context) error {
	userParam := c.Param("id")
	userID, _ := strconv.Atoi(userParam)

	dbClient := db.GetDbClient()
	var user db.User

	if err := dbClient.First(&user, userID).Error; gorm.IsRecordNotFoundError(err) {
		return c.JSON(http.StatusNotFound, nil)
	}

	return c.JSON(http.StatusOK, user)
}
