package api

import (
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/FreakyGranny/launchpad-api/db"
)

// ProjectDonation for project donations response
type ProjectDonation struct {
    ID        uint    `json:"id"`
    User      db.User `json:"user"`
    Locked    bool    `json:"locked"`
    Paid      bool    `json:"paid"`
}

// GetDonation return list of users
func GetDonation(c echo.Context) error {
	projectParam := c.QueryParam("project_id")
	projectID, _ := strconv.Atoi(projectParam)

	dbClient := db.GetDbClient()
	var donations []db.Donation
	
	if projectID != 0 {	
		dbClient.Preload("User").Where("project_id = ?", projectID).Find(&donations)
		projectDonations := make([]ProjectDonation, 0)

		for _, donation := range(donations) {
			projectDonations = append(projectDonations, ProjectDonation{
				ID: donation.ID,
				User: donation.User,
				Locked: donation.Locked,
				Paid: donation.Paid,
			})
		}
		return c.JSON(http.StatusOK, projectDonations)
	}
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["id"].(float64)

	dbClient.Where("user_id = ?", int(userID)).Find(&donations)

	return c.JSON(http.StatusOK, donations)
}
