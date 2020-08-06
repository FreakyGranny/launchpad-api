package handlers

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jonboulle/clockwork"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"github.com/FreakyGranny/launchpad-api/internal/app/auth"
	"github.com/FreakyGranny/launchpad-api/internal/app/models"
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

// AuthHandler ...
type AuthHandler struct {
	Secret    string
	UserModel models.UserImpl
	Provider  auth.Provider
	Clock     clockwork.Clock
}

// NewAuthHandler ...
func NewAuthHandler(s string, u models.UserImpl, p auth.Provider, c clockwork.Clock) *AuthHandler {
	return &AuthHandler{
		Secret:    s,
		UserModel: u,
		Provider:  p,
		Clock:     c,
	}
}

// Login route returns token
func (h *AuthHandler) Login(c echo.Context) error {
	request := new(TokenRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	data, err := h.Provider.GetAccessToken(request.Code)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusUnauthorized, nil)
	}
	user, userExist := h.UserModel.FindByID(data.UserID)
	user.ID = data.UserID
	user.Email = data.Email

	userData, err := h.Provider.GetUserData(data.UserID, data.AccessToken)
	if err != nil {
		log.Error(err)
		log.Error("unable to get user data")
		return c.JSON(http.StatusInternalServerError, errorResponse("unable create/update user"))
	}
	user.Username = userData.Username
	user.FirstName = userData.FirstName
	user.LastName = userData.LastName
	user.Avatar = userData.Avatar

	if !userExist {
		_, err = h.UserModel.Create(user)
	} else {
		_, err = h.UserModel.Update(user)
	}
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errorResponse("unable to get user data"))
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["admin"] = user.IsAdmin
	claims["exp"] = h.Clock.Now().Add(time.Second * time.Duration(data.Expires)).Unix()

	t, err := token.SignedString([]byte(h.Secret))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}
