package api

import (
	"net/http"
	"strconv"

	// "github.com/labstack/gommon/log"
	"github.com/jinzhu/gorm"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/FreakyGranny/launchpad-api/db"
	"github.com/FreakyGranny/launchpad-api/misc"
)

type createRequest struct {
	ProjectID uint `json:"project"`
	Payment   uint `json:"payment"`
}

type updateRequest struct {
	Paid    bool `json:"paid,omitempty"`
	Payment uint `json:"payment,omitempty"`
}

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

// CreateDonation return list of users
func CreateDonation(c echo.Context) error {
	request := new(createRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	dbClient := db.GetDbClient()
	var project db.Project
	
	if err := dbClient.First(&project, request.ProjectID).Error; gorm.IsRecordNotFoundError(err) {
		return c.JSON(http.StatusBadRequest, nil)
	}
	if project.Closed || project.Locked || !project.Published {
		return c.JSON(http.StatusBadRequest, nil)
	}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["id"].(float64)

	var donationCount uint
	dbClient.Model(&db.Donation{}).Where("project_id = ? AND user_id = ?", request.ProjectID, uint(userID)).Count(&donationCount)
	if donationCount > 0 {
		return c.JSON(http.StatusForbidden, nil)
	}
	newDonation := db.Donation{
		ProjectID: request.ProjectID,
		Payment: request.Payment,
		UserID: uint(userID),
	}
	dbClient.Create(&newDonation)
	ch := misc.GetRecalcPipe()
	ch <- newDonation.ProjectID

	return c.JSON(http.StatusOK, newDonation)
}

// DeleteDonation delete not locked donation
func DeleteDonation(c echo.Context) error {
	idParam := c.Param("id")
	donationID, _ := strconv.Atoi(idParam)

	dbClient := db.GetDbClient()
	var donation db.Donation
	
	if err := dbClient.Preload("User").First(&donation, donationID).Error; gorm.IsRecordNotFoundError(err) {
		return c.JSON(http.StatusNotFound, nil)
	}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["id"].(float64)

	if donation.Locked || donation.User.ID != uint(userID) {
		return c.JSON(http.StatusForbidden, nil)
	}
	dbClient.Delete(&donation)
	ch := misc.GetRecalcPipe()
	ch <- donation.ProjectID

	return c.JSON(http.StatusNoContent, nil)
}

// UpdateDonation update not locked donation
func UpdateDonation(c echo.Context) error {
	request := new(updateRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	idParam := c.Param("id")
	donationID, _ := strconv.Atoi(idParam)

	dbClient := db.GetDbClient()
	var donation db.Donation
	
	if err := dbClient.Preload("Project").Preload("User").First(&donation, donationID).Error; gorm.IsRecordNotFoundError(err) {
		return c.JSON(http.StatusNotFound, nil)
	}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["id"].(float64)

	if donation.Locked && donation.Project.OwnerID == uint(userID) {
		dbClient.Model(&donation).Update("paid", request.Paid)
		ch := misc.GetHarvestPipe()
		ch <- donation.ProjectID

		return c.JSON(http.StatusOK, donation)
	}
	if !donation.Locked && donation.User.ID == uint(userID) {
		dbClient.Model(&donation).Update("payment", request.Payment)
		ch := misc.GetRecalcPipe()
		ch <- donation.ProjectID
		
		return c.JSON(http.StatusOK, donation)
	}

	return c.JSON(http.StatusForbidden, nil)
}
