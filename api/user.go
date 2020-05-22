package api

import (
	"net/http"
	"strconv"

	"github.com/FreakyGranny/launchpad-api/db"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)


type extendedUser struct {
    ID            uint            `json:"id"`
    Username      string          `json:"username"`
    FirstName     string          `json:"first_name"`
    LastName      string          `json:"last_name"`
    Avatar        string          `json:"avatar"`
    IsStaff       bool            `json:"is_staff"`
    ProjectCount  uint            `json:"project_count"`
	SuccessRate   float32         `json:"success_rate"`
	Participation *[]participation `json:"participation"`
}

type participation struct {
	Cnt           uint `json:"count"`
	ProjectTypeID uint `json:"id"`
}


func extendUser(user db.User) extendedUser {
	dbClient := db.GetDbClient()
	var participations []participation

	dbClient.Table("donations as d").Select("count(d.id) as cnt, p.project_type_id").
						  Joins("left join projects as p on p.id = d.project_id").
						  Where("user_id = ?", user.ID).
						  Group("p.project_type_id").Scan(&participations)

	return extendedUser{
		ID: user.ID,
		Username: user.Username,
		FirstName: user.FirstName,
		LastName: user.LastName,
		Avatar: user.Avatar,
		IsStaff: user.IsStaff,
		ProjectCount: user.ProjectCount,
		SuccessRate: user.SuccessRate,
		Participation: &participations,
	}
}

// GetUsers return list of users
func GetUsers(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["id"].(float64)

	dbClient := db.GetDbClient()
	var user db.User

	if err := dbClient.First(&user, uint(userID)).Error; gorm.IsRecordNotFoundError(err) {
		return c.JSON(http.StatusNotFound, nil)
	}
	
	return c.JSON(http.StatusOK, extendUser(user))
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

	return c.JSON(http.StatusOK, extendUser(user))
}
