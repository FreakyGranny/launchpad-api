package handlers

import (
	"net/http"
	"strconv"

	"github.com/FreakyGranny/launchpad-api/internal/app/models"

	"github.com/labstack/echo/v4"
)

// type extendedUser struct {
//     ID            uint            `json:"id"`
//     Username      string          `json:"username"`
//     FirstName     string          `json:"first_name"`
//     LastName      string          `json:"last_name"`
//     Avatar        string          `json:"avatar"`
//     IsStaff       bool            `json:"is_staff"`
//     ProjectCount  uint            `json:"project_count"`
// 	SuccessRate   float32         `json:"success_rate"`
// 	Participation *[]participation `json:"participation"`
// }

// type participation struct {
// 	Cnt           uint `json:"count"`
// 	ProjectTypeID uint `json:"id"`
// }

// func extendUser(user db.User) extendedUser {
// 	dbClient := db.GetDbClient()
// 	var participations []participation

// 	dbClient.Table("donations as d").Select("count(d.id) as cnt, p.project_type_id").
// 						  Joins("left join projects as p on p.id = d.project_id").
// 						  Where("user_id = ?", user.ID).
// 						  Group("p.project_type_id").Scan(&participations)

// 	return extendedUser{
// 		ID: user.ID,
// 		Username: user.Username,
// 		FirstName: user.FirstName,
// 		LastName: user.LastName,
// 		Avatar: user.Avatar,
// 		IsStaff: user.IsStaff,
// 		ProjectCount: user.ProjectCount,
// 		SuccessRate: user.SuccessRate,
// 		Participation: &participations,
// 	}
// }

// UserHandler ...
type UserHandler struct {
	UserModel models.UserImpl
}

// NewUserHandler ...
func NewUserHandler(u models.UserImpl) *UserHandler {
	return &UserHandler{UserModel: u}
}

// // GetCurrentUser return surrent user
// func GetCurrentUser(c echo.Context) error {
// 	userToken := c.Get("user").(*jwt.Token)
// 	claims := userToken.Claims.(jwt.MapClaims)
// 	userID := claims["id"].(float64)

// 	dbClient := db.GetDbClient()
// 	var user db.User

// 	if err := dbClient.First(&user, uint(userID)).Error; gorm.IsRecordNotFoundError(err) {
// 		return c.JSON(http.StatusNotFound, nil)
// 	}

// 	return c.JSON(http.StatusOK, extendUser(user))
// }

// GetUser return specific user
func (h *UserHandler) GetUser(c echo.Context) error {
	intID, _ := strconv.Atoi(c.Param("id"))
	user, ok := h.UserModel.FindByID(intID)
	if !ok {
		return c.JSON(http.StatusNotFound, nil)
	}

	return c.JSON(http.StatusOK, user)
	// return c.JSON(http.StatusOK, extendUser(user))
}