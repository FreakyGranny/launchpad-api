package api

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/jinzhu/gorm"

	"github.com/FreakyGranny/launchpad-api/db"
	"github.com/FreakyGranny/launchpad-api/misc"
)

// TokenRequest - request for auth token
type TokenRequest struct {
	Code string `json:"code"`
}

func errorResponse(message string) map[string]string {
	return map[string]string{
		"error": message,
	}
}

// Login route returns token
func Login(c echo.Context) error {
	request := new(TokenRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	vkClient := misc.GetVkClient()
	data, err := vkClient.GetAccessToken(request.Code)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusUnauthorized, nil)
	}
	
	userData, err := vkClient.GetUserData(data.UserID, data.AccessToken)	
	if err != nil {
		log.Error("unable to get user data")
		return c.JSON(http.StatusUnauthorized, errorResponse("unable to get user data"))
	}

	user := db.User{
		ID: data.UserID,
		Email: data.Email,
		Username: userData.Username,
		FirstName: userData.FirstName,
		LastName: userData.LastName,
		Avatar: userData.Avatar,
		IsStaff: false,
	}

	client := db.GetDbClient()

	if err := client.First(&user, data.UserID).Error; gorm.IsRecordNotFoundError(err) {
		client.Create(&user)
	} else {
		client.Save(&user)
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["admin"] = user.IsStaff
	claims["exp"] = time.Now().Add(time.Second * time.Duration(data.Expires)).Unix()

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}
